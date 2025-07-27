############################
# -- go-test
############################
FROM golang:1.24-alpine AS go-test

#RUN adduser -D -g '' appuser
WORKDIR $GOPATH/src/github.com/samlitowitz/protoc-gen-crud
COPY . .

ENV GO111MODULE=on
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod download
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod verify

CMD [ "go", "test", "-v", "./test-cases/..." ]

############################
# -- go-test-dlv
############################
FROM go-test AS go-test-dlv

# Install delver
RUN --mount=type=cache,mode=0755,target=/go/bin go install github.com/go-delve/delve/cmd/dlv@latest \
  && cp /go/bin/dlv /usr/local/bin/dlv
