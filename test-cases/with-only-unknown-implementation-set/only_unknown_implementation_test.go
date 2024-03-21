package with_only_unknown_implementation_set_test

import (
	"context"
	"strings"
	"testing"

	tested "github.com/samlitowitz/protoc-gen-crud/test-cases/with-only-unknown-implementation-set"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/samlitowitz/protoc-gen-crud/expressions"
)

func TestOnlyUnknownImplementationRepository_HasCompleteInterface(t *testing.T) {
	opts := cmp.Options{
		cmpopts.IgnoreUnexported(tested.OnlyUnknownImplementation{}),
		cmpopts.SortSlices(func(x, y *tested.OnlyUnknownImplementation) bool {
			switch strings.Compare(x.GetId(), y.GetId()) {
			case -1:
				return true
			case 0:
				return true
			case 1:
				return false
			}
			panic("this should never happen")
		}),
	}

	expected := []*tested.OnlyUnknownImplementation{}
	var iface tested.OnlyUnknownImplementationRepository = &testOnlyUnknownImplementationRepository{}

	resp, err := iface.Create(context.Background(), expected)
	if err != nil {
		t.Errorf("Create(): error %s", err)
	}
	if diff := cmp.Diff(expected, resp, opts); diff != "" {
		t.Errorf("Create() mistmatch (-want +got):\n%s", diff)
	}

	resp, err = iface.Read(context.Background(), nil)
	if err != nil {
		t.Errorf("Read(): error %s", err)
	}
	if diff := cmp.Diff(expected, resp, opts); diff != "" {
		t.Errorf("Read() mistmatch (-want +got):\n%s", diff)
	}

	resp, err = iface.Update(context.Background(), expected)
	if err != nil {
		t.Errorf("Update(): error %s", err)
	}
	if diff := cmp.Diff(expected, resp, opts); diff != "" {
		t.Errorf("Update() mistmatch (-want +got):\n%s", diff)
	}

	err = iface.Delete(context.Background(), nil)
	if err != nil {
		t.Errorf("Delete(): error %s", err)
	}
}

type testOnlyUnknownImplementationRepository struct{}

func (r *testOnlyUnknownImplementationRepository) Create(ctx context.Context, implementations []*tested.OnlyUnknownImplementation) ([]*tested.OnlyUnknownImplementation, error) {
	return []*tested.OnlyUnknownImplementation{}, nil
}

func (r *testOnlyUnknownImplementationRepository) Read(ctx context.Context, expression expressions.Expression) ([]*tested.OnlyUnknownImplementation, error) {
	return []*tested.OnlyUnknownImplementation{}, nil
}

func (r *testOnlyUnknownImplementationRepository) Update(ctx context.Context, implementations []*tested.OnlyUnknownImplementation) ([]*tested.OnlyUnknownImplementation, error) {
	return []*tested.OnlyUnknownImplementation{}, nil
}

func (r *testOnlyUnknownImplementationRepository) Delete(ctx context.Context, expression expressions.Expression) error {
	return nil
}
