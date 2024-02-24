package gen_sqlite

import (
	"bytes"
	"fmt"
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

	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		if !f.HasRelationship() {
			panic(fmt.Errorf("message type without relationship defined on field %s", f.GetName()))
		}
		minUIDFields := f.Relationship.CRUD.MinimalUIDFields()
		if len(minUIDFields) != 1 {
			panic(fmt.Errorf("message type must have unique identifier with exactly one field defined on field %s", f.GetName()))
		}
		return sqliteType(minUIDFields[0])

	case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
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
		"sqliteIdent":                   gen_go_crud.SQLiteIdent,
		"sqliteTableName":               gen_go_crud.SQLiteTableName,
		"sqliteColumnName":              gen_go_crud.SQLiteColumnName,
		"SQLiteColumnNameFromFieldName": gen_go_crud.SQLiteColumnNameFromFieldName,
		"sqliteType":                    sqliteType,
	}

	// https://www.sqlite.org/lang_createtable.html
	createTableTemplate = template.Must(template.New("create-table").Funcs(funcMap).Parse(`
DROP TABLE IF EXISTS {{sqliteIdent (sqliteTableName .CRUD.GetName)}};
CREATE TABLE IF NOT EXISTS {{sqliteIdent (sqliteTableName .CRUD.GetName)}} (
    {{ range $i, $field := .CRUD.DataFields -}}
    {{if $i}},
    {{end}}{{sqliteIdent (sqliteColumnName (SQLiteColumnNameFromFieldName $field))}} {{sqliteType $field}}
    {{- end }}
    {{- if gt (len .CRUD.MinimalUIDFields) 0 -}}
        ,

    PRIMARY KEY ({{ range $i, $field := .CRUD.MinimalUIDFields -}}
        {{if $i}},{{end}}{{sqliteIdent (sqliteColumnName $field.GetName)}}
        {{- end }})
    {{- end}}
);
`))
)
