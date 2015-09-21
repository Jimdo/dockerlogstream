PROJECT := dockerlogstream
PROJECT_NAMESPACE := github.com/Jimdo

REGISTRY := registry.example.com

BIN := $(PROJECT)

SOURCE_DIR := $(CURDIR)
BUILD_DIR := $(CURDIR)/.gobuild
BUILD_DIR_SRC := $(BUILD_DIR)/src/$(PROJECT_NAMESPACE)/$(PROJECT)

VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)
DOCKER_TAG := $(shell git describe --tags --always)

ifndef GOOS
	GOOS := $(shell go env GOOS)
endif
ifndef GOARCH
	GOARCH := $(shell go env GOARCH)
endif

default: clean $(BIN)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)/src/$(PROJECT_NAMESPACE)
	ln -s $(SOURCE_DIR) $(BUILD_DIR_SRC)
	cp -R $(SOURCE_DIR)/Godeps/_workspace/src/* $(BUILD_DIR)/src/

$(BIN): VERSION $(BUILD_DIR)
	docker run \
		--rm \
		-v $(CURDIR):/usr/code \
		-e GOPATH=/usr/code/.gobuild:/usr/code/.gobuild/src/$(PROJECT_NAMESPACE)/$(PROJECT)/Godeps/_workspace \
		-e GOOS=$(GOOS) \
		-e GOARCH=$(GOARCH) \
		-w /usr/code \
		golang:1.5 \
		go build -a -ldflags "-X main.ProjectVersion=$(VERSION) -X main.ProjectBuild=$(COMMIT)" -o $(BIN)

clean:
	rm -fr $(BUILD_DIR) $(BIN)

docker-image: GOOS=linux
docker-image: GOARCH=386
docker-image: clean $(BIN)
docker-image:
	docker build -t $(REGISTRY)/$(PROJECT):$(DOCKER_TAG) .

semver-bump:
	go get github.com/giantswarm/semver-bump

release-major: semver-bump
	semver-bump major-release
	$(MAKE) -B release

release-minor: semver-bump
	semver-bump minor-release
	$(MAKE) -B release

release-patch: semver-bump
	semver-bump patch-release
	$(MAKE) -B release

release:
	git ci -m "Release v$(VERSION)" VERSION History.md
	git tag "v$(VERSION)"
