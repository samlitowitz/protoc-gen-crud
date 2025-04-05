package field_mask_test

import (
	"context"
	"testing"

	"google.golang.org/genproto/protobuf/field_mask"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	fieldMask "github.com/samlitowitz/protoc-gen-crud/test-cases/field-mask"
)

func TestSAInt32Repository_Create_WithAnyPrimeAttributeExcludedByFieldMaskFails(t *testing.T) {
	opts := saInt32DefaultCmpOpts()

	testCases := map[string][]*fieldMask.SAInt32_builder{
		"all fields excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{""},
				},
				Id:   0,
				Data: "SHOULD FAIL",
			},
		},
		"id field exlucded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"data"},
				},
				Id:   0,
				Data: "SHOULD FAIL",
			},
		},
	}

	for repoDesc, componentUnderTest := range saInt32ImplementationsToTest() {
		for testDesc, testCase := range testCases {
			// Call setup function, inject t *testing.T, and use t.Cleanup
			repoImpl := componentUnderTest(t)
			if repoImpl == nil {
				t.Fatalf(
					"%s: no implementation provided",
					repoDesc,
				)
			}

			err := sqliteSAInt32ReadCheck(opts, repoImpl, context.Background(), nil, nil)
			if err != nil {
				t.Fatalf("%s: %s", repoDesc, err)
			}
			toCreate := saInt32Build(testCase)
			_, err = repoImpl.Create(context.Background(), toCreate)
			if err == nil {
				t.Fatalf(
					"%s: %s: Create: expected field masking prime attribute error",
					repoDesc,
					testDesc,
				)
			}

			err = sqliteSAInt32ReadCheck(opts, repoImpl, context.Background(), nil, nil)
			if err != nil {
				t.Fatalf("%s: %s: %s", repoDesc, testDesc, err)
			}
		}

	}
}

func TestSAInt32Repository_Create_WithNonPrimeAttributesExcludedByFieldMaskUsesEmptyValues(t *testing.T) {
	opts := saInt32DefaultCmpOpts()

	for repoDesc, componentUnderTest := range saInt32ImplementationsToTest() {
		// Call setup function, inject t *testing.T, and use t.Cleanup
		repoImpl := componentUnderTest(t)
		if repoImpl == nil {
			t.Fatalf(
				"%s: no implementation provided",
				repoDesc,
			)
		}

		res, err := repoImpl.Read(context.Background(), nil)
		if err != nil {
			t.Fatalf(
				"%s: initial condition: empty repo: %s",
				repoDesc,
				err,
			)
		}
		if len(res) != 0 {
			t.Fatalf(
				"%s: initial condition: empty repo: got %d items",
				repoDesc,
				len(res),
			)
		}
		toCreateBuilder := []*fieldMask.SAInt32_builder{
			{
				Id:   0,
				Data: "SHOULD BE SET - 0",
			},
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id", "data"},
				},
				Id:   1,
				Data: "SHOULD SET - 1",
			},
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id", "data"},
				},
				Id:   2,
				Data: "SHOULD SET - 2",
			},
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id"},
				},
				Id:   3,
				Data: "SHOULD BE EMPTY",
			},
		}
		expectedReadBuilder := []*fieldMask.SAInt32_builder{
			{
				Id:   0,
				Data: "SHOULD BE SET - 0",
			},
			{
				Id:   1,
				Data: "SHOULD SET - 1",
			},
			{
				Id:   2,
				Data: "SHOULD SET - 2",
			},
			{
				Id:   3,
				Data: "",
			},
		}
		toCreate := saInt32Build(toCreateBuilder)
		expectedRead := saInt32Build(expectedReadBuilder)
		err = sqliteSAInt32CreateSuccessWithReadAfterCheck(opts, repoImpl, context.Background(), toCreate, expectedRead)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}
	}
}

