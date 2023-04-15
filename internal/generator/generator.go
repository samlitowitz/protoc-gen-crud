// REFURL: https://github.com/grpc-ecosystem/grpc-gateway/blob/main/internal/generator/generator.go
package generator

import "github.com/samlitowitz/protoc-gen-crud/internal/descriptor"

// Generator is an abstraction of code generators.
type Generator interface {
	// Generate generates output files from input .proto files.
	Generate(targets []*descriptor.File) ([]*descriptor.ResponseFile, error)
}
