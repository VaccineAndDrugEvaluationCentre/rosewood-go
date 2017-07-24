# EXE is binary name
EXE=carpenter
VERSION := $(shell git describe --always --long --dirty)
BUILD==`date +%FT%T%z`


PKG := github.com/drgo/carpenter 
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

README=../README.md

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

run: 
	go run carpenter.go
#macos buid
build:
	go build ${LDFLAGS} -o ${EXE}

# Windows build
cross:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${EXE}.exe	

clean:
	rm ${EXE}
	rm ${EXE}.exe
#	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi 

doc: 
#awk -F. '{$NF+=1; OFS="."; print $0}'
	sed -ne 's/[^0-9]*\(\([0-9]\.\)\{0,4\}[0-9][^.]\).*/\1/p' '${README}'
#	sed -i '' 's/Version /& ${VERSION}/' '${README}'
#	sed -i '' 's/BUILD/& ${BUILD}/' '${README}'	

aVERSION=`git describe`

bump:
	dch -v $(aVERSION)

.PHONY: clean install doc