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

func inMemoryBuildMapWrap(crud *descriptor.CRUD, name string, data *inMemoryUIDData) map[string]interface{} {
	return map[string]interface{}{
		"CRUD": crud,
		"Name": name,
		"Data": data,
	}
}

type inMemoryUIDData struct {
	KeyType        string
	KeyIsComposite bool

	keyGenerationCode string
}

func (uidData *inMemoryUIDData) SimpleKeyGenerationCode(fieldDefVarName string) string {
	if uidData.KeyIsComposite {
		return ""
	}
	return fmt.Sprintf(uidData.keyGenerationCode, fieldDefVarName)
}

func (uidData *inMemoryUIDData) CompositeKeyGenerationCode(
	fieldDefVarName string,
	hashVarName string,
) string {
	if !uidData.KeyIsComposite {
		return ""
	}
	return fmt.Sprintf(uidData.keyGenerationCode, hashVarName, fieldDefVarName)
}

type inMemory struct {
	uidDataByUIDNames map[string]*inMemoryUIDData
}

// For each uid we'll need a function to take an input of []*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}} and return a map[{{$typ}}]*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}}
func (im *inMemory) UIDDataByUIDNames(crud *descriptor.CRUD) map[string]*inMemoryUIDData {
	if im.uidDataByUIDNames != nil {
		return im.uidDataByUIDNames
	}
	im.uidDataByUIDNames = make(map[string]*inMemoryUIDData)
	for name, fieldDefs := range crud.UniqueIdentifiers {
		uidData, err := im.buildUIDData(fieldDefs)
		if err != nil {
			panic(err)
		}
		name = fmt.Sprintf(
			"%sBy%s",
			crud.CamelCaseName(),
			strcase.ToCamel(name),
		)

		im.uidDataByUIDNames[name] = uidData
	}
	return im.uidDataByUIDNames
}

func (im inMemory) buildUIDData(defs []*descriptor.Field) (*inMemoryUIDData, error) {
	if len(defs) < 1 {
		return nil, fmt.Errorf("no fields found, cannot build key type")
	}
	uidData := &inMemoryUIDData{
		KeyIsComposite: len(defs) > 1,
		KeyType:        "string",
	}

	keyGenCodeBuf := bytes.Buffer{}
	for _, def := range defs {
		uidData.keyGenerationCode = fmt.Sprintf("%%s.%s", casing.CamelIdentifier(*def.Name))

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
			return nil, &descriptor.UnsupportedTypeError{TypName: *def.TypeName}

		case descriptorpb.FieldDescriptorProto_TYPE_UINT32:
			if !uidData.KeyIsComposite {
				uidData.KeyType = "uint32"
				return uidData, nil
			}

		case descriptorpb.FieldDescriptorProto_TYPE_UINT64:
			if !uidData.KeyIsComposite {
				uidData.KeyType = "uint64"
				return uidData, nil
			}

		case descriptorpb.FieldDescriptorProto_TYPE_INT32:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_SINT32:
			if !uidData.KeyIsComposite {
				uidData.KeyType = "int32"
				return uidData, nil
			}

		case descriptorpb.FieldDescriptorProto_TYPE_INT64:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_SINT64:
			if !uidData.KeyIsComposite {
				uidData.KeyType = "int64"
				return uidData, nil
			}

		case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
			if !uidData.KeyIsComposite {
				uidData.keyGenerationCode = fmt.Sprintf("string(%%s.%s)", casing.CamelIdentifier(*def.Name))
				return uidData, nil
			}

		case descriptorpb.FieldDescriptorProto_TYPE_STRING:
			if !uidData.KeyIsComposite {
				return uidData, nil
			}
		}
		if uidData.KeyIsComposite {
			keyGenCodeBuf.WriteString(
				fmt.Sprintf(
					`
err = binary.Write(%%[1]s, binary.LittleEndian, "{{")
if err != nil {
	return nil, err
}
err = binary.Write(%%[1]s, binary.LittleEndian, %%[2]s.%s)
if err != nil {
	return nil, err
}
err = binary.Write(%%[1]s, binary.LittleEndian, "}}")
if err != nil {
	return nil, err
}
`,
					casing.CamelIdentifier(*def.Name),
				),
			)
			uidData.keyGenerationCode = keyGenCodeBuf.String()
		}
	}
	return uidData, nil
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
func build{{camelIdentifier $name}}Map({{$.CRUD.CamelCaseName}}s []*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}}) (map[{{.KeyType}}]*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}}, error) {
	{{$name}} := make(map[{{.KeyType}}]*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}})
	{{if $data.KeyIsComposite}}
	{{- template "repository-in-memory-build-map-for-complex-keys" (inMemoryBuildMapWrap $.CRUD $name $data) -}}
	{{else}}
	{{- template "repository-in-memory-build-map-for-simple-keys" (inMemoryBuildMapWrap $.CRUD $name $data) -}}
	{{end}}
	return {{$name}}, nil
}
{{end}}
`))

	_ = template.Must(repositoryTemplate.New("repository-in-memory-build-map-for-simple-keys").Funcs(funcMap).Parse(`
	for _, def := range {{.CRUD.CamelCaseName}}s {
		{{.Name}}[{{.Data.SimpleKeyGenerationCode "def"}}] = def
	}
`))
	_ = template.Must(repositoryTemplate.New("repository-in-memory-build-map-for-complex-keys").Funcs(funcMap).Parse(`
	var err error
	h := sha256.New()
	for _, def := range {{$.CRUD.CamelCaseName}}s {
		{{.Data.CompositeKeyGenerationCode "def" "h"}}
		key := string(h.Sum(nil))
		{{.Name}}[key] = def
		h.Reset()
	}
`))
)
