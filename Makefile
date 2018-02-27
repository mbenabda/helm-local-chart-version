PROJECT_BIN_NAME=local-chart-version
PROJECT_NAME=helm-$(PROJECT_BIN_NAME)

HELM_HOME ?= $(shell helm home)
VERSION ?= $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)

.PHONY: clean
clean:
	@rm -rf $(PROJECT_BIN_NAME) ./dist

HAS_GLIDE := $(shell command -v glide;)
HAS_GIT := $(shell command -v git;)
HAS_GORELEASER := $(shell command -v goreleaser;)

.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
ifndef HAS_GIT
	$(error You must install Git)
endif
ifndef HAS_GORELEASER
	go get -u github.com/goreleaser/goreleaser
endif
	glide install --strip-vendor

.PHONY: build
build:
	go build -i -v -o $(PROJECT_BIN_NAME)

.PHONY: build-cross
build-cross: bootstrap
	goreleaser --snapshot

.PHONY: install
install: bootstrap build
	mkdir -p $(HELM_HOME)/plugins/$(PROJECT_NAME)
	cp $(PROJECT_BIN_NAME) $(HELM_HOME)/plugins/$(PROJECT_NAME)/
	cp plugin.yaml $(HELM_HOME)/plugins/$(PROJECT_NAME)/

.PHONY: release
release: clean
	git tag -a v$(VERSION) -m "release v$(VERSION)"
	git push origin v$(VERSION)
	goreleaser