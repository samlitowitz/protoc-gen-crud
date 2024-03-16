package crud

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

func crossLinkFixture(f *descriptor.File) *descriptor.File {
	for _, m := range f.Messages {
		m.File = f
	}
	return f
}

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
	allOperations := []options.Operation{options.Operation_READ}
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

//func TestApplyTemplate_RepositorySQLiteRelationships(t *testing.T) {
//	//allOperations := []options.Operation{options.Operation_CREATE, options.Operation_READ, options.Operation_UPDATE, options.Operation_DELETE}
//	allOperations := []options.Operation{options.Operation_READ}
//	operationCombinations := allOperationCombinations(allOperations)
//
//	supportedScalarTypes := []string{
//		"int32",
//		"int64",
//		"uint32",
//		"uint64",
//		"sint32",
//		"sint64",
//		"fixed32",
//		"fixed64",
//		"sfixed32",
//		"sfixed64",
//		"string",
//		"bytes",
//	}
//	uidCombinations := allStringCombinations(supportedScalarTypes)
//
//	implementation := options.Implementation_SQLITE
//	for _, operations := range operationCombinations {
//		typNameString := "TYPE_STRING"
//		fieldLabelOptional := new(descriptorpb.FieldDescriptorProto_Label)
//		*fieldLabelOptional = descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
//		fieldTypString := new(descriptorpb.FieldDescriptorProto_Type)
//		*fieldTypString = descriptorpb.FieldDescriptorProto_Type(descriptorpb.FieldDescriptorProto_Type_value[typNameString])
//
//		// create child type
//		fieldDescs := []*descriptorpb.FieldDescriptorProto{
//			{
//				Name:     proto.String("id"),
//				Number:   proto.Int32(1),
//				Label:    fieldLabelOptional,
//				Type:     fieldTypString,
//				TypeName: proto.String(typNameString),
//			},
//		}
//
//		// create parent type
//		for _, uidCombination := range uidCombinations {
//			fieldDescs := make([]*descriptorpb.FieldDescriptorProto, 0, len(uidCombination))
//			fields := make([]*descriptor.Field, 0, len(uidCombination))
//			for i, typ := range uidCombination {
//				typName := "TYPE_" + strings.ToUpper(typ)
//				fieldLabel := new(descriptorpb.FieldDescriptorProto_Label)
//				*fieldLabel = descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
//				fieldTyp := new(descriptorpb.FieldDescriptorProto_Type)
//				*fieldTyp = descriptorpb.FieldDescriptorProto_Type(descriptorpb.FieldDescriptorProto_Type_value[typName])
//
//				fieldDesc := &descriptorpb.FieldDescriptorProto{
//					Name:     proto.String("Field_" + typ),
//					Number:   proto.Int32(int32(i)),
//					Label:    fieldLabel,
//					Type:     fieldTyp,
//					TypeName: proto.String(typName),
//				}
//				field := &descriptor.Field{
//					FieldDescriptorProto: fieldDesc,
//				}
//
//				fieldDescs = append(fieldDescs, fieldDesc)
//				fields = append(fields, field)
//			}
//
//			enumTypName := proto.String("EnumTyp")
//			enumDesc := &descriptorpb.EnumDescriptorProto{
//				Name: enumTypName,
//				Value: []*descriptorpb.EnumValueDescriptorProto{
//					{Name: proto.String("FALSE"), Number: proto.Int32(0)},
//					{Name: proto.String("TRUE"), Number: proto.Int32(1)},
//				},
//			}
//			enum := &descriptor.Enum{
//				EnumDescriptorProto: enumDesc,
//			}
//
//			typName := "TYPE_ENUM"
//			fieldLabel := new(descriptorpb.FieldDescriptorProto_Label)
//			*fieldLabel = descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
//			fieldTyp := new(descriptorpb.FieldDescriptorProto_Type)
//			*fieldTyp = descriptorpb.FieldDescriptorProto_Type(descriptorpb.FieldDescriptorProto_Type_value[typName])
//			fieldDescs = append(
//				fieldDescs,
//				&descriptorpb.FieldDescriptorProto{
//					Name:     proto.String("Field_Enum"),
//					Number:   proto.Int32(int32(len(uidCombination) + 1)),
//					Label:    fieldLabel,
//					Type:     fieldTyp,
//					TypeName: enumTypName,
//				},
//			)
//			uidName := "uidName"
//			uidMap := map[string][]*descriptor.Field{
//				uidName: fields,
//			}
//
//			msgDesc := &descriptorpb.DescriptorProto{
//				Name:  proto.String("ExampleMessageOne"),
//				Field: fieldDescs,
//			}
//			msg := &descriptor.Message{
//				DescriptorProto: msgDesc,
//				Fields:          fields,
//			}
//			for _, def := range fields {
//				def.Message = msg
//			}
//			crud := &descriptor.CRUD{
//				Message:            msg,
//				Operations:         make(map[options.Operation]struct{}),
//				Implementations:    map[options.Implementation]struct{}{implementation: {}},
//				UniqueIdentifiers:  uidMap,
//				FieldMaskFieldName: "test",
//			}
//			for _, operation := range operations {
//				crud.Operations[operation] = struct{}{}
//			}
//
//			file := descriptor.File{
//				FileDescriptorProto: &descriptorpb.FileDescriptorProto{
//					Name:        proto.String("example.proto"),
//					Package:     proto.String("example"),
//					EnumType:    []*descriptorpb.EnumDescriptorProto{enumDesc},
//					MessageType: []*descriptorpb.DescriptorProto{msgDesc},
//					Service:     []*descriptorpb.ServiceDescriptorProto{},
//				},
//				GoPkg: descriptor.GoPackage{
//					Path: "example.com/path/to/example/example.pb",
//					Name: "example_pb",
//				},
//				Messages: []*descriptor.Message{msg},
//				Enums:    []*descriptor.Enum{enum},
//				CRUDs:    []*descriptor.CRUD{crud},
//			}
//			got, err := applyTemplate(param{File: crossLinkFixture(&file)}, descriptor.NewRegistry())
//			if err != nil {
//				t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
//				return
//			}
//
//			uidCombinationStr := strings.Join(uidCombination, ", ")
//
//			// Assert struct definition
//			if want := "type SQLite" + *msgDesc.Name + "Repository struct {"; !strings.Contains(got, want) {
//				t.Errorf("%s: applyTemplate(%#v) = %s; want to contain %s", uidCombinationStr, file, got, want)
//			}
//
//			// Assert UID maps
//			goTyp, err := goTypeByFieldDescType(fields)
//			if err != nil {
//				t.Errorf("failed to generate go type for UID map: %s", err)
//			}
//			uidMapDef := fmt.Sprintf(
//				"%sBy%s map[%s]*%s",
//				strcase.ToLowerCamel(crud.GetName()),
//				strcase.ToCamel(uidName),
//				goTyp,
//				crud.GoType(crud.File.GoPkg.Path),
//			)
//			if want := uidMapDef; !strings.Contains(got, want) {
//				t.Errorf("%s: applyTemplate(%#v) = %s; want to contain %s", uidCombinationStr, file, got, want)
//			}
//
//			// Assert "constructor" function
//			if want := "func NewSQLite" + *msgDesc.Name + "Repository"; !strings.Contains(got, want) {
//				t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
//			}
//
//			for _, operation := range operations {
//				switch operation {
//				case options.Operation_CREATE:
//					want := fmt.Sprintf(
//						"func (repo *SQLite%sRepository) Create([]*%s) ([]*%s, error)",
//						*msgDesc.Name,
//						*msgDesc.Name,
//						*msgDesc.Name,
//					)
//					if !strings.Contains(got, want) {
//						t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
//					}
//				case options.Operation_READ:
//					want := fmt.Sprintf(
//						"func (repo *SQLite%sRepository) Read() ([]*%s, error)",
//						*msgDesc.Name,
//						*msgDesc.Name,
//					)
//					if !strings.Contains(got, want) {
//						t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
//					}
//				case options.Operation_UPDATE:
//					want := fmt.Sprintf(
//						"func (repo *SQLite%sRepository) Update([]*%s) ([]*%s, error)",
//						*msgDesc.Name,
//						*msgDesc.Name,
//						*msgDesc.Name,
//					)
//					if !strings.Contains(got, want) {
//						t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
//					}
//				case options.Operation_DELETE:
//					want := fmt.Sprintf(
//						"func (repo *SQLite%sRepository) Delete([]*%s) error",
//						*msgDesc.Name,
//						*msgDesc.Name,
//					)
//					if !strings.Contains(got, want) {
//						t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
//					}
//				}
//			}
//		}
//	}
//}

