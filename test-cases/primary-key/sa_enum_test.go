package primary_key_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/samlitowitz/protoc-gen-crud/expressions"

	"modernc.org/sqlite"
	sqliteLib "modernc.org/sqlite/lib"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	primaryKey "github.com/samlitowitz/protoc-gen-crud/test-cases/primary-key"
)

func TestSAEnumRepository_Create_WithADuplicatePrimaryKeyFails(t *testing.T) {
	opts := saEnumDefaultCmpOpts()

	for repoDesc, componentUnderTest := range saEnumImplementationsToTest() {
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
		initial := []*primaryKey.SAEnum{
			{
				Id:   primaryKey.PrimaryKeyEnum_ZERO,
				Data: primaryKey.PrimaryKeyEnum_ZERO.String(),
			},
			{
				Id:   primaryKey.PrimaryKeyEnum_ONE,
				Data: primaryKey.PrimaryKeyEnum_ONE.String(),
			},
			{
				Id:   primaryKey.PrimaryKeyEnum_TWO,
				Data: primaryKey.PrimaryKeyEnum_TWO.String(),
			},
			{
				Id:   primaryKey.PrimaryKeyEnum_THREE,
				Data: primaryKey.PrimaryKeyEnum_THREE.String(),
			},
		}
		res, err = repoImpl.Create(context.Background(), initial)
		if err != nil {
			t.Fatalf(
				"%s: Create(): %s",
				repoDesc,
				err,
			)
		}

		if diff := cmp.Diff(initial, res, opts); diff != "" {
			t.Fatal(
				mismatch(
					fmt.Sprintf(
						"%s: Create():",
						repoDesc,
					),
					diff,
				),
			)
		}

		duplicates := []*primaryKey.SAEnum{
			{
				Id:   primaryKey.PrimaryKeyEnum_ZERO,
				Data: primaryKey.PrimaryKeyEnum_ZERO.String(),
			},
			{
				Id:   primaryKey.PrimaryKeyEnum_ONE,
				Data: primaryKey.PrimaryKeyEnum_ONE.String(),
			},
			{
				Id:   primaryKey.PrimaryKeyEnum_TWO,
				Data: primaryKey.PrimaryKeyEnum_TWO.String(),
			},
			{
				Id:   primaryKey.PrimaryKeyEnum_THREE,
				Data: primaryKey.PrimaryKeyEnum_THREE.String(),
			},
		}
		res, err = repoImpl.Create(context.Background(), duplicates)
		if err == nil {
			t.Fatalf("%s: Create(): expected error", repoDesc)
		}

		sqlErr, ok := err.(*sqlite.Error)
		if !ok {
			t.Fatalf("%s: Create(): expected *sqlite.Error, got %T", repoDesc, err)
		}
		if sqlErr.Code() != sqliteLib.SQLITE_CONSTRAINT_PRIMARYKEY {
			t.Fatalf("%s: Create(): expected duplicate error code, got %d", repoDesc, sqlErr.Code())
		}

		res, err = repoImpl.Read(context.Background(), nil)
		if err != nil {
			t.Fatalf(
				"%s: initial condition: empty repo: %s",
				repoDesc,
				err,
			)
		}
		if diff := cmp.Diff(initial, res, opts); diff != "" {
			t.Fatal(
				mismatch(
					fmt.Sprintf(
						"%s: Create():",
						repoDesc,
					),
					diff,
				),
			)
		}
	}
}

func TestSAEnumRepository_Create_WithANonDuplicatePrimaryKeySucceeds(t *testing.T) {
	opts := saEnumDefaultCmpOpts()

	for repoDesc, componentUnderTest := range saEnumImplementationsToTest() {
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
		expected := []*primaryKey.SAEnum{
			{
				Id:   primaryKey.PrimaryKeyEnum_ZERO,
				Data: primaryKey.PrimaryKeyEnum_ZERO.String(),
			},
			{
				Id:   primaryKey.PrimaryKeyEnum_ONE,
				Data: primaryKey.PrimaryKeyEnum_ONE.String(),
			},
			{
				Id:   primaryKey.PrimaryKeyEnum_TWO,
				Data: primaryKey.PrimaryKeyEnum_TWO.String(),
			},
			{
				Id:   primaryKey.PrimaryKeyEnum_THREE,
				Data: primaryKey.PrimaryKeyEnum_THREE.String(),
			},
		}
		res, err = repoImpl.Create(context.Background(), expected)
		if err != nil {
			t.Fatalf(
				"%s: Create(): %s",
				repoDesc,
				err,
			)
		}

		if diff := cmp.Diff(expected, res, opts); diff != "" {
			t.Fatal(
				mismatch(
					fmt.Sprintf(
						"%s: Create():",
						repoDesc,
					),
					diff,
				),
			)
		}

		res, err = repoImpl.Read(context.Background(), nil)
		if err != nil {
			t.Fatalf(
				"%s: initial condition: empty repo: %s",
				repoDesc,
				err,
			)
		}
		if diff := cmp.Diff(expected, res, opts); diff != "" {
			t.Fatal(
				mismatch(
					fmt.Sprintf(
						"%s: Create():",
						repoDesc,
					),
					diff,
				),
			)
		}
	}
}

