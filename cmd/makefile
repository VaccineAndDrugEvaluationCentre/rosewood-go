# EXE is binary name
EXE=carpenter
VERSION=0.1.0
BUILD==`date +%FT%T%z`
#to use github commit tag
#BUILD=`git rev-parse HEAD`  
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

.PHONY: clean install 