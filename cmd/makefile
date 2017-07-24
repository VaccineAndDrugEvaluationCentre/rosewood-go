# EXE is binary name
EXE=carpenter
VERSION=0.2.0
BUILD := $(shell git describe --always --long --dirty)
##BUILD==`date +%FT%T%z`


PKG := github.com/drgo/carpenter 
PKG_LIST := $(shell go list ${PKG}/... )
GO_FILES := $(shell find . -name '*.go' )

README=../README.md

LDFLAGS=-ldflags "-X maign.Version=${VERSION} -X main.Build=${BUILD}"

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
#awk -F. '{$NF+=1; OFS="."; print $0}'
	sed -ne 's/[^0-9]*\(\([0-9]\.\)\{0,4\}[0-9][^.]\).*/\1/p' '${README}'
#	sed -i '' 's/Version /& ${VERSION}/' '${README}'
#	sed -i '' 's/BUILD/& ${BUILD}/' '${README}'	

.PHONY: clean install doc