func TestSAEnumRepository_Update_WithUnLocatablePrimaryKeyUpdatesNothing(t *testing.T) {
	for repoDesc, componentUnderTest := range saEnumImplementationsToTest() {
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

		expected := []*primaryKey.SAEnum{
			{
				Id:   primaryKey.PrimaryKeyEnum_ZERO,
				Data: primaryKey.PrimaryKeyEnum_ZERO.String(),
			},
			{
				Id:   primaryKey.PrimaryKeyEnum_ONE,
				Data: primaryKey.PrimaryKeyEnum_ONE.String(),
			},
			{
				Id:   primaryKey.PrimaryKeyEnum_TWO,
				Data: primaryKey.PrimaryKeyEnum_TWO.String(),
			},
			{
				Id:   primaryKey.PrimaryKeyEnum_THREE,
				Data: primaryKey.PrimaryKeyEnum_THREE.String(),
			},
		}

		_, err = repoImpl.Update(context.Background(), expected)
		if err != nil {
			t.Fatalf(
				"%s: Update(): %s",
				repoDesc,
				err,
			)
		}

		res, err = repoImpl.Read(context.Background(), nil)
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
	}
}

func TestSAEnumRepository_Update_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := saEnumDefaultCmpOpts()

	for repoDesc, componentUnderTest := range saEnumImplementationsToTest() {
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
		initial := []*primaryKey.SAEnum{
			{
				Id:   primaryKey.PrimaryKeyEnum_ZERO,
				Data: primaryKey.PrimaryKeyEnum_ZERO.String(),
			},
			{
				Id:   primaryKey.PrimaryKeyEnum_ONE,
				Data: primaryKey.PrimaryKeyEnum_ONE.String(),
			},
			{
				Id:   primaryKey.PrimaryKeyEnum_TWO,
				Data: primaryKey.PrimaryKeyEnum_TWO.String(),
			},
			{
				Id:   primaryKey.PrimaryKeyEnum_THREE,
				Data: primaryKey.PrimaryKeyEnum_THREE.String(),
			},
		}
		res, err = repoImpl.Create(context.Background(), initial)
		if err != nil {
			t.Fatalf(
				"%s: Create(): %s",
				repoDesc,
				err,
			)
		}
		if diff := cmp.Diff(initial, res, opts); diff != "" {
			t.Fatal(
				mismatch(
					fmt.Sprintf(
						"%s: Create():",
						repoDesc,
					),
					diff,
				),
			)
		}

		res, err = repoImpl.Read(context.Background(), nil)
		if err != nil {
			t.Fatalf(
				"%s: initial condition: empty repo: %s",
				repoDesc,
				err,
			)
		}
		if diff := cmp.Diff(initial, res, opts); diff != "" {
			t.Fatal(
				mismatch(
					fmt.Sprintf(
						"%s: Create():",
						repoDesc,
					),
					diff,
				),
			)
		}

		expected := make([]*primaryKey.SAEnum, 0, len(initial))
		for _, saenum := range initial {
			expected = append(
				expected,
				&primaryKey.SAEnum{
					Id:   saenum.GetId(),
					Data: "UPDATED",
				},
			)
		}

		res, err = repoImpl.Update(context.Background(), expected)
		if err != nil {
			t.Fatalf(
				"%s: Update(): %s",
				repoDesc,
				err,
			)
		}
		if diff := cmp.Diff(expected, res, opts); diff != "" {
			t.Fatal(
				mismatch(
					fmt.Sprintf(
						"%s: Update():",
						repoDesc,
					),
					diff,
				),
			)
		}
	}
}

