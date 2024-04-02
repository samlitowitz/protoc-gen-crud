.PHONY: clean debug release generate

VERSION=`git describe --tags|sed -e "s/\-/\./g"`
BUILD=`date +%FT%T%z`
COMMIT=`git rev-parse HEAD`
TARGETDIR=${GOPATH}/bin
PROJECT_PROTO_INCLUDE=$(shell cd .. && pwd)
PROJECT_PROTO_OUT=$(shell cd ../../../ && pwd)

LDFLAGS_DEB=-ldflags "-X main.Date=${BUILD} -X main.Commit=${COMMIT}"
LDFLAGS_REL=-ldflags "-s -w -X main.Version=${VERSION} -X main.Date=${BUILD} -X main.Commit=${COMMIT}"

clean:
	rm -f ${TARGETDIR}/protoc-gen-go-crud

release:
	go build ${LDFLAGS_REL} -o ${TARGETDIR}/protoc-gen-go-crud ./cmd/protoc-gen-go-crud

debug:
	go build ${LDFLAGS_DEB} -o ${TARGETDIR}/protoc-gen-go-crud ./cmd/protoc-gen-go-crud

generate:
	echo ${PROJECT_PROTO_INCLUDE}
	PROJECT_PROTO_INCLUDE=${PROJECT_PROTO_INCLUDE} PROJECT_PROTO_OUT=${PROJECT_PROTO_OUT} go generate -v $(go list ./... | grep -v examples)