// REFURL: https://github.com/mxschmitt/golang-combinations/blob/main/combinations.go#L8
func allOperationCombinations(set []options.Operation) (subsets [][]options.Operation) {
	length := uint(len(set))

	// Go through all possible combinations of objects
	// from 1 (only first object in subset) to 2^length (all objects in subset)
	for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
		var subset []options.Operation

		for object := uint(0); object < length; object++ {
			// checks if object is contained in subset
			// by checking if bit 'object' is set in subsetBits
			if (subsetBits>>object)&1 == 1 {
				// add object to subset
				subset = append(subset, set[object])
			}
		}
		// add subset to subsets
		subsets = append(subsets, subset)
	}
	return subsets
}

func allStringCombinations(set []string) (subsets [][]string) {
	length := uint(len(set))

	// Go through all possible combinations of objects
	// from 1 (only first object in subset) to 2^length (all objects in subset)
	for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
		var subset []string

		for object := uint(0); object < length; object++ {
			// checks if object is contained in subset
			// by checking if bit 'object' is set in subsetBits
			if (subsetBits>>object)&1 == 1 {
				// add object to subset
				subset = append(subset, set[object])
			}
		}
		// add subset to subsets
		subsets = append(subsets, subset)
	}
	return subsets
}

func goTypeByFieldDescType(defs []*descriptor.Field) (string, error) {
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
		case descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_SINT32:
			return "int32", nil

		case descriptorpb.FieldDescriptorProto_TYPE_INT64:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
			fallthrough
		case descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
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
