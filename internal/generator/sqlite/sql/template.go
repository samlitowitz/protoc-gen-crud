package sql

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/samlitowitz/protoc-gen-crud/internal/generator/crud"

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

	PrimaryKeyCols        []*sqlite.Column
	NonPrimeAttributeCols []*sqlite.Column
}

type enum struct {
	*descriptor.Enum
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

		injected := &message{
			Message:               msg,
			PrimaryKeyCols:        sqlite.ColumnsFromFields(crud.QueryableFieldsFromFields(msg.PrimaryKey())),
			NonPrimeAttributeCols: sqlite.ColumnsFromFields(crud.QueryableFieldsFromFields(msg.NonPrimeAttributes())),
		}
		if err := createTableForMessageTemplate.Execute(w, injected); err != nil {
			return "", fmt.Errorf("%s: create message table: %v", msg.GetName(), err)
		}
	}
	return w.String(), nil
}

var (
	funcMap template.FuncMap = map[string]interface{}{
		"quotedIdent": sqlite.QuotedIdent,
	}

	// https://www.sqlite.org/lang_createtable.html
	createTableForMessageTemplate = template.Must(template.New("create-table-for-message").Funcs(funcMap).Parse(`
DROP TABLE IF EXISTS {{quotedIdent .GetName}};
CREATE TABLE IF NOT EXISTS {{quotedIdent .GetName}} (
{{- range $i, $col := .PrimaryKeyCols -}}
    {{- if $i}},{{end}}
    {{- template "column-definition" $col}}
{{- end -}}
{{- if gt (len .NonPrimeAttributeCols) 0 -}},{{- end -}}

{{- range $i, $col := .NonPrimeAttributeCols -}}
    {{- if $i}},{{end}}
    {{template "column-definition" $col}}
{{- end }}
{{- if gt (len .PrimaryKey) 0 -}}
        ,

    PRIMARY KEY (
    {{- range $i, $col := .PrimaryKeyCols -}}
        {{- if $i}},{{end}}
        {{quotedIdent $col.GetName}}
    {{- end}}
    )
    {{- end}}
);
`))

	_ = template.Must(createTableForMessageTemplate.New("column-definition").Funcs(funcMap).Parse(`
    {{- quotedIdent .GetName}} {{.GetType}}{{.GetComment -}}
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
