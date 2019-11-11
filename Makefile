GO_CMD ?=go
BIN_NAME := lyrics

build:
	${GO_CMD} build -o ${BIN_NAME} -v */go

install: build
	install -d /usr/local/bin/
	install -m 755 ./${BIN_NAME} /usr/local/bin/${BIN_NAME}

clean:
	rm -f ./${BIN_NAME}*

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GO_CMD} build -o ${BIN_NAME} -v *.go

build-osx:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 ${GO_CMD} build -o ${BIN_NAME}_osx -v *.go

build-all: build-linux build-osx
