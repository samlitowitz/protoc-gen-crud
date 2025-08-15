package primary_key_test

import (
	"context"
	"fmt"
	"testing"

	test_cases "github.com/samlitowitz/protoc-gen-crud/test-cases"

	"github.com/samlitowitz/protoc-gen-crud/options"

	"github.com/samlitowitz/expressions"

	sqliteLib "modernc.org/sqlite/lib"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	primaryKey "github.com/samlitowitz/protoc-gen-crud/test-cases/primary-key"
)

func TestSAInt32Repository_Create_WithADuplicatePrimaryKeyFails(t *testing.T) {
	opts := saInt32DefaultCmpOpts()

	for repoType, componentUnderTest := range saInt32ImplementationsToTest() {
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
		initialBuilder := []*primaryKey.SAInt32_builder{
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
		initial := make([]*primaryKey.SAInt32, 0, len(initialBuilder))
		for _, builder := range initialBuilder {
			initial = append(initial, builder.Build())
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

		duplicatesBuilder := []*primaryKey.SAInt32_builder{
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
		duplicates := make([]*primaryKey.SAInt32, 0, len(duplicatesBuilder))
		for _, builder := range duplicatesBuilder {
			duplicates = append(duplicates, builder.Build())
		}
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

func TestSAInt32Repository_Create_WithANonDuplicatePrimaryKeySucceeds(t *testing.T) {
	opts := saInt32DefaultCmpOpts()

	for repoType, componentUnderTest := range saInt32ImplementationsToTest() {
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
		expectedBuilder := []*primaryKey.SAInt32_builder{
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
		expected := make([]*primaryKey.SAInt32, 0, len(expectedBuilder))
		for _, builder := range expectedBuilder {
			expected = append(expected, builder.Build())
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
	for repoType, componentUnderTest := range saInt32ImplementationsToTest() {
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

		expectedBuilder := []*primaryKey.SAInt32_builder{
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
		expected := make([]*primaryKey.SAInt32, 0, len(expectedBuilder))
		for _, builder := range expectedBuilder {
			expected = append(expected, builder.Build())
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

	for repoType, componentUnderTest := range saInt32ImplementationsToTest() {
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
		initialBuilder := []*primaryKey.SAInt32_builder{
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
		initial := make([]*primaryKey.SAInt32, 0, len(initialBuilder))
		for _, builder := range initialBuilder {
			initial = append(initial, builder.Build())
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

		expected := make([]*primaryKey.SAInt32, 0, len(initial))
		for _, saint32 := range initial {
			expected = append(
				expected,
				primaryKey.SAInt32_builder{
					Id:   saint32.GetId(),
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

func TestSAInt32Repository_Delete_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := saInt32DefaultCmpOpts()

	testCases := map[string]struct {
		initial          []*primaryKey.SAInt32_builder
		deleteExpression expressions.Expression
		expected         []*primaryKey.SAInt32_builder
	}{
		"using primary key": {
			initial: []*primaryKey.SAInt32_builder{
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
				expressions.NewEquals(
					expressions.NewIdentifier(primaryKey.Saint32_Id_Field),
					expressions.NewScalar(2),
				),
				expressions.NewEquals(
					expressions.NewIdentifier(primaryKey.Saint32_Id_Field),
					expressions.NewScalar(3),
				),
			),
			expected: []*primaryKey.SAInt32_builder{
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
			initial: []*primaryKey.SAInt32_builder{
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
				expressions.NewEquals(
					expressions.NewIdentifier(primaryKey.Saint32_Data_Field),
					expressions.NewScalar("2"),
				),
				expressions.NewEquals(
					expressions.NewIdentifier(primaryKey.Saint32_Data_Field),
					expressions.NewScalar("3"),
				),
			),
			expected: []*primaryKey.SAInt32_builder{
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
		for repoType, componentUnderTest := range saInt32ImplementationsToTest() {
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
			initial := saInt32Build(testCase.initial)
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
			expected := saInt32Build(testCase.expected)
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

func saInt32Build(in []*primaryKey.SAInt32_builder) []*primaryKey.SAInt32 {
	out := make([]*primaryKey.SAInt32, 0, len(in))
	for _, builder := range in {
		out = append(out, builder.Build())
	}
	return out
}

func saInt32ImplementationsToTest() map[options.Implementation]saInt32ComponentUnderTest {
	return map[options.Implementation]saInt32ComponentUnderTest{
		options.Implementation_IMPLEMENTATION_SQLITE: sqliteSAInt32ComponentUnderTest,
		options.Implementation_IMPLEMENTATION_PGSQL:  pgsqlSAInt32ComponentUnderTest,
	}
}

func saInt32DefaultCmpOpts() cmp.Options {
	return cmp.Options{
		cmpopts.IgnoreUnexported(primaryKey.SAInt32{}),
		cmpopts.SortSlices(func(x, y *primaryKey.SAInt32) bool {
			return x.GetId() > y.GetId()
		}),
	}
}
