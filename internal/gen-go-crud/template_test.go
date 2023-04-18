// REFURL: https://github.com/grpc-ecosystem/grpc-gateway/blob/v2.15.2/protoc-gen-grpc-gateway/internal/gengateway/template_test.go
package gen_go_crud

import (
	"strings"
	"testing"

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
	meth := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("Example"),
		InputType:  proto.String("ExampleMessage"),
		OutputType: proto.String("ExampleMessage"),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth},
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
			Service:     []*descriptorpb.ServiceDescriptorProto{svc},
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

func TestApplyTemplateInProcess(t *testing.T) {
	msgdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("ExampleMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     proto.String("nested"),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				TypeName: proto.String("NestedMessage"),
				Number:   proto.Int32(1),
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
	meth := &descriptorpb.MethodDescriptorProto{
		Name:            proto.String("Echo"),
		InputType:       proto.String("ExampleMessage"),
		OutputType:      proto.String("ExampleMessage"),
		ClientStreaming: proto.Bool(true),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth},
	}
	for _, spec := range []struct {
		clientStreaming bool
		serverStreaming bool
		sigWant         []string
	}{
		{
			clientStreaming: false,
			serverStreaming: false,
			sigWant: []string{
				`func local_request_ExampleService_Echo_0(ctx context.Context, marshaler runtime.Marshaler, server ExampleServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {`,
				`resp, md, err := local_request_ExampleService_Echo_0(annotatedContext, inboundMarshaler, server, req, pathParams)`,
			},
		},
		{
			clientStreaming: true,
			serverStreaming: true,
			sigWant: []string{
				`err := status.Error(codes.Unimplemented, "streaming calls are not yet supported in the in-process transport")`,
			},
		},
		{
			clientStreaming: true,
			serverStreaming: false,
			sigWant: []string{
				`err := status.Error(codes.Unimplemented, "streaming calls are not yet supported in the in-process transport")`,
			},
		},
		{
			clientStreaming: false,
			serverStreaming: true,
			sigWant: []string{
				`err := status.Error(codes.Unimplemented, "streaming calls are not yet supported in the in-process transport")`,
			},
		},
	} {
		meth.ClientStreaming = proto.Bool(spec.clientStreaming)
		meth.ServerStreaming = proto.Bool(spec.serverStreaming)

		msg := &descriptor.Message{
			DescriptorProto: msgdesc,
		}
		nested := &descriptor.Message{
			DescriptorProto: nesteddesc,
		}

		file := descriptor.File{
			FileDescriptorProto: &descriptorpb.FileDescriptorProto{
				Name:        proto.String("example.proto"),
				Package:     proto.String("example"),
				MessageType: []*descriptorpb.DescriptorProto{msgdesc, nesteddesc},
				Service:     []*descriptorpb.ServiceDescriptorProto{svc},
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

		for _, want := range spec.sigWant {
			if !strings.Contains(got, want) {
				t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
			}
		}

		if want := `func RegisterExampleServiceHandlerServer(ctx context.Context, mux *runtime.ServeMux, server ExampleServiceServer) error {`; !strings.Contains(got, want) {
			t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
		}
	}
}

func TestIdentifierCapitalization(t *testing.T) {
	msgdesc1 := &descriptorpb.DescriptorProto{
		Name: proto.String("Exam_pleRequest"),
	}
	msgdesc2 := &descriptorpb.DescriptorProto{
		Name: proto.String("example_response"),
	}
	meth1 := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("ExampleGe2t"),
		InputType:  proto.String("Exam_pleRequest"),
		OutputType: proto.String("example_response"),
	}
	meth2 := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("Exampl_ePost"),
		InputType:  proto.String("Exam_pleRequest"),
		OutputType: proto.String("example_response"),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("Example"),
		Method: []*descriptorpb.MethodDescriptorProto{meth1, meth2},
	}
	msg1 := &descriptor.Message{
		DescriptorProto: msgdesc1,
	}
	msg2 := &descriptor.Message{
		DescriptorProto: msgdesc2,
	}
	file := descriptor.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			Name:        proto.String("example.proto"),
			Package:     proto.String("example"),
			Dependency:  []string{"a.example/b/c.proto", "a.example/d/e.proto"},
			MessageType: []*descriptorpb.DescriptorProto{msgdesc1, msgdesc2},
			Service:     []*descriptorpb.ServiceDescriptorProto{svc},
		},
		GoPkg: descriptor.GoPackage{
			Path: "example.com/path/to/example/example.pb",
			Name: "example_pb",
		},
		Messages: []*descriptor.Message{msg1, msg2},
	}

	got, err := applyTemplate(param{File: crossLinkFixture(&file)}, descriptor.NewRegistry())
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}
	if want := `msg, err := client.ExampleGe2T(ctx, &protoReq, grpc.Header(&metadata.HeaderMD)`; !strings.Contains(got, want) {
		t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
	}
	if want := `msg, err := client.ExamplEPost(ctx, &protoReq, grpc.Header(&metadata.HeaderMD)`; !strings.Contains(got, want) {
		t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
	}
	if want := `var protoReq ExamPleRequest`; !strings.Contains(got, want) {
		t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
	}
	if want := `var protoReq ExampleResponse`; !strings.Contains(got, want) {
		t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
	}
}
