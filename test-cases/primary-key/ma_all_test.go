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

func TestMAAllRepository_Create_WithADuplicatePrimaryKeyFails(t *testing.T) {
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
		initial := []*primaryKey.MAAll{
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: "zero",
				Data:     "zero",
			},
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: "one",
				Data:     "one",
			},
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: "two",
				Data:     "two",
			},
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: "three",
				Data:     "three",
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

		duplicates := []*primaryKey.MAAll{
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: "zero",
				Data:     "zero",
			},
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: "one",
				Data:     "one",
			},
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: "two",
				Data:     "two",
			},
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: "three",
				Data:     "three",
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

func TestMAAllRepository_Create_WithANonDuplicatePrimaryKeySucceeds(t *testing.T) {
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
		expected := []*primaryKey.MAAll{
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: "zero",
				Data:     "zero",
			},
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: "one",
				Data:     "one",
			},
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: "two",
				Data:     "two",
			},
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: "three",
				Data:     "three",
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

func TestMAAllRepository_Update_WithUnLocatablePrimaryKeyUpdatesNothing(t *testing.T) {
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

		expected := []*primaryKey.MAAll{
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: "zero",
				Data:     "zero",
			},
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: "one",
				Data:     "one",
			},
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: "two",
				Data:     "two",
			},
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: "three",
				Data:     "three",
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

func TestMAAllRepository_Update_WithLocatablePrimaryKeySucceeds(t *testing.T) {
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
		initial := []*primaryKey.MAAll{
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: "zero",
				Data:     "zero",
			},
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: "one",
				Data:     "one",
			},
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: "two",
				Data:     "two",
			},
			{
				IdEnum:   primaryKey.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: "three",
				Data:     "three",
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

		expected := make([]*primaryKey.MAAll, 0, len(initial))
		for _, maall := range initial {
			expected = append(
				expected,
				&primaryKey.MAAll{
					IdEnum:   maall.GetIdEnum(),
					IdInt32:  maall.GetIdInt32(),
					IdInt64:  maall.GetIdInt64(),
					IdUint32: maall.GetIdUint32(),
					IdUint64: maall.GetIdUint64(),
					IdString: maall.GetIdString(),
					Data:     "UPDATED",
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

func TestMAAllRepository_Delete_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := maAllDefaultCmpOpts()

	testCases := map[string]struct {
		initial          []*primaryKey.MAAll
		deleteExpression expressions.Expression
		expected         []*primaryKey.MAAll
	}{
		"using primary key": {
			initial: []*primaryKey.MAAll{
				{
					IdEnum:   primaryKey.PrimaryKeyEnum_ZERO,
					IdInt32:  0,
					IdInt64:  0,
					IdUint32: 0,
					IdUint64: 0,
					IdString: "zero",
					Data:     "zero",
				},
				{
					IdEnum:   primaryKey.PrimaryKeyEnum_ONE,
					IdInt32:  1,
					IdInt64:  1,
					IdUint32: 1,
					IdUint64: 1,
					IdString: "one",
					Data:     "one",
				},
				{
					IdEnum:   primaryKey.PrimaryKeyEnum_TWO,
					IdInt32:  2,
					IdInt64:  2,
					IdUint32: 2,
					IdUint64: 2,
					IdString: "two",
					Data:     "two",
				},
				{
					IdEnum:   primaryKey.PrimaryKeyEnum_THREE,
					IdInt32:  3,
					IdInt64:  3,
					IdUint32: 3,
					IdUint64: 3,
					IdString: "three",
					Data:     "three",
				},
			},
			deleteExpression: expressions.NewOr(
				expressions.NewEqual(
					expressions.NewIdentifier(primaryKey.Maall_IdEnum_Field),
					expressions.NewScalar(primaryKey.PrimaryKeyEnum_TWO),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(primaryKey.Maall_IdEnum_Field),
					expressions.NewScalar(primaryKey.PrimaryKeyEnum_THREE),
				),
			),
			expected: []*primaryKey.MAAll{
				{
					IdEnum:   primaryKey.PrimaryKeyEnum_ZERO,
					IdInt32:  0,
					IdInt64:  0,
					IdUint32: 0,
					IdUint64: 0,
					IdString: "zero",
					Data:     "zero",
				},
				{
					IdEnum:   primaryKey.PrimaryKeyEnum_ONE,
					IdInt32:  1,
					IdInt64:  1,
					IdUint32: 1,
					IdUint64: 1,
					IdString: "one",
					Data:     "one",
				},
			},
		},
		"using non-prime attributes": {
			initial: []*primaryKey.MAAll{
				{
					IdEnum:   primaryKey.PrimaryKeyEnum_ZERO,
					IdInt32:  0,
					IdInt64:  0,
					IdUint32: 0,
					IdUint64: 0,
					IdString: "zero",
					Data:     "zero",
				},
				{
					IdEnum:   primaryKey.PrimaryKeyEnum_ONE,
					IdInt32:  1,
					IdInt64:  1,
					IdUint32: 1,
					IdUint64: 1,
					IdString: "one",
					Data:     "one",
				},
				{
					IdEnum:   primaryKey.PrimaryKeyEnum_TWO,
					IdInt32:  2,
					IdInt64:  2,
					IdUint32: 2,
					IdUint64: 2,
					IdString: "two",
					Data:     "two",
				},
				{
					IdEnum:   primaryKey.PrimaryKeyEnum_THREE,
					IdInt32:  3,
					IdInt64:  3,
					IdUint32: 3,
					IdUint64: 3,
					IdString: "three",
					Data:     "three",
				},
			},
			deleteExpression: expressions.NewOr(
				expressions.NewEqual(
					expressions.NewIdentifier(primaryKey.Maall_Data_Field),
					expressions.NewScalar("two"),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(primaryKey.Maall_Data_Field),
					expressions.NewScalar("three"),
				),
			),
			expected: []*primaryKey.MAAll{
				{
					IdEnum:   primaryKey.PrimaryKeyEnum_ZERO,
					IdInt32:  0,
					IdInt64:  0,
					IdUint32: 0,
					IdUint64: 0,
					IdString: "zero",
					Data:     "zero",
				},
				{
					IdEnum:   primaryKey.PrimaryKeyEnum_ONE,
					IdInt32:  1,
					IdInt64:  1,
					IdUint32: 1,
					IdUint64: 1,
					IdString: "one",
					Data:     "one",
				},
			},
		},
	}

	for testDesc, testCase := range testCases {
		for repoDesc, componentUnderTest := range maAllImplementationsToTest() {
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

			//res = res[:0]
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

func maAllImplementationsToTest() map[string]maAllComponentUnderTest {
	return map[string]maAllComponentUnderTest{
		"SQLite": sqliteMAAllComponentUnderTest,
	}
}

func maAllDefaultCmpOpts() cmp.Options {
	return cmp.Options{
		cmpopts.IgnoreUnexported(primaryKey.MAAll{}),
		cmpopts.SortSlices(func(x, y *primaryKey.MAAll) bool {
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
