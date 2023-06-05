PROJECT_BIN_NAME=local-chart-version
PROJECT_NAME=helm-$(PROJECT_BIN_NAME)

HELM_HOME ?= $(shell helm home)
VERSION ?= $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)

pkgs = $(shell go list ./... | grep -v /vendor/ | grep -v /test/)

.PHONY: clean
clean:
	rm -rf $(PROJECT_BIN_NAME) ./dist

HAS_GIT := $(shell command -v git;)
HAS_GORELEASER := $(shell command -v goreleaser;)

.PHONY: show-version
show-version:
	@echo $(VERSION)

.PHONY: gendoc
gendoc: build
	mkdir -p docs
	./$(PROJECT_BIN_NAME) generate-documentation
	make clean --always-make

.PHONY: build
build: test
	go build -v -o $(PROJECT_BIN_NAME)

.PHONY: test
test: 
	go test ./...

.PHONY: build-cross
build-cross: clean test gendoc
ifndef HAS_GORELEASER
	$(error You must install goreleaser)
endif
	goreleaser --snapshot

.PHONY: install
install: build
	mkdir -p $(HELM_HOME)/plugins/$(PROJECT_NAME)
	cp $(PROJECT_BIN_NAME) $(HELM_HOME)/plugins/$(PROJECT_NAME)/
	cp plugin.yaml $(HELM_HOME)/plugins/$(PROJECT_NAME)/

.PHONY: release
release: clean test gendoc
ifndef HAS_GIT
	$(error You must install Git)
endif
ifndef HAS_GORELEASER
	$(error You must install goreleaser)
endif
	git tag -a v$(VERSION) -m "release v$(VERSION)"
	git push origin v$(VERSION)
	goreleaser