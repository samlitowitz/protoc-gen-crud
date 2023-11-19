// REFURL: https://github.com/grpc-ecosystem/grpc-gateway/blob/main/protoc-gen-grpc-gateway/internal/gengateway/generator.go
package gen_go_crud

import (
	"fmt"
	"go/format"
	"path"

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
	var imports []descriptor.GoPackage
	for _, pkgpath := range []string{
		//"io",
		//"google.golang.org/protobuf/proto",
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

	for _, msg := range file.Messages {
		imports = append(imports, g.addMessagePathParamImports(file, msg, pkgSeen)...)
	}
	for _, crud := range file.CRUDs {
		imports = append(imports, g.addCrudPathParamImports(file, crud, pkgSeen)...)
	}

	params := param{
		File:    file,
		Imports: imports,
	}

	return applyTemplate(params, g.reg)
}

// addMessagePathParamImports handles adding import of message path parameter go packages
func (g *generator) addMessagePathParamImports(file *descriptor.File, m *descriptor.Message, pkgSeen map[string]bool) []descriptor.GoPackage {
	var imports []descriptor.GoPackage
	for _, f := range m.Fields {
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

func (g *generator) addCrudPathParamImports(file *descriptor.File, crud *descriptor.CRUD, pkgSeen map[string]bool) []descriptor.GoPackage {
	var imports []descriptor.GoPackage
	for _, fields := range crud.UniqueIdentifiers {
		if len(fields) <= 1 {
			continue
		}
		if !pkgSeen["encoding/binary"] {
			pkgSeen["encoding/binary"] = true
			imports = append(imports, descriptor.GoPackage{Path: "encoding/binary", Name: "binary"})
		}
		if !pkgSeen["crypto/sha256"] {
			pkgSeen["crypto/sha256"] = true
			imports = append(imports, descriptor.GoPackage{Path: "crypto/sha256", Name: "sha256"})
		}
	}
	return imports
}
