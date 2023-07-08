package gen_go_crud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
	"github.com/samlitowitz/protoc-gen-crud/options"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestApplyTemplate_RepositoryInMemory(t *testing.T) {
	allOperations := []options.Operation{options.Operation_CREATE, options.Operation_READ, options.Operation_UPDATE, options.Operation_DELETE}
	operationCombinations := allOperationCombinations(allOperations)
	implementation := options.Implementation_IN_MEMORY
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
		if want := "type InMemory" + *msgdesc.Name + "Repository struct {"; !strings.Contains(got, want) {
			t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
		}

		// Assert "constructor" function
		if want := "func NewInMemory" + *msgdesc.Name + "Repository"; !strings.Contains(got, want) {
			t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
		}

		for _, operation := range operations {
			switch operation {
			case options.Operation_CREATE:
				want := fmt.Sprintf(
					"func (repo *InMemory%sRepository) Create([]*%s) ([]*%s, error)",
					*msgdesc.Name,
					*msgdesc.Name,
					*msgdesc.Name,
				)
				if !strings.Contains(got, want) {
					t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
				}
			case options.Operation_READ:
				want := fmt.Sprintf(
					"func (repo *InMemory%sRepository) Read() ([]*%s, error)",
					*msgdesc.Name,
					*msgdesc.Name,
				)
				if !strings.Contains(got, want) {
					t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
				}
			case options.Operation_UPDATE:
				want := fmt.Sprintf(
					"func (repo *InMemory%sRepository) Update([]*%s) ([]*%s, error)",
					*msgdesc.Name,
					*msgdesc.Name,
					*msgdesc.Name,
				)
				if !strings.Contains(got, want) {
					t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
				}
			case options.Operation_DELETE:
				want := fmt.Sprintf(
					"func (repo *InMemory%sRepository) Delete([]*%s) error",
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
