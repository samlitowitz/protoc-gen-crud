package primary_key_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/samlitowitz/protoc-gen-crud/options"
	test_cases "github.com/samlitowitz/protoc-gen-crud/test-cases"

	"github.com/samlitowitz/protoc-gen-crud/expressions"

	sqliteLib "modernc.org/sqlite/lib"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	primaryKey "github.com/samlitowitz/protoc-gen-crud/test-cases/primary-key"
)

func TestSAUint64Repository_Create_WithADuplicatePrimaryKeyFails(t *testing.T) {
	opts := saUint64DefaultCmpOpts()

	for repoType, componentUnderTest := range saUint64ImplementationsToTest() {
		repoDesc := repoType.String()
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
		initialBuilder := []*primaryKey.SAUint64_builder{
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
		initial := saUint64Build(initialBuilder)
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

		duplicatesBuilder := []*primaryKey.SAUint64_builder{
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
		duplicates := saUint64Build(duplicatesBuilder)
		_, err = repoImpl.Create(context.Background(), duplicates)
		if err == nil {
			t.Fatalf("%s: Create(): expected error", repoDesc)
		}

		test_cases.AssertSQLErrorCode(
			t,
			repoType,
			map[options.Implementation]any{
				options.Implementation_IMPLEMENTATION_PGSQL:  "23505",
				options.Implementation_IMPLEMENTATION_SQLITE: sqliteLib.SQLITE_CONSTRAINT_PRIMARYKEY,
			},
			err,
			fmt.Sprintf("%s: Create(): ", repoDesc),
		)

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

func TestSAUint64Repository_Create_WithANonDuplicatePrimaryKeySucceeds(t *testing.T) {
	opts := saUint64DefaultCmpOpts()

	for repoType, componentUnderTest := range saUint64ImplementationsToTest() {
		repoDesc := repoType.String()
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
		expectedBuilder := []*primaryKey.SAUint64_builder{
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
		expected := saUint64Build(expectedBuilder)
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

func TestSAUint64Repository_Update_WithUnLocatablePrimaryKeyUpdatesNothing(t *testing.T) {
	for repoType, componentUnderTest := range saUint64ImplementationsToTest() {
		repoDesc := repoType.String()
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

		expectedBuilder := []*primaryKey.SAUint64_builder{
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
		expected := saUint64Build(expectedBuilder)
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

func TestSAUint64Repository_Update_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := saUint64DefaultCmpOpts()

	for repoType, componentUnderTest := range saUint64ImplementationsToTest() {
		repoDesc := repoType.String()
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
		initialBuilder := []*primaryKey.SAUint64_builder{
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
		initial := saUint64Build(initialBuilder)
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

		expected := make([]*primaryKey.SAUint64, 0, len(initial))
		for _, sauint64 := range initial {
			expected = append(
				expected,
				primaryKey.SAUint64_builder{
					Id:   sauint64.GetId(),
					Data: "UPDATED",
				}.Build(),
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

func TestSAUint64Repository_Delete_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := saUint64DefaultCmpOpts()

	testCases := map[string]struct {
		initial          []*primaryKey.SAUint64_builder
		deleteExpression expressions.Expression
		expected         []*primaryKey.SAUint64_builder
	}{
		"using primary key": {
			initial: []*primaryKey.SAUint64_builder{
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
					expressions.NewIdentifier(primaryKey.Sauint64_Id_Field),
					expressions.NewScalar(2),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(primaryKey.Sauint64_Id_Field),
					expressions.NewScalar(3),
				),
			),
			expected: []*primaryKey.SAUint64_builder{
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
			initial: []*primaryKey.SAUint64_builder{
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
					expressions.NewIdentifier(primaryKey.Sauint64_Data_Field),
					expressions.NewScalar("2"),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(primaryKey.Sauint64_Data_Field),
					expressions.NewScalar("3"),
				),
			),
			expected: []*primaryKey.SAUint64_builder{
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
		for repoType, componentUnderTest := range saUint64ImplementationsToTest() {
			repoDesc := repoType.String()
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
			initial := saUint64Build(testCase.initial)
			res, err = repoImpl.Create(context.Background(), initial)
			if err != nil {
				t.Fatalf(
					"%s: %s: Create(): %s",
					testDesc,
					repoDesc,
					err,
				)
			}
			if diff := cmp.Diff(initial, res, opts); diff != "" {
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
			if diff := cmp.Diff(initial, res, opts); diff != "" {
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
			expected := saUint64Build(testCase.expected)
			if diff := cmp.Diff(expected, res, opts); diff != "" {
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

func saUint64Build(in []*primaryKey.SAUint64_builder) []*primaryKey.SAUint64 {
	out := make([]*primaryKey.SAUint64, 0, len(in))
	for _, builder := range in {
		out = append(out, builder.Build())
	}
	return out
}

func saUint64ImplementationsToTest() map[options.Implementation]saUint64ComponentUnderTest {
	return map[options.Implementation]saUint64ComponentUnderTest{
		options.Implementation_IMPLEMENTATION_SQLITE: sqliteSAUint64ComponentUnderTest,
		options.Implementation_IMPLEMENTATION_PGSQL:  pgsqlSAUint64ComponentUnderTest,
	}
}

func saUint64DefaultCmpOpts() cmp.Options {
	return cmp.Options{
		cmpopts.IgnoreUnexported(primaryKey.SAUint64{}),
		cmpopts.SortSlices(func(x, y *primaryKey.SAUint64) bool {
			return x.GetId() > y.GetId()
		}),
	}
}
