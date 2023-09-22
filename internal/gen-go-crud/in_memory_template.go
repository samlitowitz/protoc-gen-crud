package gen_go_crud

import (
	"fmt"
	"text/template"

	"github.com/iancoleman/strcase"

	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
)

func init() {
	strcase.ConfigureAcronym("UID", "uid")
}

type inMemory struct {
	uidKeyTypeByUIDNames map[string]string
}

func (im *inMemory) UIDKeyTypeByUIDNames(crud *descriptor.CRUD) map[string]string {
	if im.uidKeyTypeByUIDNames != nil {
		return im.uidKeyTypeByUIDNames
	}
	im.uidKeyTypeByUIDNames = make(map[string]string)
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
		im.uidKeyTypeByUIDNames[name] = typ
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
	{{range $name, $typ := .InMemory.UIDKeyTypeByUIDNames .CRUD}}
	{{$name}} map[{{$typ}}]*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}}
	{{end}}
	iTable []*{{.CRUD.GoType .CRUD.File.GoPkg.Path}} // Internal table of all {{.CRUD.Name}}s
}

// NewInMemory creates a new InMemory{{.CRUD.Name}}Repository to be used.
func NewInMemory{{.CRUD.Name}}Repository() *InMemory{{.CRUD.Name}}Repository {
	return &InMemory{{.CRUD.Name}}Repository{
		{{range $name, $typ := .InMemory.UIDKeyTypeByUIDNames .CRUD}}{{$name}}: make(map[{{$typ}}]*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}}),{{end}}
		iTable: make([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, 0),
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
	panic("not implemented")
}
{{end}}
`))
)
