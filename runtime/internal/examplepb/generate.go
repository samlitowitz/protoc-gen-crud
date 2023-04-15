//go:build generate

//go:generate protoc -I $PROTOC_INCLUDE -I ../../../../ --go_out=../../../../../../ protoc-gen-crud/runtime/internal/examplepb/example.proto
//go:generate protoc -I $PROTOC_INCLUDE -I ../../../../ --go_out=../../../../../../ protoc-gen-crud/runtime/internal/examplepb/non_standard_names.proto
//go:generate protoc -I $PROTOC_INCLUDE -I ../../../../ --go_out=../../../../../../ protoc-gen-crud/runtime/internal/examplepb/proto3.proto
package example
