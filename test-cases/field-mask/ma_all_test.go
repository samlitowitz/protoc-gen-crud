package field_mask_test

import (
	"context"
	"strings"
	"testing"

	"google.golang.org/genproto/protobuf/field_mask"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	fieldMask "github.com/samlitowitz/protoc-gen-crud/test-cases/field-mask"
)

func TestMAAllRepository_Create_WithAnyPrimeAttributeExcludedByFieldMaskFails(t *testing.T) {
	opts := maAllDefaultCmpOpts()

	testCases := map[string][]*fieldMask.MAAll_builder{
		"all fields excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{""},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD FAIL",
			},
		},
		"id_enum excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_int32", "id_int64", "id_uint32", "id_uint64", "id_string"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD FAIL",
			},
		},
		"id_int32 excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_enum", "id_int64", "id_uint32", "id_uint64", "id_string"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD FAIL",
			},
		},
		"id_int64 excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_enum", "id_int32", "id_uint32", "id_uint64", "id_string"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD FAIL",
			},
		},
		"id_uint32 excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_enum", "id_int32", "id_int64", "id_uint64", "id_string"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD FAIL",
			},
		},
		"id_uint64 excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_enum", "id_int32", "id_int64", "id_uint32", "id_string"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD FAIL",
			},
		},
		"id_string excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_enum", "id_int32", "id_int64", "id_uint32", "id_uint64"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD FAIL",
			},
		},
	}

	for repoDesc, componentUnderTest := range maAllImplementationsToTest() {
		for testDesc, testCase := range testCases {
			// Call setup function, inject t *testing.T, and use t.Cleanup
			repoImpl := componentUnderTest(t)
			if repoImpl == nil {
				t.Fatalf(
					"%s: no implementation provided",
					repoDesc,
				)
			}

			err := sqliteMAAllReadCheck(opts, repoImpl, context.Background(), nil, nil)
			if err != nil {
				t.Fatalf("%s: %s", repoDesc, err)
			}
			input := maAllBuild(testCase)
			_, err = repoImpl.Create(context.Background(), input)
			if err == nil {
				t.Fatalf(
					"%s: %s: Create: expected field masking prime attribute error",
					repoDesc,
					testDesc,
				)
			}

			err = sqliteMAAllReadCheck(opts, repoImpl, context.Background(), nil, nil)
			if err != nil {
				t.Fatalf("%s: %s: %s", repoDesc, testDesc, err)
			}
		}

	}
}

func TestMAAllRepository_Create_WithNonPrimeAttributesExcludedByFieldMaskUsesEmptyValues(t *testing.T) {
	opts := maAllDefaultCmpOpts()

	for repoDesc, componentUnderTest := range maAllImplementationsToTest() {
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
		toCreateBuilder := []*fieldMask.MAAll_builder{
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD SET - 0",
			},
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_enum", "id_int32", "id_int64", "id_uint32", "id_uint64", "id_string", "data"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: fieldMask.PrimaryKeyEnum_ONE.String(),
				Data:     "SHOULD SET - 1",
			},
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_enum", "id_int32", "id_int64", "id_uint32", "id_uint64", "id_string", "data"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: fieldMask.PrimaryKeyEnum_TWO.String(),
				Data:     "SHOULD SET - 2",
			},
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_enum", "id_int32", "id_int64", "id_uint32", "id_uint64", "id_string"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: fieldMask.PrimaryKeyEnum_THREE.String(),
				Data:     "SHOULD BE EMPTY",
			},
		}
		expectedReadBuilder := []*fieldMask.MAAll_builder{
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD SET - 0",
			},
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: fieldMask.PrimaryKeyEnum_ONE.String(),
				Data:     "SHOULD SET - 1",
			},
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: fieldMask.PrimaryKeyEnum_TWO.String(),
				Data:     "SHOULD SET - 2",
			},
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: fieldMask.PrimaryKeyEnum_THREE.String(),
				Data:     "",
			},
		}
		toCreate := maAllBuild(toCreateBuilder)
		expectedRead := maAllBuild(expectedReadBuilder)
		err = sqliteMAAllCreateSuccessWithReadAfterCheck(opts, repoImpl, context.Background(), toCreate, expectedRead)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}
	}
}

