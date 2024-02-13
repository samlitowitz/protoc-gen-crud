package gen_go_crud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/iancoleman/strcase"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
	"github.com/samlitowitz/protoc-gen-crud/options"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestApplyTemplate_RepositorySQLite(t *testing.T) {
	allOperations := []options.Operation{options.Operation_CREATE, options.Operation_READ, options.Operation_UPDATE, options.Operation_DELETE}
	operationCombinations := allOperationCombinations(allOperations)

	implementation := options.Implementation_SQLITE
	for _, operations := range operationCombinations {
		msgdesc := &descriptorpb.DescriptorProto{
			Name: proto.String("ExampleMessageOne"),
		}
		msg := &descriptor.Message{
			DescriptorProto: msgdesc,
		}
		crud := &descriptor.CRUD{
			Message:         msg,
			Operations:      make(map[options.Operation]struct{}),
			Implementations: map[options.Implementation]struct{}{implementation: {}},
		}
		for _, operation := range operations {
			crud.Operations[operation] = struct{}{}
		}

		file := descriptor.File{
			FileDescriptorProto: &descriptorpb.FileDescriptorProto{
				Name:        proto.String("example.proto"),
				Package:     proto.String("example"),
				MessageType: []*descriptorpb.DescriptorProto{msgdesc},
				Service:     []*descriptorpb.ServiceDescriptorProto{},
			},
			GoPkg: descriptor.GoPackage{
				Path: "example.com/path/to/example/example.pb",
				Name: "example_pb",
			},
			Messages: []*descriptor.Message{msg},
			CRUDs:    []*descriptor.CRUD{crud},
		}
		got, err := applyTemplate(param{File: crossLinkFixture(&file)}, descriptor.NewRegistry())
		if err != nil {
			t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
			return
		}

		// Assert struct definition
		if want := "type SQLite" + *msgdesc.Name + "Repository struct {"; !strings.Contains(got, want) {
			t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
		}

		// Assert "constructor" function
		if want := "func NewSQLite" + *msgdesc.Name + "Repository"; !strings.Contains(got, want) {
			t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
		}

		for _, operation := range operations {
			switch operation {
			case options.Operation_CREATE:
				want := fmt.Sprintf(
					"func (repo *SQLite%sRepository) Create([]*%s) ([]*%s, error)",
					*msgdesc.Name,
					*msgdesc.Name,
					*msgdesc.Name,
				)
				if !strings.Contains(got, want) {
					t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
				}
			case options.Operation_READ:
				want := fmt.Sprintf(
					"func (repo *SQLite%sRepository) Read() ([]*%s, error)",
					*msgdesc.Name,
					*msgdesc.Name,
				)
				if !strings.Contains(got, want) {
					t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
				}
			case options.Operation_UPDATE:
				want := fmt.Sprintf(
					"func (repo *SQLite%sRepository) Update([]*%s) ([]*%s, error)",
					*msgdesc.Name,
					*msgdesc.Name,
					*msgdesc.Name,
				)
				if !strings.Contains(got, want) {
					t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
				}
			case options.Operation_DELETE:
				want := fmt.Sprintf(
					"func (repo *SQLite%sRepository) Delete([]*%s) error",
					*msgdesc.Name,
					*msgdesc.Name,
				)
				if !strings.Contains(got, want) {
					t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
				}
			}
		}
	}
}

