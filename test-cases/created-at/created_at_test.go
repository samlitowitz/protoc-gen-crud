package created_at_test

import (
	"context"
	"fmt"
	"testing"

	created_at "github.com/samlitowitz/protoc-gen-crud/test-cases/created-at"

	"google.golang.org/protobuf/types/known/timestamppb"

	test_cases "github.com/samlitowitz/protoc-gen-crud/test-cases"

	"github.com/samlitowitz/protoc-gen-crud/options"

	"github.com/samlitowitz/protoc-gen-crud/expressions"

	sqliteLib "modernc.org/sqlite/lib"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreatedAtRepository_Create_WithADuplicatePrimaryKeyFails(t *testing.T) {
	opts := createdAtDefaultCmpOpts()
	for repoType, componentUnderTest := range createdAtImplementationsToTest() {
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
		initialBuilder := []*created_at.CreatedAt_builder{
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
		initial := make([]*created_at.CreatedAt, 0, len(initialBuilder))
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

		duplicatesBuilder := []*created_at.CreatedAt_builder{
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
		duplicates := make([]*created_at.CreatedAt, 0, len(duplicatesBuilder))
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

func TestCreatedAtRepository_Create_WithANonDuplicatePrimaryKeySucceeds(t *testing.T) {
	opts := createdAtDefaultCmpOpts()
	for repoType, componentUnderTest := range createdAtImplementationsToTest() {
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
		expectedBuilder := []*created_at.CreatedAt_builder{
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
		expected := make([]*created_at.CreatedAt, 0, len(expectedBuilder))
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

func TestCreatedAtRepository_Update_WithUnLocatablePrimaryKeyUpdatesNothing(t *testing.T) {
	for repoType, componentUnderTest := range createdAtImplementationsToTest() {
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

		expectedBuilder := []*created_at.CreatedAt_builder{
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
		expected := make([]*created_at.CreatedAt, 0, len(expectedBuilder))
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

func TestCreatedAtRepository_Update_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := createdAtDefaultCmpOpts()
	for repoType, componentUnderTest := range createdAtImplementationsToTest() {
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
		initialBuilder := []*created_at.CreatedAt_builder{
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
		initial := make([]*created_at.CreatedAt, 0, len(initialBuilder))
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

		expectedBuilder := make([]*created_at.CreatedAt_builder, 0, len(initial))
		for _, inlineTimestamp := range initial {
			expectedBuilder = append(
				expectedBuilder,
				&created_at.CreatedAt_builder{
					Id:        inlineTimestamp.GetId(),
					Data:      inlineTimestamp.GetData() + " - updated",
					CreatedAt: inlineTimestamp.GetCreatedAt(),
				},
			)
		}
		expected := make([]*created_at.CreatedAt, 0, len(expectedBuilder))
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

func TestCreatedAtRepository_Delete_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := createdAtDefaultCmpOpts()
	testCases := map[string]struct {
		initial          []*created_at.CreatedAt_builder
		deleteExpression expressions.Expression
		expected         []*created_at.CreatedAt_builder
	}{
		"using primary key": {
			initial: []*created_at.CreatedAt_builder{
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
					expressions.NewIdentifier(created_at.CreatedAt_Id_Field),
					expressions.NewScalar(2),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(created_at.CreatedAt_Id_Field),
					expressions.NewScalar(3),
				),
			),
			expected: []*created_at.CreatedAt_builder{
				{
					Id:   0,
					Data: "zero",
				},
				{
					Id:   1,
					Data: "one",
				},
			},
		},
		"using non-prime attributes": {
			initial: []*created_at.CreatedAt_builder{
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
					expressions.NewIdentifier(created_at.CreatedAt_Data_Field),
					expressions.NewScalar("two"),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(created_at.CreatedAt_Data_Field),
					expressions.NewScalar("three"),
				),
			),
			expected: []*created_at.CreatedAt_builder{
				{
					Id:   0,
					Data: "zero",
				},
				{
					Id:   1,
					Data: "one",
				},
			},
		},
	}

	for testDesc, testCase := range testCases {
		for repoType, componentUnderTest := range createdAtImplementationsToTest() {
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
			initial := make([]*created_at.CreatedAt, 0, len(testCase.initial))
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
			expected := make([]*created_at.CreatedAt, 0, len(testCase.expected))
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

func createdAtImplementationsToTest() map[options.Implementation]createdAtComponentUnderTest {
	return map[options.Implementation]createdAtComponentUnderTest{
		options.Implementation_IMPLEMENTATION_SQLITE: sqliteCreatedAtComponentUnderTest,
		options.Implementation_IMPLEMENTATION_PGSQL:  pgsqlCreatedAtComponentUnderTest,
	}
}

func createdAtDefaultCmpOpts() cmp.Options {
	return cmp.Options{
		cmpopts.IgnoreUnexported(created_at.CreatedAt{}),
		cmpopts.IgnoreUnexported(timestamppb.Timestamp{}),
		cmpopts.IgnoreFields(timestamppb.Timestamp{}, "Nanos"),
		cmpopts.SortSlices(func(x, y *created_at.CreatedAt) bool {
			return x.GetId() < y.GetId()
		}),
	}
}
