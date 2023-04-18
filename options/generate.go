//go:build generate

//go:generate protoc -I $PROTOC_INCLUDE -I ../../ --go_out=../../../../ protoc-gen-crud/options/crud.proto protoc-gen-crud/options/annotations.proto

package internal
