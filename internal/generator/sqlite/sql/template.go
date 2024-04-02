package sql

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/samlitowitz/protoc-gen-crud/options"

	"github.com/samlitowitz/protoc-gen-crud/internal/generator/sqlite"

	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
)

type param struct {
	*descriptor.File
}

type message struct {
	*descriptor.Message
}

type enum struct {
	*descriptor.Enum
}

func fieldColType(f *descriptor.Field) string {
	switch f.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return "REAL"

	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
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

	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		// Enums will reference a look-up table
		return "INTEGER"

	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		fallthrough
	default:
		panic(fmt.Errorf("sqlite: sql: field %s: unsupported type %s", f.GetName(), f.GetType()))
	}
}

func fieldColComment(f *descriptor.Field) string {
	switch f.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		fallthrough
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
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return ""

	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return fmt.Sprintf(
			" /* references %s.%s */",
			sqlite.SQLiteQuotedIdent(f.FieldEnum.GetName()),
			sqlite.SQLiteQuotedIdent("id"),
		)

	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		fallthrough
	default:
		panic(fmt.Errorf("sqlite: sql: field %s: unsupported type %s", f.GetName(), f.GetType()))
	}
}

func applyTemplate(p param, reg *descriptor.Registry) (string, error) {
	completedEnums := make(map[string]struct{})

	w := bytes.NewBuffer(nil)

	//return "", fmt.Errorf("%v", p.Messages)

	for _, msg := range p.Messages {
		if !msg.GenerateCRUD {
			continue
		}
		if _, ok := msg.Implementations[options.Implementation_SQLITE]; !ok {
			continue
		}

		for _, field := range msg.Fields {
			if field.GetType() != descriptorpb.FieldDescriptorProto_TYPE_ENUM {
				continue
			}
			if field.FieldEnum == nil {
				return "", fmt.Errorf("%s: missing enum definition", field.GetName())
			}
			if _, ok := completedEnums[field.FieldEnum.FQEN()]; ok {
				continue
			}

			if err := createTableForEnumTemplate.Execute(w, &enum{field.FieldEnum}); err != nil {
				return "", fmt.Errorf("%s: %s: create enum table: %v", field.GetName(), field.FieldEnum.GetName(), err)
			}
		}

		if err := createTableForMessageTemplate.Execute(w, &message{msg}); err != nil {
			return "", fmt.Errorf("%s: create message table: %v", msg.GetName(), err)
		}
	}
	return w.String(), nil
}

var (
	funcMap template.FuncMap = map[string]interface{}{
		"quotedIdent":     sqlite.SQLiteQuotedIdent,
		"fieldColComment": fieldColComment,
		"fieldColType":    fieldColType,
	}

	// https://www.sqlite.org/lang_createtable.html
	createTableForMessageTemplate = template.Must(template.New("create-table-for-message").Funcs(funcMap).Parse(`
DROP TABLE IF EXISTS {{quotedIdent .GetName}};
CREATE TABLE IF NOT EXISTS {{quotedIdent .GetName}} (
{{- range $i, $field := .PrimaryKey -}}
    {{- if $i}},{{end}}
    {{quotedIdent $field.GetName}} {{fieldColType $field}}{{fieldColComment $field}}
{{- end -}}
{{- if gt (len .NonPrimeAttributes) 0 -}},{{- end -}}

{{- range $i, $field := .NonPrimeAttributes -}}
    {{- if $i}},{{end}}
    {{quotedIdent $field.GetName}} {{fieldColType $field}}{{fieldColComment $field}}
{{- end }}
{{- if gt (len .PrimaryKey) 0 -}}
        ,

    PRIMARY KEY (
    {{- range $i, $field := .PrimaryKey -}}
        {{- if $i}},{{end}}
        {{quotedIdent $field.GetName}}
    {{- end}}
    )
    {{- end}}
);
`))

	createTableForEnumTemplate = template.Must(template.New("create-table-for-enum").Funcs(funcMap).Parse(`
DROP TABLE IF EXISTS {{quotedIdent .GetName}};
CREATE TABLE IF NOT EXISTS {{quotedIdent .GetName}} (
    "id" INTEGER PRIMARY KEY,
    "value" TEXT
);

INSERT INTO {{quotedIdent .GetName}} ("id", "value") VALUES
{{- range $i, $valDesc := .GetValue}}
    {{- if $i}},{{end}}
    ({{$valDesc.GetNumber}}, "{{$valDesc.GetName}}")
{{- end}}
;
`))
)