func TestSAInt32Repository_Create_WithNoFieldMaskUsedSetsAllAttributes(t *testing.T) {
	opts := saInt32DefaultCmpOpts()

	for repoDesc, componentUnderTest := range saInt32ImplementationsToTest() {
		// Call setup function, inject t *testing.T, and use t.Cleanup
		repoImpl := componentUnderTest(t)
		if repoImpl == nil {
			t.Fatalf(
				"%s: no implementation provided",
				repoDesc,
			)
		}

		res, err := repoImpl.Read(context.Background(), nil)
		if err != nil {
			t.Fatalf(
				"%s: initial condition: empty repo: %s",
				repoDesc,
				err,
			)
		}
		if len(res) != 0 {
			t.Fatalf(
				"%s: initial condition: empty repo: got %d items",
				repoDesc,
				len(res),
			)
		}
		toCreateBuilder := []*fieldMask.SAInt32_builder{
			{
				Id:   0,
				Data: "0",
			},
			{
				Id:   1,
				Data: "1",
			},
			{
				Id:   2,
				Data: "2",
			},
			{
				Id:   3,
				Data: "3",
			},
		}
		toCreate := saInt32Build(toCreateBuilder)
		err = sqliteSAInt32CreateSuccessWithReadAfterCheck(opts, repoImpl, context.Background(), toCreate, toCreate)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}
	}
}

func TestSAInt32Repository_Update_WithAnyPrimeAttributeExcludedByFieldMaskFails(t *testing.T) {
	opts := saInt32DefaultCmpOpts()

	testCases := map[string][]*fieldMask.SAInt32_builder{
		"all fields excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{""},
				},
				Id:   0,
				Data: "SHOULD FAIL",
			},
		},
		"id field exlucded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"data"},
				},
				Id:   0,
				Data: "SHOULD FAIL",
			},
		},
	}

	for repoDesc, componentUnderTest := range saInt32ImplementationsToTest() {
		for testDesc, testCase := range testCases {
			// Call setup function, inject t *testing.T, and use t.Cleanup
			repoImpl := componentUnderTest(t)
			if repoImpl == nil {
				t.Fatalf(
					"%s: no implementation provided",
					repoDesc,
				)
			}

			err := sqliteSAInt32ReadCheck(opts, repoImpl, context.Background(), nil, nil)
			if err != nil {
				t.Fatalf("%s: %s: %s", repoDesc, testDesc, err)
			}

			toCreateBuilder := []*fieldMask.SAInt32_builder{
				{
					Id:   0,
					Data: "0",
				},
				{
					Id:   1,
					Data: "1",
				},
				{
					Id:   2,
					Data: "2",
				},
				{
					Id:   3,
					Data: "3",
				},
			}
			toCreate := saInt32Build(toCreateBuilder)
			err = sqliteSAInt32CreateSuccessWithReadAfterCheck(opts, repoImpl, context.Background(), toCreate, toCreate)
			if err != nil {
				t.Fatalf("%s: %s: %s", repoDesc, testDesc, err)
			}

			toUpdate := saInt32Build(testCase)
			_, err = repoImpl.Update(context.Background(), toUpdate)
			if err == nil {
				t.Fatalf(
					"%s: %s: Update: expected prime attribute masked error",
					repoDesc,
					testDesc,
				)
			}

			err = sqliteSAInt32ReadCheck(opts, repoImpl, context.Background(), nil, toCreate)
			if err != nil {
				t.Fatalf("%s: %s: %s", repoDesc, testDesc, err)
			}
		}

	}
}

