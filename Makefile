VERSION := $(git describe --tags)
BUILD_TIME := $(date +%Y-%m-%d_%H:%M:%S)
GIT_COMMIT := $(git rev-parse HEAD)

.PHONY: build
build:
	CGO_ENABLED=0 go build \
		-ldflags="-w -s \
		-X 'main.Version=${VERSION}' \
		-X 'main.BuildTime=${BUILD_TIME}' \
		-X 'main.GitCommit=${GIT_COMMIT}'" \
		-trimpath \
		-o build/${VERSION}/adc

.PHONY: build-all
build-all:
	# Linux
	GOOS=linux GOARCH=amd64 make build
	mv build/${VERSION}/adc build/${VERSION}/adc-linux-amd64
	
	# MacOS
	GOOS=darwin GOARCH=amd64 make build
	mv build/${VERSION}/adc build/${VERSION}/adc-darwin-amd64
	
	# Windows
	GOOS=windows GOARCH=amd64 make build
	mv build/${VERSION}/adc build/${VERSION}/adc-windows-amd64.exe