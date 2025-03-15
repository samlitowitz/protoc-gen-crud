// REFURL: https://github.com/grpc-ecosystem/grpc-gateway/blob/main/protoc-gen-grpc-gateway/internal/gengateway/generator.go
package crud

import (
	"fmt"
	"go/format"
	"path"

	crudOptions "github.com/samlitowitz/protoc-gen-crud/options"

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
	for _, pkgpath := range []string{} {
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
	}
}

func (g *generator) Generate(targets []*descriptor.File) ([]*descriptor.ResponseFile, error) {
	var files []*descriptor.ResponseFile
	for _, file := range targets {
		if len(file.Implementations) == 0 {
			continue
		}
		if _, ok := file.Implementations[crudOptions.Implementation_IMPLEMENTATION_SQLITE]; !ok {
			continue
		}
		code, err := g.generate(file)
		if err != nil {
			return nil, fmt.Errorf("sqlite: generate: %s: %v", file.GetName(), err)
		}

		output := code
		if g.formatOutput {
			formatted, err := format.Source([]byte(code))
			if err != nil {
				return nil, fmt.Errorf("sqlite: format: %s: %v", file.GetName(), err)
			}
			output = string(formatted)
		}

		files = append(files, &descriptor.ResponseFile{
			CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(file.GeneratedFilenamePrefix + ".pb.crud.sqlite.go"),
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

	for _, msg := range file.Messages {
		imports = append(imports, g.addMessagePathParamImports(file, msg, pkgSeen)...)
		imports = append(imports, g.addCrudPathParamImports(msg, pkgSeen)...)
	}

	params := param{
		File:    file,
		Imports: imports,
	}

	return applyTemplate(params, g.reg)
}

// addMessagePathParamImports handles adding import of message path parameter go packages
func (g *generator) addMessagePathParamImports(file *descriptor.File, msg *descriptor.Message, pkgSeen map[string]bool) []descriptor.GoPackage {
	var imports []descriptor.GoPackage
	for _, f := range msg.Fields {
		if f.Ignore {
			continue
		}
		t, err := g.reg.LookupMsg("", f.GetTypeName())
		if err != nil {
			continue
		}
		pkg := t.File.GoPkg
		if pkg == file.GoPkg || pkgSeen[pkg.Path] {
			continue
		}
		pkgSeen[pkg.Path] = true
		imports = append(imports, pkg)
	}
	return imports
}
func (g *generator) addCrudPathParamImports(msg *descriptor.Message, pkgSeen map[string]bool) []descriptor.GoPackage {
	if !msg.GenerateCRUD {
		return []descriptor.GoPackage{}
	}
	var imports []descriptor.GoPackage

	if _, ok := msg.Implementations[crudOptions.Implementation_IMPLEMENTATION_SQLITE]; ok {
		if !pkgSeen["context"] {
			pkgSeen["context"] = true
			imports = append(imports, descriptor.GoPackage{Path: "context", Name: "context"})
		}
		if !pkgSeen["database/sql"] {
			pkgSeen["database/sql"] = true
			imports = append(imports, descriptor.GoPackage{Path: "database/sql", Name: "sql"})
		}
		if !pkgSeen["fmt"] {
			pkgSeen["fmt"] = true
			imports = append(imports, descriptor.GoPackage{Path: "fmt", Name: "fmt"})
		}
		if !pkgSeen["strings"] {
			pkgSeen["strings"] = true
			imports = append(imports, descriptor.GoPackage{Path: "strings", Name: "strings"})
		}
		if msg.HasFieldMask() && !pkgSeen["github.com/mennanov/fmutils"] {
			pkgSeen["github.com/mennanov/fmutils"] = true
			imports = append(imports, descriptor.GoPackage{Path: "github.com/mennanov/fmutils", Name: "fmutils"})
		}
		if msg.HasFieldMask() && !pkgSeen["google.golang.org/protobuf/types/known/fieldmaskpb"] {
			pkgSeen["google.golang.org/protobuf/types/known/fieldmaskpb"] = true
			imports = append(imports, descriptor.GoPackage{Path: "google.golang.org/protobuf/types/known/fieldmaskpb", Name: "fieldmaskpb"})
		}
		if !pkgSeen["modernc.org/sqlite"] {
			pkgSeen["modernc.org/sqlite"] = true
			imports = append(imports, descriptor.GoPackage{Path: "modernc.org/sqlite", Name: "sqlite"})
		}
		if !pkgSeen["github.com/samlitowitz/protoc-gen-crud/expressions"] {
			pkgSeen["github.com/samlitowitz/protoc-gen-crud/expressions"] = true
			imports = append(imports, descriptor.GoPackage{Path: "github.com/samlitowitz/protoc-gen-crud/expressions", Name: "expressions"})
		}
		if !pkgSeen["time"] {
			pkgSeen["time"] = true
			imports = append(imports, descriptor.GoPackage{Path: "time", Name: "time"})
		}
	}

	return imports
}
