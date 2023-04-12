//go:build generate

//go:generate protoc -I $PROTOC_INCLUDE -I ../../../ --go_out=../../../../../ protoc-gen-crud/assets/protobuf/examples/book-list/books.proto
package book_list_template
