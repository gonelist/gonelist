
# repo上最新tag的commitID
TAG_COMMITID    := $(shell git rev-list --tags --max-count=1)
BUILD_VERSION   := $(shell git describe --tags $(TAG_COMMITID))
ifdef VERSION
	git checkout $(VERSION)
	BUILD_VERSION := $(VERSION)
endif

HEAD_COMMIT     := $(shell git rev-parse HEAD)
BUILD_DATE		:= $(shell date +'%Y-%m-%dT%H:%M:%S')
GIT_TREE_STATE := $(shell git status --porcelain 2>/dev/null)

ifeq ($(GIT_TREE_STATE),)
	GIT_TREE_STATE := clean
	COMMIT_ID := $(HEAD_COMMIT)
else
	GIT_TREE_STATE := dirty
	BUILD_VERSION := $(BUILD_VERSION)-dirty
endif

ifneq ($(HEAD_COMMIT),$(TAG_COMMITID))
	BUILD_VERSION := $(BUILD_VERSION)-dirty
	COMMIT_ID := $(HEAD_COMMIT)
endif


GO ?= go
LDFLAGS := -ldflags "-X main.Version=${BUILD_VERSION} \
                             -X main.buildDate=${BUILD_DATE} \
                             -X main.gitCommit=${COMMIT_ID} \
                             -X main.gitTreeState=${GIT_TREE_STATE}"
fmt:
	go fmt ./...
download:
	@$(GO) mod download
version: download build

build: clean
	$(GO) build -o gonelist $(LDFLAGS) main.go
test:
	@echo $(COMMIT_SHA1) $(TAG_COMMIT) $(BUILD_VERSION) $(GIT_TREE_STATE)
status:
	@git status --porcelain 2>/dev/null
release:
	bash build.sh
clean:
	-rm -f gonelist main
	-rm -rf release/*

.PHONY : fmt download build test release clean status version tag

.EXPORT_ALL_VARIABLES:
GO111MODULE = on
CGO_ENABLED = 0
GOPROXY = https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,https://goproxy.io,direct



