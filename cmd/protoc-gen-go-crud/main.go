package main

import (
	"flag"
	"fmt"
	"os"

	genPgSQLCRUD "github.com/samlitowitz/protoc-gen-crud/internal/generator/pgsql/crud"
	genPgSQLSQL "github.com/samlitowitz/protoc-gen-crud/internal/generator/pgsql/sql"
	genGoRelationship "github.com/samlitowitz/protoc-gen-crud/internal/generator/relationship"
	genSQLiteCRUD "github.com/samlitowitz/protoc-gen-crud/internal/generator/sqlite/crud"
	genSQLiteSQL "github.com/samlitowitz/protoc-gen-crud/internal/generator/sqlite/sql"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
	genGoCRUD "github.com/samlitowitz/protoc-gen-crud/internal/generator/crud"
	genGen "github.com/samlitowitz/protoc-gen-crud/internal/generator/generator"
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

		crudGen := genGoCRUD.New(reg, genGoCRUD.WithFormatOutput(*formatOutput))
		relationshipGen := genGoRelationship.New(reg)
		pgsqlCRUDGen := genPgSQLCRUD.New(reg)
		pgsqlSQLGen := genPgSQLSQL.New(reg)
		sqliteCRUDGen := genSQLiteCRUD.New(reg)
		sqliteSQLGen := genSQLiteSQL.New(reg)

		gg := genGen.New(crudGen, relationshipGen, pgsqlCRUDGen, pgsqlSQLGen, sqliteCRUDGen, sqliteSQLGen)

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

		files, err := gg.Generate(targets)
		for _, f := range files {
			genFile := gen.NewGeneratedFile(f.GetName(), protogen.GoImportPath(f.GoPkg.Path))
			if _, err := genFile.Write([]byte(f.GetContent())); err != nil {
				return err
			}
		}
		return err
	})
}
