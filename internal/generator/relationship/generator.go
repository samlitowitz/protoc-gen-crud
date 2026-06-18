// REFURL: https://github.com/grpc-ecosystem/grpc-gateway/blob/main/protoc-gen-grpc-gateway/internal/gengateway/generator.go
package relationship

import (
	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
	gen "github.com/samlitowitz/protoc-gen-crud/internal/generator"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

type generator struct {
	reg        *descriptor.Registry
	optionsPkg string
}

func New(reg *descriptor.Registry, opts ...Option) gen.Generator {
	options := options{}
	for _, o := range opts {
		o.apply(&options)
	}

	return &generator{
		reg:        reg,
		optionsPkg: options.optionsPkg,
	}
}

func (g *generator) Generate(targets []*descriptor.File) ([]*descriptor.ResponseFile, error) {
	var files []*descriptor.ResponseFile
	for _, file := range targets {
		if len(file.Relationships) == 0 {
			continue
		}
		code, err := g.generate(file)
		if err != nil {
			return nil, err
		}
		if code == "" {
			continue
		}
		files = append(files, &descriptor.ResponseFile{
			CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(file.GeneratedFilenamePrefix + ".crud.proto"),
				Content: proto.String(code),
			},
			GoPkg: file.GoPkg,
		})
	}
	return files, nil
}

func (g *generator) generate(file *descriptor.File) (string, error) {
	params := param{
		File:       file,
		OptionsPkg: g.optionsPkg,
	}

	return applyTemplate(params, g.reg)
}
