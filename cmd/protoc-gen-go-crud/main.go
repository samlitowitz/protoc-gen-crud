package main

import (
	"flag"
	"fmt"
	gen_sqlite "github.com/samlitowitz/protoc-gen-crud/internal/gen-sqlite"
	"os"

	gen_gen "github.com/samlitowitz/protoc-gen-crud/internal/gen-gen"

	gen_go_crud "github.com/samlitowitz/protoc-gen-crud/internal/gen-go-crud"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"

	"google.golang.org/protobuf/compiler/protogen"
)

var (
	versionFlag = flag.Bool("version", false, "print protoc-gen-go-crud Version")
)

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Version %v, commit %v, built at %v\n", Version, Commit, Date)
		os.Exit(0)
	}

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		reg := descriptor.NewRegistry()

		crudGen := gen_go_crud.New(reg)
		sqliteGen := gen_sqlite.New(reg)

		genGen := gen_gen.New(crudGen, sqliteGen)

		if err := reg.LoadFromPlugin(gen); err != nil {
			return err
		}

		targets := make([]*descriptor.File, 0, len(gen.Request.FileToGenerate))
		for _, target := range gen.Request.FileToGenerate {
			f, err := reg.LookupFile(target)
			if err != nil {
				return err
			}
			targets = append(targets, f)
		}

		files, err := genGen.Generate(targets)
		for _, f := range files {
			genFile := gen.NewGeneratedFile(f.GetName(), protogen.GoImportPath(f.GoPkg.Path))
			if _, err := genFile.Write([]byte(f.GetContent())); err != nil {
				return err
			}
		}

		return err
	})
}