func TestMAAllRepository_Create_WithNoFieldMaskUsedSetsAllAttributes(t *testing.T) {
	opts := maAllDefaultCmpOpts()

	for repoDesc, componentUnderTest := range maAllImplementationsToTest() {
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
		toCreateBuilder := []*fieldMask.MAAll_builder{
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD SET - 0",
			},
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: fieldMask.PrimaryKeyEnum_ONE.String(),
				Data:     "SHOULD SET - 1",
			},
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: fieldMask.PrimaryKeyEnum_TWO.String(),
				Data:     "SHOULD SET - 2",
			},
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: fieldMask.PrimaryKeyEnum_THREE.String(),
				Data:     "SHOULD SET - 3",
			},
		}
		toCreate := maAllBuild(toCreateBuilder)
		err = sqliteMAAllCreateSuccessWithReadAfterCheck(opts, repoImpl, context.Background(), toCreate, toCreate)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}
	}
}

func TestMAAllRepository_Update_WithAnyPrimeAttributeExcludedByFieldMaskFails(t *testing.T) {
	opts := maAllDefaultCmpOpts()

	testCases := map[string][]*fieldMask.MAAll_builder{
		"all fields excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{""},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD FAIL",
			},
		},
		"id_enum excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_int32", "id_int64", "id_uint32", "id_uint64", "id_string"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD FAIL",
			},
		},
		"id_int32 excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_enum", "id_int64", "id_uint32", "id_uint64", "id_string"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD FAIL",
			},
		},
		"id_int64 excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_enum", "id_int32", "id_uint32", "id_uint64", "id_string"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD FAIL",
			},
		},
		"id_uint32 excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_enum", "id_int32", "id_int64", "id_uint64", "id_string"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD FAIL",
			},
		},
		"id_uint64 excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_enum", "id_int32", "id_int64", "id_uint32", "id_string"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD FAIL",
			},
		},
		"id_string excluded by field mask": {
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_enum", "id_int32", "id_int64", "id_uint32", "id_uint64"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD FAIL",
			},
		},
	}

	for repoDesc, componentUnderTest := range maAllImplementationsToTest() {
		for testDesc, testCase := range testCases {
			// Call setup function, inject t *testing.T, and use t.Cleanup
			repoImpl := componentUnderTest(t)
			if repoImpl == nil {
				t.Fatalf(
					"%s: no implementation provided",
					repoDesc,
				)
			}

			err := sqliteMAAllReadCheck(opts, repoImpl, context.Background(), nil, nil)
			if err != nil {
				t.Fatalf("%s: %s: %s", repoDesc, testDesc, err)
			}

			toCreateBuilder := []*fieldMask.MAAll_builder{
				{
					IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
					IdInt32:  0,
					IdInt64:  0,
					IdUint32: 0,
					IdUint64: 0,
					IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
					Data:     "SHOULD SET - 0",
				},
				{
					IdEnum:   fieldMask.PrimaryKeyEnum_ONE,
					IdInt32:  1,
					IdInt64:  1,
					IdUint32: 1,
					IdUint64: 1,
					IdString: fieldMask.PrimaryKeyEnum_ONE.String(),
					Data:     "SHOULD SET - 1",
				},
				{
					IdEnum:   fieldMask.PrimaryKeyEnum_TWO,
					IdInt32:  2,
					IdInt64:  2,
					IdUint32: 2,
					IdUint64: 2,
					IdString: fieldMask.PrimaryKeyEnum_TWO.String(),
					Data:     "SHOULD SET - 2",
				},
				{
					IdEnum:   fieldMask.PrimaryKeyEnum_THREE,
					IdInt32:  3,
					IdInt64:  3,
					IdUint32: 3,
					IdUint64: 3,
					IdString: fieldMask.PrimaryKeyEnum_THREE.String(),
					Data:     "SHOULD SET - 3",
				},
			}
			toCreate := maAllBuild(toCreateBuilder)
			err = sqliteMAAllCreateSuccessWithReadAfterCheck(opts, repoImpl, context.Background(), toCreate, toCreate)
			if err != nil {
				t.Fatalf("%s: %s: %s", repoDesc, testDesc, err)
			}
			update := maAllBuild(testCase)
			_, err = repoImpl.Update(context.Background(), update)
			if err == nil {
				t.Fatalf(
					"%s: %s: Update: expected prime attribute masked error",
					repoDesc,
					testDesc,
				)
			}

			err = sqliteMAAllReadCheck(opts, repoImpl, context.Background(), nil, toCreate)
			if err != nil {
				t.Fatalf("%s: %s: %s", repoDesc, testDesc, err)
			}
		}

	}
}

