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

	testCases := map[string][]*fieldMask.SAInt32{
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

			_, err = repoImpl.Create(context.Background(), testCase)
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
		toCreate := []*fieldMask.SAInt32{
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
		expectedRead := []*fieldMask.SAInt32{
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
		toCreate := []*fieldMask.SAInt32{
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
		err = sqliteSAInt32CreateSuccessWithReadAfterCheck(opts, repoImpl, context.Background(), toCreate, toCreate)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}
	}
}

func TestSAInt32Repository_Update_WithAnyPrimeAttributeExcludedByFieldMaskFails(t *testing.T) {
	opts := saInt32DefaultCmpOpts()

	testCases := map[string][]*fieldMask.SAInt32{
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

			toCreate := []*fieldMask.SAInt32{
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
			err = sqliteSAInt32CreateSuccessWithReadAfterCheck(opts, repoImpl, context.Background(), toCreate, toCreate)
			if err != nil {
				t.Fatalf("%s: %s: %s", repoDesc, testDesc, err)
			}

			_, err = repoImpl.Update(context.Background(), testCase)
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
		toCreate := []*fieldMask.SAInt32{
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
		err = sqliteSAInt32CreateSuccessWithReadAfterCheck(opts, repoImpl, context.Background(), toCreate, toCreate)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}

		toUpdate := []*fieldMask.SAInt32{
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
		_, err = repoImpl.Update(context.Background(), toUpdate)
		if err != nil {
			t.Fatalf(
				"%s: Update: %s",
				repoDesc,
				err,
			)
		}

		expected := []*fieldMask.SAInt32{
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
		toCreate := []*fieldMask.SAInt32{
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
		err = sqliteSAInt32CreateSuccessWithReadAfterCheck(opts, repoImpl, context.Background(), toCreate, toCreate)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}

		toUpdate := []*fieldMask.SAInt32{
			{
				Id:   0,
				Data: "SHOULD UPDATE",
			},
			{
				Id:   1,
				Data: "SHOULD UPDATE",
			},
		}
		_, err = repoImpl.Update(context.Background(), toUpdate)
		if err != nil {
			t.Fatalf(
				"%s: Update: %s",
				repoDesc,
				err,
			)
		}

		expected := []*fieldMask.SAInt32{
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
		err = sqliteSAInt32ReadCheck(opts, repoImpl, context.Background(), nil, expected)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}
	}
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
		cmpopts.IgnoreFields(fieldMask.SAInt32{}, "FieldMask"),
		cmpopts.SortSlices(func(x, y *fieldMask.SAInt32) bool {
			return x.GetId() > y.GetId()
		}),
	}
}
