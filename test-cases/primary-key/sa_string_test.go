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

func TestSAStringRepository_Create_WithADuplicatePrimaryKeyFails(t *testing.T) {
	opts := saStringDefaultCmpOpts()

	for repoDesc, componentUnderTest := range saStringImplementationsToTest() {
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
		initial := []*primaryKey.SAString{
			{
				Id:   "zero",
				Data: "zero",
			},
			{
				Id:   "one",
				Data: "one",
			},
			{
				Id:   "two",
				Data: "two",
			},
			{
				Id:   "three",
				Data: "three",
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

		duplicates := []*primaryKey.SAString{
			{
				Id:   "zero",
				Data: "zero",
			},
			{
				Id:   "one",
				Data: "one",
			},
			{
				Id:   "two",
				Data: "two",
			},
			{
				Id:   "three",
				Data: "three",
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

func TestSAStringRepository_Create_WithANonDuplicatePrimaryKeySucceeds(t *testing.T) {
	opts := saStringDefaultCmpOpts()

	for repoDesc, componentUnderTest := range saStringImplementationsToTest() {
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
		expected := []*primaryKey.SAString{
			{
				Id:   "zero",
				Data: "zero",
			},
			{
				Id:   "one",
				Data: "one",
			},
			{
				Id:   "two",
				Data: "two",
			},
			{
				Id:   "three",
				Data: "three",
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

func TestSAStringRepository_Update_WithUnLocatablePrimaryKeyUpdatesNothing(t *testing.T) {
	for repoDesc, componentUnderTest := range saStringImplementationsToTest() {
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

		expected := []*primaryKey.SAString{
			{
				Id:   "zero",
				Data: "zero",
			},
			{
				Id:   "one",
				Data: "one",
			},
			{
				Id:   "two",
				Data: "two",
			},
			{
				Id:   "three",
				Data: "three",
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

func TestSAStringRepository_Update_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := saStringDefaultCmpOpts()

	for repoDesc, componentUnderTest := range saStringImplementationsToTest() {
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
		initial := []*primaryKey.SAString{
			{
				Id:   "zero",
				Data: "zero",
			},
			{
				Id:   "one",
				Data: "one",
			},
			{
				Id:   "two",
				Data: "two",
			},
			{
				Id:   "three",
				Data: "three",
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

		expected := make([]*primaryKey.SAString, 0, len(initial))
		for _, sastring := range initial {
			expected = append(
				expected,
				&primaryKey.SAString{
					Id:   sastring.GetId(),
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

func TestSAStringRepository_Delete_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := saStringDefaultCmpOpts()

	testCases := map[string]struct {
		initial          []*primaryKey.SAString
		deleteExpression expressions.Expression
		expected         []*primaryKey.SAString
	}{
		"using primary key": {
			initial: []*primaryKey.SAString{
				{
					Id:   "zero",
					Data: "zero",
				},
				{
					Id:   "one",
					Data: "one",
				},
				{
					Id:   "two",
					Data: "two",
				},
				{
					Id:   "three",
					Data: "three",
				},
			},
			deleteExpression: expressions.NewOr(
				expressions.NewEqual(
					expressions.NewIdentifier(primaryKey.Sastring_Id_Field),
					expressions.NewScalar("two"),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(primaryKey.Sastring_Id_Field),
					expressions.NewScalar("three"),
				),
			),
			expected: []*primaryKey.SAString{
				{
					Id:   "zero",
					Data: "zero",
				},
				{
					Id:   "one",
					Data: "one",
				},
			},
		},
		"using non-prime attributes": {
			initial: []*primaryKey.SAString{
				{
					Id:   "zero",
					Data: "zero",
				},
				{
					Id:   "one",
					Data: "one",
				},
				{
					Id:   "two",
					Data: "two",
				},
				{
					Id:   "three",
					Data: "three",
				},
			},
			deleteExpression: expressions.NewOr(
				expressions.NewEqual(
					expressions.NewIdentifier(primaryKey.Sastring_Data_Field),
					expressions.NewScalar("two"),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(primaryKey.Sastring_Data_Field),
					expressions.NewScalar("three"),
				),
			),
			expected: []*primaryKey.SAString{
				{
					Id:   "zero",
					Data: "zero",
				},
				{
					Id:   "one",
					Data: "one",
				},
			},
		},
	}

	for testDesc, testCase := range testCases {
		for repoDesc, componentUnderTest := range saStringImplementationsToTest() {
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

func saStringImplementationsToTest() map[string]saStringComponentUnderTest {
	return map[string]saStringComponentUnderTest{
		"SQLite": sqliteSAStringComponentUnderTest,
	}
}

func saStringDefaultCmpOpts() cmp.Options {
	return cmp.Options{
		cmpopts.IgnoreUnexported(primaryKey.SAString{}),
		cmpopts.SortSlices(func(x, y *primaryKey.SAString) bool {
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
}
