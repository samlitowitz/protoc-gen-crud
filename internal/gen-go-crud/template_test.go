// REFURL: https://github.com/grpc-ecosystem/grpc-gateway/blob/v2.15.2/protoc-gen-grpc-gateway/internal/gengateway/template_test.go
package gen_go_crud

import (
	"strings"
	"testing"

	"github.com/samlitowitz/protoc-gen-crud/options"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func crossLinkFixture(f *descriptor.File) *descriptor.File {
	for _, m := range f.Messages {
		m.File = f
	}
	return f
}

func TestApplyTemplateHeader(t *testing.T) {
	msgdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("ExampleMessage"),
	}
	msg := &descriptor.Message{
		DescriptorProto: msgdesc,
	}
	file := descriptor.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			Name:        proto.String("example.proto"),
			Package:     proto.String("example"),
			Dependency:  []string{"a.example/b/c.proto", "a.example/d/e.proto"},
			MessageType: []*descriptorpb.DescriptorProto{msgdesc},
			Service:     []*descriptorpb.ServiceDescriptorProto{},
		},
		GoPkg: descriptor.GoPackage{
			Path: "example.com/path/to/example/example.pb",
			Name: "example_pb",
		},
		Messages: []*descriptor.Message{msg},
	}
	got, err := applyTemplate(param{File: crossLinkFixture(&file)}, descriptor.NewRegistry())
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}
	if want := "package example_pb\n"; !strings.Contains(got, want) {
		t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
	}
}

func TestApplyRepositoryTemplate(t *testing.T) {
	msgdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("ExampleMessageOne"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     proto.String("id_one"),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				TypeName: proto.String(".google.protobuf.StringValue"),
				Number:   proto.Int32(1),
			},
			{
				Name:     proto.String("id_two"),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				TypeName: proto.String(".google.protobuf.StringValue"),
				Number:   proto.Int32(2),
			},
			{
				Name:     proto.String("id_three"),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				TypeName: proto.String(".google.protobuf.StringValue"),
				Number:   proto.Int32(3),
			},
			{
				Name:     proto.String("nested"),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				TypeName: proto.String("NestedMessage"),
				Number:   proto.Int32(4),
			},
		},
	}
	nesteddesc := &descriptorpb.DescriptorProto{
		Name: proto.String("NestedMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("int32"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
				Number: proto.Int32(1),
			},
			{
				Name:   proto.String("bool"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_BOOL.Enum(),
				Number: proto.Int32(2),
			},
		},
	}
	msg := &descriptor.Message{
		DescriptorProto: msgdesc,
	}
	crud := &descriptor.CRUD{
		Message: msg,
		Operations: map[options.Operation]struct{}{
			options.Operation_CREATE: {},
			options.Operation_READ:   {},
			options.Operation_UPDATE: {},
			options.Operation_DELETE: {},
		},
		UniqueIdentifiers: map[string][]*descriptor.Field{
			"one": {
				{
					FieldDescriptorProto: msgdesc.Field[0],
					Message:              msg,
					ForcePrefixedName:    false,
				},
			},
			"two": {
				{
					FieldDescriptorProto: msgdesc.Field[0],
					Message:              msg,
					ForcePrefixedName:    false,
				},
				{
					FieldDescriptorProto: msgdesc.Field[1],
					Message:              msg,
					ForcePrefixedName:    false,
				},
			},
			"three": {
				{
					FieldDescriptorProto: msgdesc.Field[0],
					Message:              msg,
					ForcePrefixedName:    false,
				},
				{
					FieldDescriptorProto: msgdesc.Field[1],
					Message:              msg,
					ForcePrefixedName:    false,
				},
				{
					FieldDescriptorProto: msgdesc.Field[2],
					Message:              msg,
					ForcePrefixedName:    false,
				},
			},
		},
	}
	msg.CRUD = crud
	nested := &descriptor.Message{
		DescriptorProto: nesteddesc,
	}

	file := descriptor.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			Name:        proto.String("example.proto"),
			Package:     proto.String("example"),
			MessageType: []*descriptorpb.DescriptorProto{msgdesc, nesteddesc},
			Service:     []*descriptorpb.ServiceDescriptorProto{},
		},
		GoPkg: descriptor.GoPackage{
			Path: "example.com/path/to/example/example.pb",
			Name: "example_pb",
		},
		Messages: []*descriptor.Message{msg, nested},
	}
	got, err := applyTemplate(param{File: crossLinkFixture(&file)}, descriptor.NewRegistry())
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}

	//for _, want := range spec.sigWant {
	//	if !strings.Contains(got, want) {
	//		t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
	//	}
	//}

	if want := `func RegisterExampleServiceHandlerServer(ctx context.Context, mux *runtime.ServeMux, server ExampleServiceServer) error {`; !strings.Contains(got, want) {
		t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
	}
}