func TestSAInt32Repository_Update_WithNonPrimeAttributesExcludedByFieldMaskAreNotModified(t *testing.T) {
	opts := saInt32DefaultCmpOpts()

	for repoDesc, componentUnderTest := range saInt32ImplementationsToTest() {
		// Call setup function, inject t *testing.T, and use t.Cleanup
		repoImpl := componentUnderTest(t)
		if repoImpl == nil {
			t.Fatalf(
				"%s: no implementation provided",
				repoDesc,
			)
		}

		err := sqliteSAInt32ReadCheck(opts, repoImpl, context.Background(), nil, nil)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}
		toCreateBuilder := []*fieldMask.SAInt32_builder{
			{
				Id:   0,
				Data: "0",
			},
			{
				Id:   1,
				Data: "1",
			},
			{
				Id:   2,
				Data: "2",
			},
			{
				Id:   3,
				Data: "3",
			},
		}
		toCreate := saInt32Build(toCreateBuilder)
		err = sqliteSAInt32CreateSuccessWithReadAfterCheck(opts, repoImpl, context.Background(), toCreate, toCreate)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}

		toUpdateBuilder := []*fieldMask.SAInt32_builder{
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id"},
				},
				Id:   0,
				Data: "SHOULD NOT UPDATE",
			},
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id", "data"},
				},
				Id:   1,
				Data: "SHOULD UPDATE",
			},
		}
		toUpdate := saInt32Build(toUpdateBuilder)
		_, err = repoImpl.Update(context.Background(), toUpdate)
		if err != nil {
			t.Fatalf(
				"%s: Update: %s",
				repoDesc,
				err,
			)
		}

		expectedBuilder := []*fieldMask.SAInt32_builder{
			{
				Id:   0,
				Data: "0",
			},
			{
				Id:   1,
				Data: "SHOULD UPDATE",
			},
			{
				Id:   2,
				Data: "2",
			},
			{
				Id:   3,
				Data: "3",
			},
		}
		expected := saInt32Build(expectedBuilder)
		err = sqliteSAInt32ReadCheck(opts, repoImpl, context.Background(), nil, expected)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}
	}
}

func TestSAInt32Repository_Update_WithNoFieldMaskUsedModifiesAllNonPrimeAttributes(t *testing.T) {
	opts := saInt32DefaultCmpOpts()

	for repoDesc, componentUnderTest := range saInt32ImplementationsToTest() {
		// Call setup function, inject t *testing.T, and use t.Cleanup
		repoImpl := componentUnderTest(t)
		if repoImpl == nil {
			t.Fatalf(
				"%s: no implementation provided",
				repoDesc,
			)
		}

		err := sqliteSAInt32ReadCheck(opts, repoImpl, context.Background(), nil, nil)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}
		toCreateBuilder := []*fieldMask.SAInt32_builder{
			{
				Id:   0,
				Data: "0",
			},
			{
				Id:   1,
				Data: "1",
			},
			{
				Id:   2,
				Data: "2",
			},
			{
				Id:   3,
				Data: "3",
			},
		}
		toCreate := saInt32Build(toCreateBuilder)
		err = sqliteSAInt32CreateSuccessWithReadAfterCheck(opts, repoImpl, context.Background(), toCreate, toCreate)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}

		toUpdateBuilder := []*fieldMask.SAInt32_builder{
			{
				Id:   0,
				Data: "SHOULD UPDATE",
			},
			{
				Id:   1,
				Data: "SHOULD UPDATE",
			},
		}
		toUpdate := saInt32Build(toUpdateBuilder)
		_, err = repoImpl.Update(context.Background(), toUpdate)
		if err != nil {
			t.Fatalf(
				"%s: Update: %s",
				repoDesc,
				err,
			)
		}

		expectedBuilder := []*fieldMask.SAInt32_builder{
			{
				Id:   0,
				Data: "SHOULD UPDATE",
			},
			{
				Id:   1,
				Data: "SHOULD UPDATE",
			},
			{
				Id:   2,
				Data: "2",
			},
			{
				Id:   3,
				Data: "3",
			},
		}
		expected := saInt32Build(expectedBuilder)
		err = sqliteSAInt32ReadCheck(opts, repoImpl, context.Background(), nil, expected)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}
	}
}

func saInt32Build(in []*fieldMask.SAInt32_builder) []*fieldMask.SAInt32 {
	out := make([]*fieldMask.SAInt32, 0, len(in))
	for _, builder := range in {
		out = append(out, builder.Build())
	}
	return out
}

func saInt32ImplementationsToTest() map[string]saInt32ComponentUnderTest {
	return map[string]saInt32ComponentUnderTest{
		"SQLite": sqliteSAInt32ComponentUnderTest,
		"PgSQL":  pgsqlSAInt32ComponentUnderTest,
	}
}

func saInt32DefaultCmpOpts() cmp.Options {
	return cmp.Options{
		cmpopts.IgnoreUnexported(fieldMask.SAInt32{}),
		cmpopts.SortSlices(func(x, y *fieldMask.SAInt32) bool {
			return x.GetId() > y.GetId()
		}),
	}
}
