package relationships_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/samlitowitz/protoc-gen-crud/expressions"

	"modernc.org/sqlite"
	sqliteLib "modernc.org/sqlite/lib"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/samlitowitz/protoc-gen-crud/test-cases/relationships"
)

func TestSAInt32Repository_Create_WithADuplicatePrimaryKeyFails(t *testing.T) {
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
		initial := []*relationships.SAInt32{
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

		duplicates := []*relationships.SAInt32{
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

func TestSAInt32Repository_Create_WithANonDuplicatePrimaryKeySucceeds(t *testing.T) {
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
		expected := []*relationships.SAInt32{
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

func TestSAInt32Repository_Update_WithUnLocatablePrimaryKeyUpdatesNothing(t *testing.T) {
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

		expected := []*relationships.SAInt32{
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

func TestSAInt32Repository_Update_WithLocatablePrimaryKeySucceeds(t *testing.T) {
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
		initial := []*relationships.SAInt32{
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

		expected := make([]*relationships.SAInt32, 0, len(initial))
		for _, saint32 := range initial {
			expected = append(
				expected,
				&relationships.SAInt32{
					Id:   saint32.GetId(),
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

func TestSAInt32Repository_Delete_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := saInt32DefaultCmpOpts()

	testCases := map[string]struct {
		initial          []*relationships.SAInt32
		deleteExpression expressions.Expression
		expected         []*relationships.SAInt32
	}{
		"using primary key": {
			initial: []*relationships.SAInt32{
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
			},
			deleteExpression: expressions.NewOr(
				expressions.NewEqual(
					expressions.NewIdentifier(relationships.Saint32_Id_Field),
					expressions.NewScalar(2),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(relationships.Saint32_Id_Field),
					expressions.NewScalar(3),
				),
			),
			expected: []*relationships.SAInt32{
				{
					Id:   0,
					Data: "0",
				},
				{
					Id:   1,
					Data: "1",
				},
			},
		},
		"using non-prime attributes": {
			initial: []*relationships.SAInt32{
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
			},
			deleteExpression: expressions.NewOr(
				expressions.NewEqual(
					expressions.NewIdentifier(relationships.Saint32_Data_Field),
					expressions.NewScalar("2"),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(relationships.Saint32_Data_Field),
					expressions.NewScalar("3"),
				),
			),
			expected: []*relationships.SAInt32{
				{
					Id:   0,
					Data: "0",
				},
				{
					Id:   1,
					Data: "1",
				},
			},
		},
	}

	for testDesc, testCase := range testCases {
		for repoDesc, componentUnderTest := range saInt32ImplementationsToTest() {
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

func saInt32ImplementationsToTest() map[string]saInt32ComponentUnderTest {
	return map[string]saInt32ComponentUnderTest{
		"SQLite": sqliteSAInt32ComponentUnderTest,
	}
}

func saInt32DefaultCmpOpts() cmp.Options {
	return cmp.Options{
		cmpopts.IgnoreUnexported(relationships.SAInt32{}),
		cmpopts.SortSlices(func(x, y *relationships.SAInt32) bool {
			return x.GetId() > y.GetId()
		}),
	}
}
