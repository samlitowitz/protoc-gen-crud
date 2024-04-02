package main

import (
	"flag"
	"fmt"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
	gen_go_crud "github.com/samlitowitz/protoc-gen-crud/internal/generator/crud"
	gen_gen "github.com/samlitowitz/protoc-gen-crud/internal/generator/generator"
	gen_go_relationship "github.com/samlitowitz/protoc-gen-crud/internal/generator/relationship"
	gen_sqlite_crud "github.com/samlitowitz/protoc-gen-crud/internal/generator/sqlite/crud"
	gen_sqlite_sql "github.com/samlitowitz/protoc-gen-crud/internal/generator/sqlite/sql"
)

var (
	formatOutput = flag.Bool("format_output", true, "format code before writing to file")
	versionFlag  = flag.Bool("version", false, "print protoc-gen-go-crud Version")
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

		crudGen := gen_go_crud.New(reg, gen_go_crud.WithFormatOutput(*formatOutput))
		relationshipGen := gen_go_relationship.New(reg)
		sqliteCRUDGen := gen_sqlite_crud.New(reg)
		sqliteSQLGen := gen_sqlite_sql.New(reg)

		genGen := gen_gen.New(crudGen, relationshipGen, sqliteCRUDGen, sqliteSQLGen)

		if err := reg.LoadFromPlugin(gen); err != nil {
			return err
		}

		targets := make([]*descriptor.File, 0, len(gen.Request.FileToGenerate))
		for _, target := range gen.Request.FileToGenerate {
			f, err := reg.LookupFile(target)
			if err != nil {
				return err
			}
			if f.FileDescriptorProto.GetSyntax() != protoreflect.Proto3.String() {
				return fmt.Errorf(
					"%s: unsupported syntax %s, must be %s",
					target,
					f.FileDescriptorProto.GetSyntax(),
					protoreflect.Proto3,
				)
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
