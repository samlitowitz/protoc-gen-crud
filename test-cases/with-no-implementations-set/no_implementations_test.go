package with_no_implementations_set_test

import (
	"context"
	"strings"
	"testing"

	tested "github.com/samlitowitz/protoc-gen-crud/test-cases/with-no-implementations-set"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/samlitowitz/expressions"
)

func TestNoImplementationsRepository_HasCompleteInterface(t *testing.T) {
	opts := cmp.Options{
		cmpopts.IgnoreUnexported(tested.NoImplementations{}),
		cmpopts.SortSlices(func(x, y *tested.NoImplementations) bool {
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

	expected := []*tested.NoImplementations{}
	var iface tested.NoImplementationsRepository = &testNoImplementationsRepository{}

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

type testNoImplementationsRepository struct{}

func (r *testNoImplementationsRepository) Create(ctx context.Context, implementations []*tested.NoImplementations) ([]*tested.NoImplementations, error) {
	return []*tested.NoImplementations{}, nil
}

func (r *testNoImplementationsRepository) Read(ctx context.Context, expression expressions.Expression) ([]*tested.NoImplementations, error) {
	return []*tested.NoImplementations{}, nil
}

func (r *testNoImplementationsRepository) Update(ctx context.Context, implementations []*tested.NoImplementations) ([]*tested.NoImplementations, error) {
	return []*tested.NoImplementations{}, nil
}

func (r *testNoImplementationsRepository) Delete(ctx context.Context, expression expressions.Expression) error {
	return nil
}
