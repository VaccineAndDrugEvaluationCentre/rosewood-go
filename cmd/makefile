# EXE is binary name
EXE=carpenter
VERSION=0.2.1
BUILD := $(shell git describe --always --long --dirty)
BUILD_DATE=$(shell date)


PKG := github.com/drgo/carpenter 
PKG_LIST := $(shell go list ${PKG}/... )
GO_FILES := $(shell find . -name '*.go' )

README=../README.md

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

run: 
	go run carpenter.go
#default buid
build:
	go build ${LDFLAGS} -o ${EXE}

# Windows build
buildwin:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${EXE}.exe	

clean:
	rm ${EXE}
	rm ${EXE}.exe

doc: 
#update README with version and build number in line 4 only and only if it had Version
	sed -i '' '4s/.*Version.*/Version ${VERSION} build ${BUILD} on ${BUILD_DATE}/' '${README}'


.PHONY: clean install doc