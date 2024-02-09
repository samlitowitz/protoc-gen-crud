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
	KeyType           string
	KeyIsComposite    bool
	BuildMapFnName    string
	IndexByKeyMapName string
	KeyByIndexMapName string

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
		uidMapName := fmt.Sprintf(
			"%sBy%s",
			crud.CamelCaseName(),
			strcase.ToCamel(name),
		)
		uidData.BuildMapFnName = fmt.Sprintf("build%s%sMap", strcase.ToCamel(crud.GetName()), strcase.ToCamel(name))
		uidData.IndexByKeyMapName = fmt.Sprintf("indexBy%s", strcase.ToCamel(name))
		uidData.KeyByIndexMapName = fmt.Sprintf("%sByIndex", strcase.ToLowerCamel(name))
		im.uidDataByUIDNames[uidMapName] = uidData
	}
	return im.uidDataByUIDNames
}

func (im inMemory) buildUIDData(defs []*descriptor.Field) (*inMemoryUIDData, error) {
	if len(defs) < 1 {
		return nil, fmt.Errorf("no fields found, cannot build key type")
	}
	uidData := &inMemoryUIDData{
		KeyType:        "string",
		KeyIsComposite: len(defs) > 1,
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
					`_, err = %%[1]s.Write([]byte("{{"+%%[2]s.%s+"}}"))
if err != nil {
	return nil, nil, nil, err
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
	{{- range $name, $data := .InMemory.UIDDataByUIDNames .CRUD}}
	{{$name}} map[{{$data.KeyType}}]uint
	{{- end}}
	iTable map[uint]*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}}
	next uint
}

// NewInMemory creates a new InMemory{{.CRUD.Name}}Repository to be used.
func NewInMemory{{.CRUD.Name}}Repository() *InMemory{{.CRUD.Name}}Repository {
	return &InMemory{{.CRUD.Name}}Repository{
		{{- range $name, $data := .InMemory.UIDDataByUIDNames .CRUD}}
		{{$name}}: make(map[{{$data.KeyType}}]uint),
		{{- end}}
		iTable: make(map[uint]*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}}),
	}
}

{{if .CRUD.Create}}
// Create creates new {{.CRUD.Name}}s.
// Successfully created {{.CRUD.Name}}s are returned along with any errors that may have occurred.
func (repo *InMemory{{.CRUD.Name}}Repository) Create(ctx context.Context, toCreate []*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}) ([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, error) {
	indicesToCreate := make(map[uint]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}})
	{{range $name, $data := .InMemory.UIDDataByUIDNames .CRUD}}
	{{$data.IndexByKeyMapName}}, {{$data.KeyByIndexMapName}}, {{$name}}, err := {{$data.BuildMapFnName}}(toCreate)
	if err != nil {
		return nil, err
	}
	for key, val := range {{$name}} {
		if _, ok := {{$data.IndexByKeyMapName}}[key]; !ok {
			// internal error, should never happen
			continue
		}
		if _, ok := {{$data.KeyByIndexMapName}}[{{$data.IndexByKeyMapName}}[key]]; !ok {
			// internal error, should never happen
			continue

		}
		if _, ok := repo.{{$name}}[key]; ok {
			// add error about duplicate
			delete(indicesToCreate, {{$data.IndexByKeyMapName}}[key])
			continue
		}
		if _, ok := indicesToCreate[{{$data.IndexByKeyMapName}}[key]]; !ok {
			// mark index as to be created
			indicesToCreate[{{$data.IndexByKeyMapName}}[key]] = val
		}
	}
	{{end}}

	created := make([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, 0, len(indicesToCreate))
	for i, val := range indicesToCreate {
		repo.iTable[repo.next] = val
		{{- range $name, $data := .InMemory.UIDDataByUIDNames .CRUD}}
		repo.{{$name}}[{{$data.KeyByIndexMapName}}[i]] = repo.next
		{{- end}}
		created = append(created, val)
		repo.next++
	}
	return created, nil
}
{{end}}

{{if .CRUD.Read}}
// Read returns a set of {{.CRUD.Name}}s matching the provided criteria
// Read is incomplete and it should be considered unstable
func (repo *InMemory{{.CRUD.Name}}Repository) Read(ctx context.Context, expr expressions.Expression) ([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, error) {
	found := make([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, 0, len(repo.iTable))

	for _, {{.CRUD.CamelCaseName}} := range repo.iTable {
		include, err := testExprOn{{.CRUD.Name}}({{.CRUD.CamelCaseName}}, expr)
		if err != nil {
			return nil, err
		}
		if !include {
			continue
		}
		found = append(found, {{.CRUD.CamelCaseName}})
	}
	return found, nil
}
{{end}}

{{if .CRUD.Update}}
// Update modifies existing {{.CRUD.Name}}s based on the defined unique identifiers.
func (repo *InMemory{{.CRUD.Name}}Repository) Update(ctx context.Context, toUpdate []*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}) ([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, error) {
	indicesToUpdate := make(map[uint]uint)
	{{range $name, $data := .InMemory.UIDDataByUIDNames .CRUD}}
	{{$data.IndexByKeyMapName}}, {{$data.KeyByIndexMapName}}, {{$name}}, err := {{$data.BuildMapFnName}}(toUpdate)
	if err != nil {
		return nil, err
	}
	for key, _ := range {{$name}} {
		if _, ok := {{$data.IndexByKeyMapName}}[key]; !ok {
			// internal error, should never happen
			continue
		}
		if _, ok := {{$data.KeyByIndexMapName}}[{{$data.IndexByKeyMapName}}[key]]; !ok {
			// internal error, should never happen
			continue

		}
		if _, ok := repo.{{$name}}[key]; !ok {
			// add error about missing
			delete(indicesToUpdate, {{$data.IndexByKeyMapName}}[key])
			continue
		}
		if _, ok := indicesToUpdate[{{$data.IndexByKeyMapName}}[key]]; ok {
			continue
		}
		// mark index as to be created
		indicesToUpdate[{{$data.IndexByKeyMapName}}[key]] = repo.{{$name}}[{{$data.KeyByIndexMapName}}[{{$data.IndexByKeyMapName}}[key]]]
	}
	{{end}}

	updated := make([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, 0, len(indicesToUpdate))
	for i, j := range indicesToUpdate {
		repo.iTable[j] = toUpdate[i]
		updated = append(updated, toUpdate[i])
	}
	return updated, nil
}
{{end}}

{{if .CRUD.Delete}}
// Delete deletes {{.CRUD.Name}}s matching the provided criteria
// Delete is incomplete and it should be considered unstable
func (repo *InMemory{{.CRUD.Name}}Repository) Delete(ctx context.Context, expr expressions.Expression) error {
	indicesToDelete := make(map[uint]struct{})
	toDelete := []*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}{}

	for i, {{.CRUD.CamelCaseName}} := range repo.iTable {
		include, err := testExprOn{{.CRUD.Name}}({{.CRUD.CamelCaseName}}, expr)
		if err != nil {
			return err
		}
		if !include {
			continue
		}
		indicesToDelete[i] = struct{}{}
		toDelete = append(toDelete, {{.CRUD.CamelCaseName}})
	}

	{{range $name, $data := .InMemory.UIDDataByUIDNames .CRUD}}
	_, {{$data.KeyByIndexMapName}}, _, err := {{$data.BuildMapFnName}}(toDelete)
	if err != nil {
		return err
	}
	{{end}}

	for i, _ := range indicesToDelete {
		{{- range $name, $data := .InMemory.UIDDataByUIDNames .CRUD}}
		// remove iTable entry indexed by {{$name}}
		delete(repo.iTable, repo.{{$name}}[{{$data.KeyByIndexMapName}}[i]])
		// remove {hash, iTable index} from {{$name}}
		delete(repo.{{$name}}, {{$data.KeyByIndexMapName}}[i])
		{{- end}}
		delete(repo.iTable, i)
	}
	return nil
}
{{end}}
{{if or .CRUD.Read .CRUD.Delete}}
func testExprOn{{.CRUD.Name}}({{.CRUD.CamelCaseName}} *{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, expr expressions.Expression) (bool, error) {
	if {{.CRUD.CamelCaseName}} == nil {
		return false, nil
	}
	if expr == nil {
		return true, nil
	}
	switch expr := expr.(type) {
		case *expressions.And:
			left, err := testExprOn{{.CRUD.Name}}({{.CRUD.CamelCaseName}}, expr.Left())
			if err != nil {
				return true, err
			}
			right, err := testExprOn{{.CRUD.Name}}({{.CRUD.CamelCaseName}}, expr.Right())
			if err != nil {
				return true, err
			}
			return left && right, nil

		case *expressions.Or:
			left, err := testExprOn{{.CRUD.Name}}({{.CRUD.CamelCaseName}}, expr.Left())
			if err != nil {
				return true, err
			}
			right, err := testExprOn{{.CRUD.Name}}({{.CRUD.CamelCaseName}}, expr.Right())
			if err != nil {
				return true, err
			}
			return left || right, nil
		case *expressions.Not:
			operand, err := testExprOn{{.CRUD.Name}}({{.CRUD.CamelCaseName}}, expr.Operand())
			if err != nil {
				return true, err
			}
			return !operand, nil

		case *expressions.Equal:
			result, err := testEqualExprOn{{.CRUD.Name}}({{.CRUD.CamelCaseName}}, expr)
			if err != nil {
				return true, err
			}
			return result, nil

		case *expressions.Identifier:
			return true, fmt.Errorf("identifiers and scalar values are only supported as operands of the equal expression")
		case *expressions.Scalar:
			return true, fmt.Errorf("identifiers and scalar values are only supported as operands of the equal expression")
		default:
			return true, fmt.Errorf("unknown expression")
	}
}

func testEqualExprOn{{.CRUD.Name}}({{.CRUD.CamelCaseName}} *{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, expr *expressions.Equal) (bool, error) {
	{{if .FieldByFieldConstants}}var ident *expressions.Identifier
	var scalar *expressions.Scalar

	left := expr.Left()
	switch left := left.(type) {
	case *expressions.Identifier:
		if _, ok := valid{{.CRUD.Name}}Fields[left.ID()]; !ok {
			return true, fmt.Errorf("invalid field id: %s", left.ID())
		}
		ident = left
	case *expressions.Scalar:
		scalar = left
	default:
		return true, fmt.Errorf("left operand must an identifier or a scalar value")
	}

	right := expr.Right()
	switch right := right.(type) {
	case *expressions.Identifier:
		if ident != nil {
			return true, fmt.Errorf("left operand is an identifier, right operand must be a scalar value")
		}
		ident = right
	case *expressions.Scalar:
		if scalar != nil {
			return true, fmt.Errorf("left operand is an scalar value, right operand must be an identifier")
		}
		scalar = right
	default:
		return true, fmt.Errorf("right operand must an scalar value or an identifier")
	}
	switch ident.ID() {
	{{- range $name, $data := .FieldByFieldConstants}}
	case {{$name}}:
		typedVal, ok := any(scalar.Value()).({{ .Def.GoType }})
		if !ok {
			return true, fmt.Errorf(
				"invalid type on field {{$.CRUD.CamelCaseName}}.{{camelIdentifier .Def.GetName }}: expected type {{ .Def.GoType }}: got value %#v",
				scalar.Value(),
			)
		}
		return {{$.CRUD.CamelCaseName}}.{{camelIdentifier .Def.GetName }} == typedVal , nil
	{{- end}}
	}
	return false, nil{{else}}return true, fmt.Errorf("not implemented"){{end}}
}
{{end}}

{{if or .CRUD.Create .CRUD.Update .CRUD.Delete}}
{{range $name, $data := .InMemory.UIDDataByUIDNames .CRUD}}
func {{$data.BuildMapFnName}}({{$.CRUD.CamelCaseName}}s []*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}}) (
	map[{{.KeyType}}]uint,
	map[uint]{{.KeyType}},
	map[{{.KeyType}}]*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}},
	error,
) {
	{{$data.IndexByKeyMapName}} := make(map[{{.KeyType}}]uint)
	{{$data.KeyByIndexMapName}} := make(map[uint]{{.KeyType}})
	{{$name}} := make(map[{{.KeyType}}]*{{$.CRUD.GoType $.CRUD.File.GoPkg.Path}})
	{{if $data.KeyIsComposite}}
	{{- template "repository-in-memory-build-map-for-complex-keys" (inMemoryBuildMapWrap $.CRUD $name $data) -}}
	{{else}}
	{{- template "repository-in-memory-build-map-for-simple-keys" (inMemoryBuildMapWrap $.CRUD $name $data) -}}
	{{end}}
	return {{$data.IndexByKeyMapName}}, {{$data.KeyByIndexMapName}}, {{$name}}, nil
}
{{end}}
{{end}}
`))

	_ = template.Must(repositoryTemplate.New("repository-in-memory-build-map-for-simple-keys").Funcs(funcMap).Parse(`
	for i, def := range {{.CRUD.CamelCaseName}}s {
		{{.Data.IndexByKeyMapName}}[{{.Data.SimpleKeyGenerationCode "def"}}] = uint(i)
		{{.Data.KeyByIndexMapName}}[uint(i)] = {{.Data.SimpleKeyGenerationCode "def"}}
		{{.Name}}[{{.Data.SimpleKeyGenerationCode "def"}}] = def
	}
`))
	_ = template.Must(repositoryTemplate.New("repository-in-memory-build-map-for-complex-keys").Funcs(funcMap).Parse(`
	var err error
	h := sha256.New()
	for i, def := range {{$.CRUD.CamelCaseName}}s {
		{{.Data.CompositeKeyGenerationCode "def" "h"}}
		key := string(h.Sum(nil))
		{{.Data.IndexByKeyMapName}}[key] = uint(i)
		{{.Data.KeyByIndexMapName}}[uint(i)] = key
		{{.Name}}[key] = def
		h.Reset()
	}
`))
)
