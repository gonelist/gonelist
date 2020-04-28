
SHELL := /bin/bash

ifdef $(TAG)
	TAG_COMMITID := $(shell git rev-parse $TAG 2>/dev/null)
else
	# 取最新的tag作为版本号
	TAG_COMMITID    := $(shell git rev-list --tags --max-count=1)
	TAG   := $(shell git describe --tags $(TAG_COMMITID))
endif

$(shell git checkout $TAG 1>/dev/null)
BUILD_VERSION :=$(TAG)

BUILD_DATE		:= $(shell date +'%Y-%m-%dT%H:%M:%S')
GIT_TREE_STATE := $(shell git status --porcelain 2>/dev/null)

ifeq ($(GIT_TREE_STATE),)
	GIT_TREE_STATE := clean
else
	GIT_TREE_STATE := dirty
endif

HEAD_COMMIT     := $(shell git rev-parse HEAD)
ifneq ($(HEAD_COMMIT),$(TAG_COMMITID))
	BUILD_VERSION := $(BUILD_VERSION)-dirty
	COMMIT_ID := $(HEAD_COMMIT)
endif


GO ?= go
LDFLAGS := -ldflags "-X main.Version=${BUILD_VERSION} \
                             -X main.buildDate=${BUILD_DATE} \
                             -X main.gitCommit=${COMMIT_ID} \
                             -X main.gitTreeState=${GIT_TREE_STATE}"

version: download clean build

build: clean
	$(GO) build -o gonelist $(LDFLAGS) main.go
var:
	@echo $(COMMIT_SHA1) $(TAG_COMMIT) $(BUILD_VERSION) $(GIT_TREE_STATE)
status:
	@git status --porcelain 2>/dev/null
release:
	bash build.sh
clean:
	-rm -f gonelist main
	-rm -rf release/*

.PHONY : fmt build var release clean status version tag

.EXPORT_ALL_VARIABLES:
GO111MODULE = on
CGO_ENABLED = 0
GOPROXY = https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,https://goproxy.io,direct



