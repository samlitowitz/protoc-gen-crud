//go:build generate

//go:generate sh -c "protoc -I $PROTOC_INCLUDE -I $PROJECT_PROTO_INCLUDE --go_out=$PROJECT_PROTO_INCLUDE --go_opt=paths=source_relative --go_opt=default_api_level=API_OPAQUE $PROJECT_PROTO_INCLUDE/protoc_gen_crud/options/v1/*.proto $PROJECT_PROTO_INCLUDE/protoc_gen_crud/options/v1/relationships/*.proto"

package v1
