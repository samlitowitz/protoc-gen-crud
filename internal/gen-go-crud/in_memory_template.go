package gen_go_crud

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/samlitowitz/protoc-gen-crud/internal/casing"

	"github.com/iancoleman/strcase"

	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
)

func init() {
	strcase.ConfigureAcronym("UID", "uid")
}

type inMemoryUIDData struct {
	KeyType string
	fields  []*descriptor.Field
}

func (uidData *inMemoryUIDData) KeyGenerationCode(fieldDefVarName string, keyVarName string) string {
	if len(uidData.fields) < 1 {
		return ""
	}
	if len(uidData.fields) == 1 {
		def := uidData.fields[0]
		switch *def.Type {
		case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
			return ""

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
		case descriptorpb.FieldDescriptorProto_TYPE_STRING:
			return fmt.Sprintf("%s := %s.%s", keyVarName, fieldDefVarName, casing.CamelIdentifier(*def.Name))

		case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
			return fmt.Sprintf("%s := string(%s.%s)", keyVarName, fieldDefVarName, *def.Name)
		}
	}
	buf := bytes.Buffer{}
	buf.WriteString("buf := bytes.Buffer{}")
	for _, def := range uidData.fields {
		switch *def.Type {
		case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
			continue

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
		case descriptorpb.FieldDescriptorProto_TYPE_STRING:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
			buf.WriteString(
				fmt.Sprintf(`
err := binary.Write(buf, binary.LittleEndian, %s.%s)
if err != nil {
	return nil, err
}
`,
					fieldDefVarName,
					*def.Name,
				),
			)
		}
	}
	buf.WriteString(
		fmt.Sprintf("%s := buf.String()", keyVarName),
	)
	return buf.String()
}

type inMemory struct {
	uidKeyTypeByUIDNames map[string]*inMemoryUIDData
}

// For each uid we'll need a function to take an input of []*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}} and return a map[{{$typ}}]*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}}

func (im *inMemory) UIDDataByUIDNames(crud *descriptor.CRUD) map[string]*inMemoryUIDData {
	if im.uidKeyTypeByUIDNames != nil {
		return im.uidKeyTypeByUIDNames
	}
	im.uidKeyTypeByUIDNames = make(map[string]*inMemoryUIDData)
	for name, fieldDefs := range crud.UniqueIdentifiers {
		typ, err := im.uidKeyTypeByFieldDefs(fieldDefs)
		if err != nil {
			panic(err)
		}
		name = fmt.Sprintf(
			"%sBy%s",
			crud.CamelCaseName(),
			strcase.ToCamel(name),
		)

		im.uidKeyTypeByUIDNames[name] = &inMemoryUIDData{
			KeyType: typ,
			fields:  fieldDefs,
		}
	}
	return im.uidKeyTypeByUIDNames
}

func (im inMemory) uidKeyTypeByFieldDefs(defs []*descriptor.Field) (string, error) {
	if len(defs) < 1 {
		return "", fmt.Errorf("no fields found, cannot build key type")
	}
	if len(defs) == 1 {
		def := defs[0]
		switch *def.Type {
		case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
			return "", &descriptor.UnsupportedTypeError{TypName: *def.TypeName}

		case descriptorpb.FieldDescriptorProto_TYPE_UINT32:
			return "uint32", nil
		case descriptorpb.FieldDescriptorProto_TYPE_UINT64:
			return "uint64", nil

		case descriptorpb.FieldDescriptorProto_TYPE_INT32:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_SINT32:
			return "int32", nil

		case descriptorpb.FieldDescriptorProto_TYPE_INT64:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_SINT64:
			return "int64", nil

		case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
			return "string", nil

		case descriptorpb.FieldDescriptorProto_TYPE_STRING:
			return "string", nil
		}
	}
	return "string", nil
}

var (
	_ = template.Must(repositoryTemplate.New("repository-in-memory").Funcs(funcMap).Parse(`

// InMemory{{.CRUD.Name}}Repository is an in memory implementation of the {{.CRUD.Name}}Repository interface.
type InMemory{{.CRUD.Name}}Repository struct {
	{{range $name, $data := .InMemory.UIDDataByUIDNames .CRUD}}
	{{$name}} map[{{$data.KeyType}}]*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}}
	{{end}}
}

// NewInMemory creates a new InMemory{{.CRUD.Name}}Repository to be used.
func NewInMemory{{.CRUD.Name}}Repository() *InMemory{{.CRUD.Name}}Repository {
	return &InMemory{{.CRUD.Name}}Repository{
		{{range $name, $data := .InMemory.UIDDataByUIDNames .CRUD}}{{$name}}: make(map[{{$data.KeyType}}]*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}}),{{end}}
	}
}

{{if .CRUD.Create}}
// Create creates new {{.CRUD.Name}}s.
// Successfully created {{.CRUD.Name}}s are returned along with any errors that may have occurred.
func (repo *InMemory{{.CRUD.Name}}Repository) Create([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}) ([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, error) {
	panic("not implemented")
}
{{end}}

{{if .CRUD.Read}}
// Read returns a set of {{.CRUD.Name}}s matching the provided criteria
// Read is incomplete and it should be considered unstable
// Use where clauses
func (repo *InMemory{{.CRUD.Name}}Repository) Read() ([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, error) {
	panic("not implemented")
}
{{end}}

{{if .CRUD.Update}}
// Update modifies existing {{.CRUD.Name}}s based on the defined unique identifiers.
func (repo *InMemory{{.CRUD.Name}}Repository) Update([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}) ([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, error) {
	panic("not implemented")
}
{{end}}

{{if .CRUD.Delete}}
// Delete deletes {{.CRUD.Name}}s based on the defined unique identifiers
// Delete is incomplete and it should be considered unstable
// Use where clauses
func (repo *InMemory{{.CRUD.Name}}Repository) Delete([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}) error {
	// TODO: Get structs by uid(s)
	// TODO: Remove found structs
	// TODO: Return error(s)
	panic("not implemented")
}
{{end}}

{{range $name, $data := .InMemory.UIDDataByUIDNames .CRUD}}
	func get{{camelIdentifier $name}}({{$.CRUD.CamelCaseName}}s []*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}}) (map[{{.KeyType}}]*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}}, error) {
	{{$name}} := make(map[{{.KeyType}}]*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}})
	for _, def := range {{$.CRUD.CamelCaseName}}s {
		{{$data.KeyGenerationCode "def" "key"}}
		{{$name}}[key] = def
	}
	return {{$name}}, nil
}
{{end}}
`))
)
