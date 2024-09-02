package sql

import (
	"fmt"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
	gen "github.com/samlitowitz/protoc-gen-crud/internal/generator"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

type generator struct {
	reg *descriptor.Registry
}

func New(reg *descriptor.Registry) gen.Generator {
	return &generator{
		reg: reg,
	}
}

func (g *generator) Generate(targets []*descriptor.File) ([]*descriptor.ResponseFile, error) {
	var files []*descriptor.ResponseFile
	for _, file := range targets {
		if len(file.Implementations) == 0 {
			continue
		}
		code, err := g.generate(file)
		if err != nil {
			return nil, fmt.Errorf("pgsql: generate: %s: %v", file.GetName(), err)
		}
		files = append(files, &descriptor.ResponseFile{
			CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(file.GeneratedFilenamePrefix + ".pgsql.sql"),
				Content: proto.String(code),
			},
			GoPkg: file.GoPkg,
		})
	}
	return files, nil
}

func (g *generator) generate(file *descriptor.File) (string, error) {
	param := param{
		File: file,
	}
	return applyTemplate(param, g.reg)
}
