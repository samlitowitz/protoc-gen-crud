//go:build generate

//go:generate sh -c "protoc -I $PROTOC_INCLUDE -I $PROJECT_PROTO_INCLUDE  --go_out=$PROJECT_PROTO_OUT --go-crud_out=$PROJECT_PROTO_OUT $PROJECT_PROTO_INCLUDE/protoc-gen-crud/test-cases/candidate-key/*.proto"

package candidate_key
