VERSION=`git describe --tags|sed -e "s/\-/\./g"`
BUILD=`date +%FT%T%z`
COMMIT=`git rev-parse HEAD`
TARGETDIR=${GOPATH}/bin

LDFLAGS_DEB=-ldflags "-X main.Date=${BUILD} -X main.Commit=${COMMIT}"
LDFLAGS_REL=-ldflags "-s -w -X main.Version=${VERSION} -X main.Date=${BUILD} -X main.Commit=${COMMIT}"

clean:
	rm -f ${TARGETDIR}/protoc-gen-go-crud

release:
	go build ${LDFLAGS_REL} -o ${TARGETDIR}/protoc-gen-go-crud ./cmd/protoc-gen-go-crud

debug:
	go build ${LDFLAGS_DEB} -o ${TARGETDIR}/protoc-gen-go-crud ./cmd/protoc-gen-go-crud

.PHONY: clean debug release
