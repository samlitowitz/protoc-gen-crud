// REFURL: https://github.com/grpc-ecosystem/grpc-gateway/blob/v2.15.2/protoc-gen-grpc-gateway/internal/gengateway/template_test.go
package crud

import (
	"fmt"
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
	fileName := "example.proto"
	goPkgName := "example_pb"

	msgdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("ExampleMessage"),
	}
	msg := &descriptor.Message{
		DescriptorProto: msgdesc,
	}
	file := descriptor.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			Name:        proto.String(fileName),
			Package:     proto.String("example"),
			Dependency:  []string{"a.example/b/c.proto", "a.example/d/e.proto"},
			MessageType: []*descriptorpb.DescriptorProto{msgdesc},
			Service:     []*descriptorpb.ServiceDescriptorProto{},
		},
		GoPkg: descriptor.GoPackage{
			Path: "example.com/path/to/example/example.pb",
			Name: goPkgName,
		},
		Messages: []*descriptor.Message{msg},
	}
	got, err := applyTemplate(param{File: crossLinkFixture(&file)}, descriptor.NewRegistry())
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}
	if want := "// source: " + fileName + "\n"; !strings.Contains(got, want) {
		t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
	}
	if want := "package " + goPkgName + "\n"; !strings.Contains(got, want) {
		t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
	}
}

func TestApplyTemplate_RepositoryInterface(t *testing.T) {
	allOperations := []options.Operation{options.Operation_CREATE, options.Operation_READ, options.Operation_UPDATE, options.Operation_DELETE}
	combinations := allOperationCombinations(allOperations)
	for _, operations := range combinations {
		msgdesc := &descriptorpb.DescriptorProto{
			Name: proto.String("ExampleMessageOne"),
		}
		msg := &descriptor.Message{
			DescriptorProto: msgdesc,
		}
		crud := &descriptor.CRUD{
			Message:    msg,
			Operations: make(map[options.Operation]struct{}),
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

		if want := "type " + *msgdesc.Name + "Repository interface {"; !strings.Contains(got, want) {
			t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
		}

		for _, operation := range operations {
			switch operation {
			case options.Operation_CREATE:
				want := fmt.Sprintf(
					"Create([]*%s) ([]*%s, error)",
					*msgdesc.Name,
					*msgdesc.Name,
				)
				if !strings.Contains(got, want) {
					t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
				}
			case options.Operation_READ:
				want := fmt.Sprintf(
					"Read() ([]*%s, error)",
					*msgdesc.Name,
				)
				if !strings.Contains(got, want) {
					t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
				}
			case options.Operation_UPDATE:
				want := fmt.Sprintf(
					"Update([]*%s) ([]*%s, error)",
					*msgdesc.Name,
					*msgdesc.Name,
				)
				if !strings.Contains(got, want) {
					t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
				}
			case options.Operation_DELETE:
				want := fmt.Sprintf(
					"Delete([]*%s) error",
					*msgdesc.Name,
				)
				if !strings.Contains(got, want) {
					t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
				}
			}
		}
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