func TestMAAllRepository_Update_WithNonPrimeAttributesExcludedByFieldMaskAreNotModified(t *testing.T) {
	opts := maAllDefaultCmpOpts()

	for repoDesc, componentUnderTest := range maAllImplementationsToTest() {
		// Call setup function, inject t *testing.T, and use t.Cleanup
		repoImpl := componentUnderTest(t)
		if repoImpl == nil {
			t.Fatalf(
				"%s: no implementation provided",
				repoDesc,
			)
		}

		err := sqliteMAAllReadCheck(opts, repoImpl, context.Background(), nil, nil)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}
		toCreateBuilder := []*fieldMask.MAAll_builder{
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD SET - 0",
			},
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: fieldMask.PrimaryKeyEnum_ONE.String(),
				Data:     "SHOULD SET - 1",
			},
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: fieldMask.PrimaryKeyEnum_TWO.String(),
				Data:     "SHOULD SET - 2",
			},
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: fieldMask.PrimaryKeyEnum_THREE.String(),
				Data:     "SHOULD SET - 3",
			},
		}
		toCreate := maAllBuild(toCreateBuilder)
		err = sqliteMAAllCreateSuccessWithReadAfterCheck(opts, repoImpl, context.Background(), toCreate, toCreate)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}

		toUpdateBuilder := []*fieldMask.MAAll_builder{
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_enum", "id_int32", "id_int64", "id_uint32", "id_uint64", "id_string"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD NOT UPDATE",
			},
			{
				FieldMask: &field_mask.FieldMask{
					Paths: []string{"id_enum", "id_int32", "id_int64", "id_uint32", "id_uint64", "id_string", "data"},
				},
				IdEnum:   fieldMask.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: fieldMask.PrimaryKeyEnum_ONE.String(),
				Data:     "SHOULD UPDATE",
			},
		}
		toUpdate := maAllBuild(toUpdateBuilder)
		_, err = repoImpl.Update(context.Background(), toUpdate)
		if err != nil {
			t.Fatalf(
				"%s: Update: %s",
				repoDesc,
				err,
			)
		}

		expectedBuilder := []*fieldMask.MAAll_builder{
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD SET - 0",
			},
			{

				IdEnum:   fieldMask.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: fieldMask.PrimaryKeyEnum_ONE.String(),
				Data:     "SHOULD UPDATE",
			},
			{

				IdEnum:   fieldMask.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: fieldMask.PrimaryKeyEnum_TWO.String(),
				Data:     "SHOULD SET - 2",
			},
			{

				IdEnum:   fieldMask.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: fieldMask.PrimaryKeyEnum_THREE.String(),
				Data:     "SHOULD SET - 3",
			},
		}
		expected := maAllBuild(expectedBuilder)
		err = sqliteMAAllReadCheck(opts, repoImpl, context.Background(), nil, expected)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}
	}
}

