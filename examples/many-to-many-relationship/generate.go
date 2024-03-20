//go:build generate

//go:generate protoc -I $PROTOC_INCLUDE -I ../../../ --go_out=../../../../../ --go-crud_out=../../../../../ protoc-gen-crud/examples/many-to-many-relationship/user.proto protoc-gen-crud/examples/many-to-many-relationship/do_not_export.proto
package link_via_unique_id_single
