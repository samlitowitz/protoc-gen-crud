//go:build generate

//go:generate sh -c "protoc -I $PROTOC_INCLUDE -I $PROJECT_PROTO_INCLUDE --go_out=$PROJECT_PROTO_OUT --go-crud_out=options_pkg=options.v1:$PROJECT_PROTO_OUT --go_opt=default_api_level=API_OPAQUE $PROJECT_PROTO_INCLUDE/test-cases/relationships-one-to-one/*.proto"

package relationships_one_to_one
