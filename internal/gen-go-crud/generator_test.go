// REFURL: https://github.com/grpc-ecosystem/grpc-gateway/blob/v2.15.2/protoc-gen-grpc-gateway/internal/gengateway/generator_test.go
package gen_go_crud

import (
	"testing"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func newExampleFileDescriptorWithGoPkg(gp *descriptor.GoPackage, filenamePrefix string) *descriptor.File {
	msgdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("ExampleMessage"),
	}
	msg := &descriptor.Message{
		DescriptorProto: msgdesc,
	}
	meth := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("Example"),
		InputType:  proto.String("ExampleMessage"),
		OutputType: proto.String("ExampleMessage"),
	}
	meth1 := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("ExampleWithoutBindings"),
		InputType:  proto.String("empty.Empty"),
		OutputType: proto.String("empty.Empty"),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth, meth1},
	}
	return &descriptor.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			Name:        proto.String("example.proto"),
			Package:     proto.String("example"),
			Dependency:  []string{"a.example/b/c.proto", "a.example/d/e.proto"},
			MessageType: []*descriptorpb.DescriptorProto{msgdesc},
			Service:     []*descriptorpb.ServiceDescriptorProto{svc},
		},
		GoPkg:                   *gp,
		GeneratedFilenamePrefix: filenamePrefix,
		Messages:                []*descriptor.Message{msg},
	}
}

func TestGenerator_Generate(t *testing.T) {
	g := new(generator)
	g.reg = descriptor.NewRegistry()
	result, err := g.Generate([]*descriptor.File{
		crossLinkFixture(newExampleFileDescriptorWithGoPkg(&descriptor.GoPackage{
			Path: "example.com/path/to/example",
			Name: "example_pb",
		}, "path/to/example")),
	})
	if err != nil {
		t.Fatalf("failed to generate stubs: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected to generate one file, got: %d", len(result))
	}
	expectedName := "path/to/example.pb.gw.go"
	gotName := result[0].GetName()
	if gotName != expectedName {
		t.Fatalf("invalid name %q, expected %q", gotName, expectedName)
	}
}
