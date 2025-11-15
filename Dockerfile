############################
# -- go-build
############################
FROM golang:1.24-alpine AS go-build

RUN apk add --update curl git && \
    rm -rf /var/cache/apk/*

# Install protoc
ENV PROTOC_INCLUDE=/usr/local/include
RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v30.2/protoc-30.2-linux-x86_64.zip && \
    unzip protoc-30.2-linux-x86_64.zip -d /usr/local


# Install protoc-gen-go
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Install Task
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

WORKDIR $GOPATH/src/github.com/samlitowitz/protoc-gen-crud
COPY . .

ENV GOBIN=/usr/local/bin
ENV GO111MODULE=on
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod download
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod verify

RUN task generate-options && \
    task release

############################
# -- go-test
############################
FROM go-build AS go-test

RUN task generate-test-cases

WORKDIR $GOPATH/src/github.com/samlitowitz/protoc-gen-crud

CMD [ "go", "test", "-v", "./test-cases/..." ]

############################
# -- go-test-dlv
############################
FROM go-test AS go-test-dlv

# Install delver
RUN --mount=type=cache,mode=0755,target=/go/bin go install github.com/go-delve/delve/cmd/dlv@latest \
  && cp /go/bin/dlv /usr/local/bin/dlv
