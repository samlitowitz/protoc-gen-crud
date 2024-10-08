//go:build generate

//go:generate sh -c "protoc -I $PROTOC_INCLUDE -I $PROJECT_PROTO_INCLUDE  --go_out=$PROJECT_PROTO_OUT --go-crud_out=$PROJECT_PROTO_OUT --go-crud_opt=format_output=false $PROJECT_PROTO_INCLUDE/protoc-gen-crud/test-cases/as-timestamp-field/*.proto"
//go:generate sh -c "protoc -I $PROTOC_INCLUDE -I $PROJECT_PROTO_INCLUDE  --go_out=$PROJECT_PROTO_OUT --go-crud_out=$PROJECT_PROTO_OUT $PROJECT_PROTO_INCLUDE/protoc-gen-crud/test-cases/as-timestamp-field/*.proto"

package as_timestamp_field
