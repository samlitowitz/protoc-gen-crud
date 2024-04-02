// REFURL: https://github.com/grpc-ecosystem/grpc-gateway/blob/main/protoc-gen-grpc-gateway/internal/gengateway/generator.go
package crud

import (
	"fmt"
	"go/format"
	"path"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
	gen "github.com/samlitowitz/protoc-gen-crud/internal/generator"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	defaultFormatOutput = true
)

type generator struct {
	reg         *descriptor.Registry
	baseImports []descriptor.GoPackage

	formatOutput bool
}

func New(reg *descriptor.Registry, opts ...Option) gen.Generator {
	options := options{
		formatOutput: defaultFormatOutput,
	}
	for _, o := range opts {
		o.apply(&options)
	}

	var imports []descriptor.GoPackage
	for _, pkgpath := range []string{
		"context",
		"github.com/samlitowitz/protoc-gen-crud/expressions",
	} {
		pkg := descriptor.GoPackage{
			Path: pkgpath,
			Name: path.Base(pkgpath),
		}
		if err := reg.ReserveGoPackageAlias(pkg.Name, pkg.Path); err != nil {
			for i := 0; ; i++ {
				alias := fmt.Sprintf("%s_%d", pkg.Name, i)
				if err := reg.ReserveGoPackageAlias(alias, pkg.Path); err != nil {
					continue
				}
				pkg.Alias = alias
				break
			}
		}
		imports = append(imports, pkg)
	}
	return &generator{
		reg:         reg,
		baseImports: imports,

		formatOutput: options.formatOutput,
	}
}

func (g *generator) Generate(targets []*descriptor.File) ([]*descriptor.ResponseFile, error) {
	var files []*descriptor.ResponseFile
	for _, file := range targets {
		code, err := g.generate(file)
		if err != nil {
			return nil, err
		}
		output := code
		if g.formatOutput {
			formatted, err := format.Source([]byte(code))
			if err != nil {
				return nil, err
			}
			output = string(formatted)
		}
		files = append(files, &descriptor.ResponseFile{
			CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(file.GeneratedFilenamePrefix + ".pb.crud.go"),
				Content: proto.String(output),
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
		File:    file,
		Imports: imports,
	}

	return applyTemplate(params, g.reg)
}