func TestMAAllRepository_Update_WithNoFieldMaskUsedModifiesAllNonPrimeAttributes(t *testing.T) {
	opts := maAllDefaultCmpOpts()

	for repoDesc, componentUnderTest := range maAllImplementationsToTest() {
		// Call setup function, inject t *testing.T, and use t.Cleanup
		repoImpl := componentUnderTest(t)
		if repoImpl == nil {
			t.Fatalf(
				"%s: no implementation provided",
				repoDesc,
			)
		}

		err := sqliteMAAllReadCheck(opts, repoImpl, context.Background(), nil, nil)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}
		toCreateBuilder := []*fieldMask.MAAll_builder{
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD SET - 0",
			},
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: fieldMask.PrimaryKeyEnum_ONE.String(),
				Data:     "SHOULD SET - 1",
			},
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: fieldMask.PrimaryKeyEnum_TWO.String(),
				Data:     "SHOULD SET - 2",
			},
			{

				IdEnum:   fieldMask.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: fieldMask.PrimaryKeyEnum_THREE.String(),
				Data:     "SHOULD SET - 3",
			},
		}
		toCreate := maAllBuild(toCreateBuilder)
		err = sqliteMAAllCreateSuccessWithReadAfterCheck(opts, repoImpl, context.Background(), toCreate, toCreate)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}

		toUpdateBuilder := []*fieldMask.MAAll_builder{
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD UPDATE - 0",
			},
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: fieldMask.PrimaryKeyEnum_ONE.String(),
				Data:     "SHOULD UPDATE - 1",
			},
		}
		toUpdate := maAllBuild(toUpdateBuilder)
		_, err = repoImpl.Update(context.Background(), toUpdate)
		if err != nil {
			t.Fatalf(
				"%s: Update: %s",
				repoDesc,
				err,
			)
		}

		expectedBuilder := []*fieldMask.MAAll_builder{
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: fieldMask.PrimaryKeyEnum_ZERO.String(),
				Data:     "SHOULD UPDATE - 0",
			},
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: fieldMask.PrimaryKeyEnum_ONE.String(),
				Data:     "SHOULD UPDATE - 1",
			},
			{
				IdEnum:   fieldMask.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: fieldMask.PrimaryKeyEnum_TWO.String(),
				Data:     "SHOULD SET - 2",
			},
			{

				IdEnum:   fieldMask.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: fieldMask.PrimaryKeyEnum_THREE.String(),
				Data:     "SHOULD SET - 3",
			},
		}
		expected := maAllBuild(expectedBuilder)
		err = sqliteMAAllReadCheck(opts, repoImpl, context.Background(), nil, expected)
		if err != nil {
			t.Fatalf("%s: %s", repoDesc, err)
		}
	}
}

func maAllBuild(in []*fieldMask.MAAll_builder) []*fieldMask.MAAll {
	out := make([]*fieldMask.MAAll, 0, len(in))
	for _, builder := range in {
		out = append(out, builder.Build())
	}
	return out
}

func maAllImplementationsToTest() map[string]maAllComponentUnderTest {
	return map[string]maAllComponentUnderTest{
		"SQLite": sqliteMAAllComponentUnderTest,
		"PgSQL":  pgsqlMAAllComponentUnderTest,
	}
}

func maAllDefaultCmpOpts() cmp.Options {
	return cmp.Options{
		cmpopts.IgnoreUnexported(fieldMask.MAAll{}),
		cmpopts.SortSlices(func(x, y *fieldMask.MAAll) bool {
			switch strings.Compare(x.GetIdEnum().String(), y.GetIdEnum().String()) {
			case -1:
				return true
			case 1:
				return false
			}

			if x.GetIdInt32() > y.GetIdInt32() {
				return false
			}
			if x.GetIdInt32() < y.GetIdInt32() {
				return true
			}

			if x.GetIdInt64() > y.GetIdInt64() {
				return false
			}
			if x.GetIdInt64() < y.GetIdInt64() {
				return true
			}

			if x.GetIdUint32() > y.GetIdUint32() {
				return false
			}
			if x.GetIdUint32() < y.GetIdUint32() {
				return true
			}

			if x.GetIdUint64() > y.GetIdUint64() {
				return false
			}
			if x.GetIdUint64() < y.GetIdUint64() {
				return true
			}

			switch strings.Compare(x.GetIdString(), y.GetIdString()) {
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
}
