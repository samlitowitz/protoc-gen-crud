//go:build generate

//go:generate protoc -I $PROTOC_INCLUDE -I ../../ --go_out=../../../../ --go_opt=default_api_level=API_OPAQUE protoc-gen-crud/options/relationships/direction.proto protoc-gen-crud/options/relationships/type.proto
//go:generate protoc -I $PROTOC_INCLUDE -I ../../ --go_out=../../../../ --go_opt=default_api_level=API_OPAQUE protoc-gen-crud/options/relationship.proto protoc-gen-crud/options/crud.proto protoc-gen-crud/options/annotations.proto

package internal