func TestApplyTemplate_RepositorySQLiteUIDs(t *testing.T) {
	//allOperations := []options.Operation{options.Operation_CREATE, options.Operation_READ, options.Operation_UPDATE, options.Operation_DELETE}
	allOperations := []options.Operation{options.Operation_CREATE}
	operationCombinations := allOperationCombinations(allOperations)

	supportedScalarTypes := []string{
		"int32",
		"int64",
		"uint32",
		"uint64",
		"sint32",
		"sint64",
		"fixed32",
		"fixed64",
		"sfixed32",
		"sfixed64",
		"string",
		"bytes",
	}
	uidCombinations := allStringCombinations(supportedScalarTypes)

	implementation := options.Implementation_SQLITE
	for _, operations := range operationCombinations {
		for _, uidCombination := range uidCombinations {
			fieldDescs := make([]*descriptorpb.FieldDescriptorProto, 0, len(uidCombination))
			fields := make([]*descriptor.Field, 0, len(uidCombination))
			for i, typ := range uidCombination {
				typName := "TYPE_" + strings.ToUpper(typ)
				fieldLabel := new(descriptorpb.FieldDescriptorProto_Label)
				*fieldLabel = descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
				fieldTyp := new(descriptorpb.FieldDescriptorProto_Type)
				*fieldTyp = descriptorpb.FieldDescriptorProto_Type(descriptorpb.FieldDescriptorProto_Type_value[typName])

				fieldDesc := &descriptorpb.FieldDescriptorProto{
					Name:     proto.String("Field_" + typ),
					Number:   proto.Int32(int32(i)),
					Label:    fieldLabel,
					Type:     fieldTyp,
					TypeName: proto.String(typName),
				}
				field := &descriptor.Field{
					FieldDescriptorProto: fieldDesc,
				}

				fieldDescs = append(fieldDescs, fieldDesc)
				fields = append(fields, field)
			}

			enumTypName := proto.String("EnumTyp")
			enumDesc := &descriptorpb.EnumDescriptorProto{
				Name: enumTypName,
				Value: []*descriptorpb.EnumValueDescriptorProto{
					{Name: proto.String("FALSE"), Number: proto.Int32(0)},
					{Name: proto.String("TRUE"), Number: proto.Int32(1)},
				},
			}
			enum := &descriptor.Enum{
				EnumDescriptorProto: enumDesc,
			}

			typName := "TYPE_ENUM"
			fieldLabel := new(descriptorpb.FieldDescriptorProto_Label)
			*fieldLabel = descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
			fieldTyp := new(descriptorpb.FieldDescriptorProto_Type)
			*fieldTyp = descriptorpb.FieldDescriptorProto_Type(descriptorpb.FieldDescriptorProto_Type_value[typName])
			fieldDescs = append(
				fieldDescs,
				&descriptorpb.FieldDescriptorProto{
					Name:     proto.String("Field_Enum"),
					Number:   proto.Int32(int32(len(uidCombination) + 1)),
					Label:    fieldLabel,
					Type:     fieldTyp,
					TypeName: enumTypName,
				},
			)
			uidName := "uidName"
			uidMap := map[string][]*descriptor.Field{
				uidName: fields,
			}

			msgDesc := &descriptorpb.DescriptorProto{
				Name:  proto.String("ExampleMessageOne"),
				Field: fieldDescs,
			}
			msg := &descriptor.Message{
				DescriptorProto: msgDesc,
				Fields:          fields,
			}
			for _, def := range fields {
				def.Message = msg
			}
			crud := &descriptor.CRUD{
				Message:            msg,
				Operations:         make(map[options.Operation]struct{}),
				Implementations:    map[options.Implementation]struct{}{implementation: {}},
				UniqueIdentifiers:  uidMap,
				FieldMaskFieldName: "test",
			}
			for _, operation := range operations {
				crud.Operations[operation] = struct{}{}
			}

			file := descriptor.File{
				FileDescriptorProto: &descriptorpb.FileDescriptorProto{
					Name:        proto.String("example.proto"),
					Package:     proto.String("example"),
					EnumType:    []*descriptorpb.EnumDescriptorProto{enumDesc},
					MessageType: []*descriptorpb.DescriptorProto{msgDesc},
					Service:     []*descriptorpb.ServiceDescriptorProto{},
				},
				GoPkg: descriptor.GoPackage{
					Path: "example.com/path/to/example/example.pb",
					Name: "example_pb",
				},
				Messages: []*descriptor.Message{msg},
				Enums:    []*descriptor.Enum{enum},
				CRUDs:    []*descriptor.CRUD{crud},
			}
			got, err := applyTemplate(param{File: crossLinkFixture(&file)}, descriptor.NewRegistry())
			if err != nil {
				t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
				return
			}

			uidCombinationStr := strings.Join(uidCombination, ", ")

			// Assert struct definition
			if want := "type SQLite" + *msgDesc.Name + "Repository struct {"; !strings.Contains(got, want) {
				t.Errorf("%s: applyTemplate(%#v) = %s; want to contain %s", uidCombinationStr, file, got, want)
			}

			// Assert UID maps
			goTyp, err := goTypeByFieldDescType(fields)
			if err != nil {
				t.Errorf("failed to generate go type for UID map: %s", err)
			}
			uidMapDef := fmt.Sprintf(
				"%sBy%s map[%s]*%s",
				strcase.ToLowerCamel(crud.GetName()),
				strcase.ToCamel(uidName),
				goTyp,
				crud.GoType(crud.File.GoPkg.Path),
			)
			if want := uidMapDef; !strings.Contains(got, want) {
				t.Errorf("%s: applyTemplate(%#v) = %s; want to contain %s", uidCombinationStr, file, got, want)
			}

			// Assert "constructor" function
			if want := "func NewSQLite" + *msgDesc.Name + "Repository"; !strings.Contains(got, want) {
				t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
			}

			for _, operation := range operations {
				switch operation {
				case options.Operation_CREATE:
					want := fmt.Sprintf(
						"func (repo *SQLite%sRepository) Create([]*%s) ([]*%s, error)",
						*msgDesc.Name,
						*msgDesc.Name,
						*msgDesc.Name,
					)
					if !strings.Contains(got, want) {
						t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
					}
				case options.Operation_READ:
					want := fmt.Sprintf(
						"func (repo *SQLite%sRepository) Read() ([]*%s, error)",
						*msgDesc.Name,
						*msgDesc.Name,
					)
					if !strings.Contains(got, want) {
						t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
					}
				case options.Operation_UPDATE:
					want := fmt.Sprintf(
						"func (repo *SQLite%sRepository) Update([]*%s) ([]*%s, error)",
						*msgDesc.Name,
						*msgDesc.Name,
						*msgDesc.Name,
					)
					if !strings.Contains(got, want) {
						t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
					}
				case options.Operation_DELETE:
					want := fmt.Sprintf(
						"func (repo *SQLite%sRepository) Delete([]*%s) error",
						*msgDesc.Name,
						*msgDesc.Name,
					)
					if !strings.Contains(got, want) {
						t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
					}
				}
			}
		}
	}
}
