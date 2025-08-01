package updated_at_test

import (
	"context"
	"fmt"
	"testing"

	updated_at "github.com/samlitowitz/protoc-gen-crud/test-cases/updated-at"

	"google.golang.org/protobuf/types/known/timestamppb"

	test_cases "github.com/samlitowitz/protoc-gen-crud/test-cases"

	"github.com/samlitowitz/protoc-gen-crud/options"

	"github.com/samlitowitz/protoc-gen-crud/expressions"

	sqliteLib "modernc.org/sqlite/lib"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestUpdatedAtRepository_Create_WithADuplicatePrimaryKeyFails(t *testing.T) {
	opts := updatedAtDefaultCmpOpts()

	for repoType, componentUnderTest := range updatedAtImplementationsToTest() {
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
		initialBuilder := []*updated_at.UpdatedAt_builder{
			{
				Id:   0,
				Data: "zero",
			},
			{
				Id:   1,
				Data: "one",
			},
			{
				Id:   2,
				Data: "two",
			},
			{
				Id:   3,
				Data: "three",
			},
		}
		initial := make([]*updated_at.UpdatedAt, 0, len(initialBuilder))
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

		duplicatesBuilder := []*updated_at.UpdatedAt_builder{
			{
				Id:   0,
				Data: "zero",
			},
			{
				Id:   1,
				Data: "one",
			},
			{
				Id:   2,
				Data: "two",
			},
			{
				Id:   3,
				Data: "three",
			},
		}
		duplicates := make([]*updated_at.UpdatedAt, 0, len(duplicatesBuilder))
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

func TestUpdatedAtRepository_Create_WithANonDuplicatePrimaryKeySucceeds(t *testing.T) {
	opts := updatedAtDefaultCmpOpts()

	for repoType, componentUnderTest := range updatedAtImplementationsToTest() {
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
		expectedBuilder := []*updated_at.UpdatedAt_builder{
			{
				Id:   0,
				Data: "zero",
			},
			{
				Id:   1,
				Data: "one",
			},
			{
				Id:   2,
				Data: "two",
			},
			{
				Id:   3,
				Data: "three",
			},
		}
		expected := make([]*updated_at.UpdatedAt, 0, len(expectedBuilder))
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

func TestUpdatedAtRepository_Update_WithUnLocatablePrimaryKeyUpdatesNothing(t *testing.T) {
	for repoType, componentUnderTest := range updatedAtImplementationsToTest() {
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

		expectedBuilder := []*updated_at.UpdatedAt_builder{
			{
				Id:   0,
				Data: "zero",
			},
			{
				Id:   1,
				Data: "one",
			},
			{
				Id:   2,
				Data: "two",
			},
			{
				Id:   3,
				Data: "three",
			},
		}
		expected := make([]*updated_at.UpdatedAt, 0, len(expectedBuilder))
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

func TestUpdatedAtRepository_Update_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := updatedAtDefaultCmpOpts()

	for repoType, componentUnderTest := range updatedAtImplementationsToTest() {
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
		initialBuilder := []*updated_at.UpdatedAt_builder{
			{
				Id:   0,
				Data: "zero",
			},
			{
				Id:   1,
				Data: "one",
			},
			{
				Id:   2,
				Data: "two",
			},
			{
				Id:   3,
				Data: "three",
			},
		}
		initial := make([]*updated_at.UpdatedAt, 0, len(initialBuilder))
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

		expectedBuilder := make([]*updated_at.UpdatedAt_builder, 0, len(initial))
		for _, inlineTimestamp := range initial {
			expectedBuilder = append(
				expectedBuilder,
				&updated_at.UpdatedAt_builder{
					Id: inlineTimestamp.GetId(),
				},
			)
		}
		expected := make([]*updated_at.UpdatedAt, 0, len(expectedBuilder))
		for _, builder := range expectedBuilder {
			expected = append(expected, builder.Build())
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

func TestUpdatedAtRepository_Delete_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := updatedAtDefaultCmpOpts()

	testCases := map[string]struct {
		initial          []*updated_at.UpdatedAt_builder
		deleteExpression expressions.Expression
		expected         []*updated_at.UpdatedAt_builder
	}{
		"using primary key": {
			initial: []*updated_at.UpdatedAt_builder{
				{
					Id: 0,
				},
				{
					Id: 1,
				},
				{
					Id: 2,
				},
				{
					Id: 3,
				},
			},
			deleteExpression: expressions.NewOr(
				expressions.NewEqual(
					expressions.NewIdentifier(updated_at.UpdatedAt_Id_Field),
					expressions.NewScalar(2),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(updated_at.UpdatedAt_Id_Field),
					expressions.NewScalar(3),
				),
			),
			expected: []*updated_at.UpdatedAt_builder{
				{
					Id: 0,
				},
				{
					Id: 1,
				},
			},
		},
		"using non-prime attributes": {
			initial: []*updated_at.UpdatedAt_builder{
				{
					Id:   0,
					Data: "zero",
				},
				{
					Id:   1,
					Data: "one",
				},
				{
					Id:   2,
					Data: "two",
				},
				{
					Id:   3,
					Data: "three",
				},
			},
			deleteExpression: expressions.NewOr(
				expressions.NewEqual(
					expressions.NewIdentifier(updated_at.UpdatedAt_Id_Field),
					expressions.NewScalar(2),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(updated_at.UpdatedAt_Id_Field),
					expressions.NewScalar(3),
				),
			),
			expected: []*updated_at.UpdatedAt_builder{
				{
					Id: 0,
				},
				{
					Id: 1,
				},
			},
		},
	}

	for testDesc, testCase := range testCases {
		for repoType, componentUnderTest := range updatedAtImplementationsToTest() {
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
			initial := make([]*updated_at.UpdatedAt, 0, len(testCase.initial))
			for _, builder := range testCase.initial {
				initial = append(initial, builder.Build())
			}
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
							"%s: %s: Read():",
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
					"%s: %s: Delete(): %s",
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
			expected := make([]*updated_at.UpdatedAt, 0, len(testCase.expected))
			for _, builder := range testCase.expected {
				expected = append(expected, builder.Build())
			}
			if diff := cmp.Diff(expected, res, opts); diff != "" {
				t.Fatal(
					mismatch(
						fmt.Sprintf(
							"%s: %s: Read():",
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

func TestUpdatedAtRepository_Create_CorrectlySetsUpdatedAt(t *testing.T) {
	opts := updatedAtDefaultCmpOpts()

	for repoType, componentUnderTest := range updatedAtImplementationsToTest() {
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
		expectedBuilder := []*updated_at.UpdatedAt_builder{
			{
				Id:   0,
				Data: "zero",
			},
			{
				Id:   1,
				Data: "one",
			},
			{
				Id:   2,
				Data: "two",
			},
			{
				Id:   3,
				Data: "three",
			},
		}
		expected := make([]*updated_at.UpdatedAt, 0, len(expectedBuilder))
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

func updatedAtImplementationsToTest() map[options.Implementation]updatedAtComponentUnderTest {
	return map[options.Implementation]updatedAtComponentUnderTest{
		options.Implementation_IMPLEMENTATION_SQLITE: sqliteUpdatedAtComponentUnderTest,
		options.Implementation_IMPLEMENTATION_PGSQL:  pgsqlUpdatedAtComponentUnderTest,
	}
}

func updatedAtDefaultCmpOpts() cmp.Options {
	return cmp.Options{
		cmpopts.IgnoreUnexported(updated_at.UpdatedAt{}),
		cmpopts.IgnoreUnexported(timestamppb.Timestamp{}),
		cmpopts.IgnoreFields(timestamppb.Timestamp{}, "Nanos"),
		cmpopts.SortSlices(func(x, y *updated_at.UpdatedAt) bool {
			return x.GetId() < y.GetId()
		}),
	}
}
