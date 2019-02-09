EXE=carpenter
VERSION=0.6.3
BUILD := $(shell git describe --always --long --dirty)
BUILD_DATE=$(shell date)

# disable gomodule for now
# export GO111MODULE=off
#PKG := github.com/drgo/rosewood 
#PKG_LIST := $(shell go list ${PKG}/... )
#GO_FILES := $(shell find . -name '*.go' )

README=../README.md

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

#default buid
build:
	go build -race ${LDFLAGS} -o ${EXE}

buildprod: 
	go build ${LDFLAGS} -o ${EXE}

# Windows build
buildwin:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build  -race ${LDFLAGS} -o ${EXE}.exe	

buildwinprod:
	GOOS=windows GOARCH=amd64 go build  ${LDFLAGS} -o ${EXE}.exe	

clean:
	rm ${EXE}
	rm ${EXE}.exe
	# go clean

doc: 
#update README with version and build number in line 4 only and only if it had Version
	sed -i '' '4s/.*Version.*/Version ${VERSION} build ${BUILD} on ${BUILD_DATE}/' '${README}'
#update README with the content of the usage.go file
	sed -i 'bak' "/### How to use/,\$$d"  '${README}'
	echo '### How to use' >> '${README}'
	sed -e '1,/Usage/d' usage.go >> '${README}'


cover:
	cd ..
	go test -cover -coverprofile=c.out  
	go tool cover -html=c.out -o coverage.html 
	rm c.out	

showcover:	
	go test -coverprofile=c.out && go tool cover -html=c.out

graph:
	godepgraph "github.com/drgo/rosewood/cmd"  | dot -Tpng -o godepgraph.png	

.PHONY: clean install doc testcover graph
