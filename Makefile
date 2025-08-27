.PHONY: all  clean build


TARGET_OS := linux
TARGET_ARCH := amd64
# ---
MOSTVERSION = $(shell git describe --tags --always)

GIT_COMMIT = $(shell git rev-parse --short HEAD)

BUILD_TIME := $(shell date  "+%Y-%m-%d %H:%M:%S")

VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo v0.0.1)

OUT = detection

all: $(OUT)

build:  main.go
	@go build  -ldflags "-s -w  -X 'main.MostVersion=${MOSTVERSION}' -X 'main.Version=${VERSION}' -X 'main.GitCommit=${GIT_COMMIT}' -X 'main.BuildTime=${BUILD_TIME}'" -o $(OUT)_$(TARGET_OS)_$(TARGET_ARCH)_${VERSION} $<


$(OUT): main.go
	@go build  -ldflags "-s -w  -X 'main.MostVersion=${MOSTVERSION}' -X 'main.Version=${VERSION}' -X 'main.GitCommit=${GIT_COMMIT}' -X 'main.BuildTime=${BUILD_TIME}'" -o ./build/$@_$(TARGET_OS)_$(TARGET_ARCH)_${VERSION} $<
	@cp -r ./config ./build/
	@cd ./build && zip -r ./$@_$(TARGET_OS)_$(TARGET_ARCH)_${VERSION}.zip ./$@_$(TARGET_OS)_$(TARGET_ARCH)_${VERSION} ./config/
	@rm ./build/$@_$(TARGET_OS)_$(TARGET_ARCH)_${VERSION}
	@rm -rf ./build/config
clean:
	@rm -f $(OUT)_$(TARGET_OS)_$(TARGET_ARCH)_v*
	@rm -rf ./build
