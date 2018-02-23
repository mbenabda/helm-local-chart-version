PROJECT_BIN_NAME=local-chart-version
PROJECT_NAME=helm-$(PROJECT_BIN_NAME)
PROJECT_GH=mbenabda/$(PROJECT_NAME)
PKG:= github.com/$(PROJECT_GH)

HELM_HOME ?= $(shell helm home)
VERSION ?= $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)
TARGETS ?= $(shell cat .build_targets)

LDFLAGS := -X main.Version=$(VERSION)
# Clear the "unreleased" string in BuildMetadata
LDFLAGS += -X $(PKG)/vendor/k8s.io/helm/pkg/version.BuildMetadata=
LDFLAGS += -X $(PKG)/vendor/k8s.io/helm/pkg/version.Version=$(shell grep -A1 "package: k8s.io/helm" glide.yaml | sed -n -e 's/[ ]*version:.*\(v[.0-9]*\).*/\1/p')

.PHONY: clean
clean:
	@rm -rf $(PROJECT_BIN_NAME) ./build ./dist

.PHONY: build-cross
build-cross: LDFLAGS += -s -w -extldflags "-static"
build-cross:
	CGO_ENABLED=0 gox \
		-verbose \
		-output="dist/$(PROJECT_NAME)-{{.OS}}-{{.Arch}}/$(PROJECT_BIN_NAME)" \
		-ldflags "$(LDFLAGS)" \
		-osarch="$(TARGETS)" \
		$(PKG)

.PHONY: dist
dist: export COPYFILE_DISABLE=1 #teach OSX tar to not put ._* files in tar archive
dist: 
	( \
		cd dist && \
		find * -maxdepth 1 -type d -exec cp ../README.md {} \; && \
		find * -maxdepth 1 -type d -exec cp ../plugin.yaml {} \; && \
		ls -d * | xargs -I {} -n1 bash -c "cd {} && tar -zcf ../{}.tgz * && cd .." \
	)

HAS_GLIDE := $(shell command -v glide;)
HAS_GOX := $(shell command -v gox;)
HAS_GIT := $(shell command -v git;)
HAS_GITHUB_RELEASE := $(shell command -v github-release;)

.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
ifndef HAS_GOX
	go get -u github.com/mitchellh/gox
endif
ifndef HAS_GIT
	$(error You must install Git)
endif
ifndef HAS_GITHUB_RELEASE
	go get -u github.com/c4milo/github-release
endif
	glide install --strip-vendor

.PHONY: build
build:
	go build -i -v -o ${PROJECT_BIN_NAME} -ldflags="$(LDFLAGS)"

.PHONY: install
install: bootstrap build
	mkdir -p $(HELM_HOME)/plugins/${PROJECT_NAME}
	cp ${PROJECT_BIN_NAME} $(HELM_HOME)/plugins/${PROJECT_NAME}/
	cp plugin.yaml $(HELM_HOME)/plugins/${PROJECT_NAME}/

.PHONY: release
release: clean bootstrap build-cross dist
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN is undefined)
endif
	git push
	github-release ${PROJECT_GH} v$(VERSION) master "v$(VERSION)" "dist/*.*"
