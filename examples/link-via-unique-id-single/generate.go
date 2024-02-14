//go:build generate

//go:generate protoc -I $PROTOC_INCLUDE -I ../../../ --go_out=../../../../../ --go-crud_out=../../../../../ protoc-gen-crud/examples/link-via-unique-id-single/user.proto
package link_via_unique_id_single
