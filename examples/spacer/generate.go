//go:build generate

//go:generate protoc -I $PROTOC_INCLUDE -I ../../../ --go_out=../../../../../ --go-crud_out=../../../../../  protoc-gen-crud/examples/spacer/tag.proto protoc-gen-crud/examples/spacer/entity.proto
package partial_creates_updates
