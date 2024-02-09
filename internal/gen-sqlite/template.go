package gen_sqlite

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"google.golang.org/protobuf/types/descriptorpb"

	gen_go_crud "github.com/samlitowitz/protoc-gen-crud/internal/gen-go-crud"

	"github.com/samlitowitz/protoc-gen-crud/internal/casing"
	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
)

type param struct {
	*descriptor.File
}

type crud struct {
	*descriptor.CRUD
}

func colDefs(crud *descriptor.CRUD) string {
	var colDefs []string

	for _, field := range crud.Fields {
		colDefs = append(
			colDefs,
			fmt.Sprintf(
				"%s %s",
				gen_go_crud.SQLiteColumnIdentifier(field.GetName()),
				sqliteType(field),
			),
		)
	}

	return strings.Join(colDefs, ",\n\t")
}

func sqliteType(f *descriptor.Field) string {
	switch *f.Type {
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return "REAL"

	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return "INTEGER"
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return "INTEGER"

	case descriptorpb.FieldDescriptorProto_TYPE_UINT32:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_INT32:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_SINT32:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_INT64:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		return "INTEGER"

	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return "BLOB"

	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return "TEXT"

	case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		fallthrough
	default:
		panic(fmt.Errorf("non-scalar type on field %s", f.GetName()))
	}
}

func applyTemplate(p param, reg *descriptor.Registry) (string, error) {
	w := bytes.NewBuffer(nil)

	for _, msg := range p.Messages {
		msgName := casing.Camel(*msg.Name)
		msg.Name = &msgName
	}
	for _, def := range p.CRUDs {
		if !def.SQLiteImplementation() {
			continue
		}
		inject := &crud{
			CRUD: def,
		}

		if err := createTableTemplate.Execute(w, inject); err != nil {
			return "", err
		}
	}

	return w.String(), nil
}

var (
	funcMap template.FuncMap = map[string]interface{}{
		"tableName": gen_go_crud.SQLiteTableName,
		"colDefs":   colDefs,
	}

	// https://www.sqlite.org/lang_createtable.html
	createTableTemplate = template.Must(template.New("create-table").Funcs(funcMap).Parse(`
CREATE TABLE IF NOT EXISTS "{{tableName .CRUD.GetName}}" (
	{{colDefs .CRUD}}
	-- TODO: column definitions with constraints
);
-- TODO: table constraints, unique and pkey
`))
)
