package descriptor

import (
	"fmt"
	"strings"
	"testing"

	"github.com/samlitowitz/protoc-gen-crud/options"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func testExtractCRUDs(t *testing.T, input []*descriptorpb.FileDescriptorProto, target string, wantCRUDs []*CRUD) {
	testExtractCRUDsWithRegistry(t, NewRegistry(), input, target, wantCRUDs)
}

func testExtractCRUDsWithRegistry(t *testing.T, reg *Registry, input []*descriptorpb.FileDescriptorProto, target string, wantCRUDs []*CRUD) {
	for _, file := range input {
		reg.loadFile(file.GetName(), &protogen.File{
			Proto: file,
		})
	}
	err := reg.loadCRUDs(reg.files[target])
	if err != nil {
		t.Errorf("loadCRUDs(%q) failed with %v; want success; files=%v", target, err, input)
	}

	file := reg.files[target]
	cruds := file.CRUDs
	var i int
	for i = 0; i < len(cruds) && i < len(wantCRUDs); i++ {
		crud, wantCRUD := cruds[i], wantCRUDs[i]
		// Ensure same message proto
		if got, want := crud.Message.DescriptorProto, wantCRUD.Message.DescriptorProto; !proto.Equal(got, want) {
			t.Errorf("cruds[%d].DescriptorProto = %v; want %v; input = %v", i, got, want, input)
			continue
		}

		// Ensure same operations
		if len(crud.Operations) != len(wantCRUD.Operations) {
			t.Errorf("len(cruds[%d].Operations) = %d; want %d; input %v", i, len(crud.Operations), len(wantCRUD.Operations), input)
		}
		for wantOperation := range wantCRUD.Operations {
			if _, ok := crud.Operations[wantOperation]; !ok {
				t.Errorf("cruds[%d].Operations = nil; want %s; input %v", i, wantOperation, input)
			}
		}

		// Ensure same implementations
		if len(crud.Implementations) != len(wantCRUD.Implementations) {
			t.Errorf("len(cruds[%d].Implementations) = %d; want %d; input %v", i, len(crud.Implementations), len(wantCRUD.Implementations), input)
		}
		for wantImplementation := range wantCRUD.Implementations {
			if _, ok := crud.Implementations[wantImplementation]; !ok {
				t.Errorf("cruds[%d].Implementations = nil; want %s; input %v", i, wantImplementation, input)
			}
		}

		// Ensure same field options
	}

	for ; i < len(cruds); i++ {
		got := cruds[i].Message.DescriptorProto
		t.Errorf("cruds[%d] = %v; want it to be missing; input = %v", i, got, input)
	}
	for ; i < len(wantCRUDs); i++ {
		want := wantCRUDs[i].Message.DescriptorProto
		t.Errorf("cruds[%d] missing; want %v; input = %v", i, want, input)
	}
}

func crossLinkFixture(f *File) *File {
	for _, m := range f.Messages {
		m.File = f
		for _, f := range m.Fields {
			f.Message = m
		}
	}
	for _, crud := range f.CRUDs {
		crud.Message.File = f
	}
	for _, e := range f.Enums {
		e.File = f
	}
	return f
}

func TestExtractCRUDsWithoutAnnotation(t *testing.T) {
	src := `
		name: "path/to/example.proto",
		package: "example"
		message_type <
			name: "StringMessage"
			field <
				name: "string"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
			>
		>
	`
	var fd descriptorpb.FileDescriptorProto
	if err := prototext.Unmarshal([]byte(src), &fd); err != nil {
		t.Fatalf("proto.UnmarshalText(%s, &fd) failed with %v; want success", src, err)
	}
	msg := &Message{
		DescriptorProto: fd.MessageType[0],
		Fields: []*Field{
			{
				FieldDescriptorProto: fd.MessageType[0].Field[0],
			},
		},
	}
	file := &File{
		FileDescriptorProto: &fd,
		GoPkg: GoPackage{
			Path: "path/to/example.pb",
			Name: "example_pb",
		},
		Messages: []*Message{msg},
		CRUDs:    []*CRUD{},
	}

	crossLinkFixture(file)
	testExtractCRUDs(t, []*descriptorpb.FileDescriptorProto{&fd}, "path/to/example.proto", file.CRUDs)
}

func TestExtractCRUDOperations(t *testing.T) {
	allOperations := []options.Operation{options.Operation_CREATE, options.Operation_READ, options.Operation_UPDATE, options.Operation_DELETE}
	combinations := allOperationCombinations(allOperations)
	for _, operations := range combinations {
		src := `
		name: "path/to/example.proto",
		package: "example"
		message_type <
			name: "StringMessage"
			field <
				name: "string"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
			>
			options <
				[protoc_gen_crud.options.crud_message_options] <
					operations: [%s]
				>
			>
		>
	`
		operationsLiteral := make([]string, 0, len(operations))
		operationsMap := make(map[options.Operation]struct{})
		for _, operation := range operations {
			operationsLiteral = append(operationsLiteral, operation.String())
			operationsMap[operation] = struct{}{}
		}

		src = fmt.Sprintf(src, strings.Join(operationsLiteral, ", "))

		var fd descriptorpb.FileDescriptorProto
		if err := prototext.Unmarshal([]byte(src), &fd); err != nil {
			t.Fatalf("proto.UnmarshalText(%s, &fd) failed with %v; want success", src, err)
		}
		msg := &Message{
			DescriptorProto: fd.MessageType[0],
			Fields: []*Field{
				{
					FieldDescriptorProto: fd.MessageType[0].Field[0],
				},
			},
		}
		file := &File{
			FileDescriptorProto: &fd,
			GoPkg: GoPackage{
				Path: "path/to/example.pb",
				Name: "example_pb",
			},
			Messages: []*Message{msg},
			CRUDs: []*CRUD{
				{
					Message:    msg,
					Operations: operationsMap,
				},
			},
		}

		crossLinkFixture(file)
		testExtractCRUDs(t, []*descriptorpb.FileDescriptorProto{&fd}, "path/to/example.proto", file.CRUDs)
	}
}

func TestExtractCRUDImplementations(t *testing.T) {
	allImplementations := []options.Implementation{options.Implementation_IN_MEMORY}
	combinations := allImplementationCombinations(allImplementations)
	for _, implementations := range combinations {
		src := `
		name: "path/to/example.proto",
		package: "example"
		message_type <
			name: "StringMessage"
			field <
				name: "string"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
			>
			options <
				[protoc_gen_crud.options.crud_message_options] <
					implementations: [%s]
				>
			>
		>
	`
		implementationsLiteral := make([]string, 0, len(implementations))
		implementationsMap := make(map[options.Implementation]struct{})
		for _, implementation := range implementations {
			implementationsLiteral = append(implementationsLiteral, implementation.String())
			implementationsMap[implementation] = struct{}{}
		}

		src = fmt.Sprintf(src, strings.Join(implementationsLiteral, ", "))

		var fd descriptorpb.FileDescriptorProto
		if err := prototext.Unmarshal([]byte(src), &fd); err != nil {
			t.Fatalf("proto.UnmarshalText(%s, &fd) failed with %v; want success", src, err)
		}
		msg := &Message{
			DescriptorProto: fd.MessageType[0],
			Fields: []*Field{
				{
					FieldDescriptorProto: fd.MessageType[0].Field[0],
				},
			},
		}
		file := &File{
			FileDescriptorProto: &fd,
			GoPkg: GoPackage{
				Path: "path/to/example.pb",
				Name: "example_pb",
			},
			Messages: []*Message{msg},
			CRUDs: []*CRUD{
				{
					Message:         msg,
					Implementations: implementationsMap,
				},
			},
		}

		crossLinkFixture(file)
		testExtractCRUDs(t, []*descriptorpb.FileDescriptorProto{&fd}, "path/to/example.proto", file.CRUDs)
	}
}

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

func allImplementationCombinations(set []options.Implementation) (subsets [][]options.Implementation) {
	length := uint(len(set))

	// Go through all possible combinations of objects
	// from 1 (only first object in subset) to 2^length (all objects in subset)
	for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
		var subset []options.Implementation

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
