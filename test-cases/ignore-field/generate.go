//go:build generate

//go:generate sh -c "protoc -I $PROTOC_INCLUDE -I $PROJECT_PROTO_INCLUDE  --go_out=$PROJECT_PROTO_OUT --go-crud_out=$PROJECT_PROTO_OUT $PROJECT_PROTO_INCLUDE/protoc-gen-crud/test-cases/ignore-field/*.proto"

package ignore_field
