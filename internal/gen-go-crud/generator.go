// REFURL: https://github.com/grpc-ecosystem/grpc-gateway/blob/main/protoc-gen-grpc-gateway/internal/gengateway/generator.go
package gen_go_crud

import (
	"go/format"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
	gen "github.com/samlitowitz/protoc-gen-crud/internal/generator"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

type generator struct {
	reg         *descriptor.Registry
	baseImports []descriptor.GoPackage
}

func New(reg *descriptor.Registry) gen.Generator {
	return &generator{
		reg:         reg,
		baseImports: make([]descriptor.GoPackage, 0),
	}
}

func (g *generator) Generate(targets []*descriptor.File) ([]*descriptor.ResponseFile, error) {
	var files []*descriptor.ResponseFile
	for _, file := range targets {
		code, err := g.generate(file)
		if err != nil {
			return nil, err
		}
		formatted, err := format.Source([]byte(code))
		if err != nil {
			return nil, err
		}
		files = append(files, &descriptor.ResponseFile{
			CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(file.GeneratedFilenamePrefix + ".pb.crud.go"),
				Content: proto.String(string(formatted)),
			},
			GoPkg: file.GoPkg,
		})
	}
	return files, nil
}

func (g *generator) generate(file *descriptor.File) (string, error) {
	pkgSeen := make(map[string]bool)
	var imports []descriptor.GoPackage
	for _, pkg := range g.baseImports {
		pkgSeen[pkg.Path] = true
		imports = append(imports, pkg)
	}

	params := param{
		File: file,
	}

	return applyTemplate(params, g.reg)
}