func TestSAEnumRepository_Delete_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := saEnumDefaultCmpOpts()

	testCases := map[string]struct {
		initial          []*primaryKey.SAEnum
		deleteExpression expressions.Expression
		expected         []*primaryKey.SAEnum
	}{
		"using primary key": {
			initial: []*primaryKey.SAEnum{
				{
					Id:   primaryKey.PrimaryKeyEnum_ZERO,
					Data: primaryKey.PrimaryKeyEnum_ZERO.String(),
				},
				{
					Id:   primaryKey.PrimaryKeyEnum_ONE,
					Data: primaryKey.PrimaryKeyEnum_ONE.String(),
				},
				{
					Id:   primaryKey.PrimaryKeyEnum_TWO,
					Data: primaryKey.PrimaryKeyEnum_TWO.String(),
				},
				{
					Id:   primaryKey.PrimaryKeyEnum_THREE,
					Data: primaryKey.PrimaryKeyEnum_THREE.String(),
				},
			},
			deleteExpression: expressions.NewOr(
				expressions.NewEqual(
					expressions.NewIdentifier(primaryKey.Saenum_Id_Field),
					expressions.NewScalar(primaryKey.PrimaryKeyEnum_TWO),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(primaryKey.Saenum_Id_Field),
					expressions.NewScalar(primaryKey.PrimaryKeyEnum_THREE),
				),
			),
			expected: []*primaryKey.SAEnum{
				{
					Id:   primaryKey.PrimaryKeyEnum_ZERO,
					Data: primaryKey.PrimaryKeyEnum_ZERO.String(),
				},
				{
					Id:   primaryKey.PrimaryKeyEnum_ONE,
					Data: primaryKey.PrimaryKeyEnum_ONE.String(),
				},
			},
		},
		"using non-prime attributes": {
			initial: []*primaryKey.SAEnum{
				{
					Id:   primaryKey.PrimaryKeyEnum_ZERO,
					Data: primaryKey.PrimaryKeyEnum_ZERO.String(),
				},
				{
					Id:   primaryKey.PrimaryKeyEnum_ONE,
					Data: primaryKey.PrimaryKeyEnum_ONE.String(),
				},
				{
					Id:   primaryKey.PrimaryKeyEnum_TWO,
					Data: primaryKey.PrimaryKeyEnum_TWO.String(),
				},
				{
					Id:   primaryKey.PrimaryKeyEnum_THREE,
					Data: primaryKey.PrimaryKeyEnum_THREE.String(),
				},
			},
			deleteExpression: expressions.NewOr(
				expressions.NewEqual(
					expressions.NewIdentifier(primaryKey.Saenum_Data_Field),
					expressions.NewScalar(primaryKey.PrimaryKeyEnum_TWO.String()),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(primaryKey.Saenum_Data_Field),
					expressions.NewScalar(primaryKey.PrimaryKeyEnum_THREE.String()),
				),
			),
			expected: []*primaryKey.SAEnum{
				{
					Id:   primaryKey.PrimaryKeyEnum_ZERO,
					Data: primaryKey.PrimaryKeyEnum_ZERO.String(),
				},
				{
					Id:   primaryKey.PrimaryKeyEnum_ONE,
					Data: primaryKey.PrimaryKeyEnum_ONE.String(),
				},
			},
		},
	}

	for testDesc, testCase := range testCases {
		for repoDesc, componentUnderTest := range saEnumImplementationsToTest() {
			// Call setup function, inject t *testing.T, and use t.Cleanup
			repoImpl := componentUnderTest(t)
			if repoImpl == nil {
				t.Fatalf(
					"%s: %s: no implementation provided",
					testDesc,
					repoDesc,
				)
			}

			res, err := repoImpl.Read(context.Background(), nil)
			if err != nil {
				t.Fatalf(
					"%s: %s: initial condition: empty repo: %s",
					testDesc,
					repoDesc,
					err,
				)
			}
			if len(res) != 0 {
				t.Fatalf(
					"%s: %s: initial condition: empty repo: got %d items",
					testDesc,
					repoDesc,
					len(res),
				)
			}
			res, err = repoImpl.Create(context.Background(), testCase.initial)
			if err != nil {
				t.Fatalf(
					"%s: %s: Create(): %s",
					testDesc,
					repoDesc,
					err,
				)
			}
			if diff := cmp.Diff(testCase.initial, res, opts); diff != "" {
				t.Fatal(
					mismatch(
						fmt.Sprintf(
							"%s: %s: Create():",
							testDesc,
							repoDesc,
						),
						diff,
					),
				)
			}

			res, err = repoImpl.Read(context.Background(), nil)
			if err != nil {
				t.Fatalf(
					"%s: %s: initial condition: %s",
					testDesc,
					repoDesc,
					err,
				)
			}
			if diff := cmp.Diff(testCase.initial, res, opts); diff != "" {
				t.Fatal(
					mismatch(
						fmt.Sprintf(
							"%s: %s: Create():",
							testDesc,
							repoDesc,
						),
						diff,
					),
				)
			}

			err = repoImpl.Delete(context.Background(), testCase.deleteExpression)
			if err != nil {
				t.Fatalf(
					"%s: %s: Update(): %s",
					testDesc,
					repoDesc,
					err,
				)
			}

			res, err = repoImpl.Read(context.Background(), nil)
			if err != nil {
				t.Fatalf(
					"%s: %s: initial condition: %s",
					testDesc,
					repoDesc,
					err,
				)
			}
			if diff := cmp.Diff(testCase.expected, res, opts); diff != "" {
				t.Fatal(
					mismatch(
						fmt.Sprintf(
							"%s: %s: Create():",
							testDesc,
							repoDesc,
						),
						diff,
					),
				)
			}
		}
	}
}

func saEnumImplementationsToTest() map[string]saEnumComponentUnderTest {
	return map[string]saEnumComponentUnderTest{
		"SQLite": sqliteSAEnumComponentUnderTest,
	}
}

func saEnumDefaultCmpOpts() cmp.Options {
	return cmp.Options{
		cmpopts.IgnoreUnexported(primaryKey.SAEnum{}),
		cmpopts.SortSlices(func(x, y *primaryKey.SAEnum) bool {
			switch strings.Compare(x.GetId().String(), y.GetId().String()) {
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
