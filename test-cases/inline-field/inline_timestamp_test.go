package inline_field_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	inline_field "github.com/samlitowitz/protoc-gen-crud/test-cases/inline-field"

	test_cases "github.com/samlitowitz/protoc-gen-crud/test-cases"

	"github.com/samlitowitz/protoc-gen-crud/options"

	"github.com/samlitowitz/expressions"

	sqliteLib "modernc.org/sqlite/lib"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestInlineTimestampRepository_Create_WithADuplicatePrimaryKeyFails(t *testing.T) {
	opts := inlineTimestampDefaultCmpOpts()

	for repoType, componentUnderTest := range inlineTimestampImplementationsToTest() {
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
		initialBuilder := []*inline_field.InlineTimestamp_builder{
			{
				Id:        0,
				Timestamp: timestamppb.New(time.Now()),
			},
			{
				Id:        1,
				Timestamp: timestamppb.New(time.Now()),
			},
			{
				Id:        2,
				Timestamp: timestamppb.New(time.Now()),
			},
			{
				Id:        3,
				Timestamp: timestamppb.New(time.Now()),
			},
		}
		initial := make([]*inline_field.InlineTimestamp, 0, len(initialBuilder))
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

		duplicatesBuilder := []*inline_field.InlineTimestamp_builder{
			{
				Id:        0,
				Timestamp: timestamppb.New(time.Now()),
			},
			{
				Id:        1,
				Timestamp: timestamppb.New(time.Now()),
			},
			{
				Id:        2,
				Timestamp: timestamppb.New(time.Now()),
			},
			{
				Id:        3,
				Timestamp: timestamppb.New(time.Now()),
			},
		}
		duplicates := make([]*inline_field.InlineTimestamp, 0, len(duplicatesBuilder))
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

func TestInlineTimestamp_DescriptorRepository_Create_WithANonDuplicatePrimaryKeySucceeds(t *testing.T) {
	opts := inlineTimestampDefaultCmpOpts()

	for repoType, componentUnderTest := range inlineTimestampImplementationsToTest() {
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
		expectedBuilder := []*inline_field.InlineTimestamp_builder{
			{
				Id:        0,
				Timestamp: timestamppb.New(time.Now()),
			},
			{
				Id:        1,
				Timestamp: timestamppb.New(time.Now()),
			},
			{
				Id:        2,
				Timestamp: timestamppb.New(time.Now()),
			},
			{
				Id:        3,
				Timestamp: timestamppb.New(time.Now()),
			},
		}
		expected := make([]*inline_field.InlineTimestamp, 0, len(expectedBuilder))
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

func TestInlineTimestamp_DescriptorRepository_Update_WithUnLocatablePrimaryKeyUpdatesNothing(t *testing.T) {
	for repoType, componentUnderTest := range inlineTimestampImplementationsToTest() {
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

		expectedBuilder := []*inline_field.InlineTimestamp_builder{
			{
				Id:        0,
				Timestamp: timestamppb.New(time.Now()),
			},
			{
				Id:        1,
				Timestamp: timestamppb.New(time.Now()),
			},
			{
				Id:        2,
				Timestamp: timestamppb.New(time.Now()),
			},
			{
				Id:        3,
				Timestamp: timestamppb.New(time.Now()),
			},
		}
		expected := make([]*inline_field.InlineTimestamp, 0, len(expectedBuilder))
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

func TestInlineTimestamp_DescriptorRepository_Update_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := inlineTimestampDefaultCmpOpts()

	for repoType, componentUnderTest := range inlineTimestampImplementationsToTest() {
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
		initialBuilder := []*inline_field.InlineTimestamp_builder{
			{
				Id:        0,
				Timestamp: timestamppb.New(time.Now()),
			},
			{
				Id:        1,
				Timestamp: timestamppb.New(time.Now()),
			},
			{
				Id:        2,
				Timestamp: timestamppb.New(time.Now()),
			},
			{
				Id:        3,
				Timestamp: timestamppb.New(time.Now()),
			},
		}
		initial := make([]*inline_field.InlineTimestamp, 0, len(initialBuilder))
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

		expected := make([]*inline_field.InlineTimestamp, 0, len(initial))
		for _, inlineTimestamp := range initial {
			expected = append(
				expected,
				inline_field.InlineTimestamp_builder{
					Id:        inlineTimestamp.GetId(),
					Timestamp: timestamppb.New(time.Now().Add(time.Hour * 5)),
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

func TestInlineTimestamp_DescriptorRepository_Delete_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := inlineTimestampDefaultCmpOpts()

	now := timestamppb.New(time.Now())
	one_hour_ago := timestamppb.New(time.Now().Add(time.Hour * -1))
	two_hour_ago := timestamppb.New(time.Now().Add(time.Hour * -2))

	testCases := map[string]struct {
		initial          []*inline_field.InlineTimestamp_builder
		deleteExpression expressions.Expression
		expected         []*inline_field.InlineTimestamp_builder
	}{
		"using primary key": {
			initial: []*inline_field.InlineTimestamp_builder{
				{
					Id:        0,
					Timestamp: now,
				},
				{
					Id:        1,
					Timestamp: now,
				},
				{
					Id:        2,
					Timestamp: now,
				},
				{
					Id:        3,
					Timestamp: now,
				},
			},
			deleteExpression: expressions.NewOr(
				expressions.NewEquals(
					expressions.NewIdentifier(inline_field.InlineTimestamp_Id_Field),
					expressions.NewScalar(2),
				),
				expressions.NewEquals(
					expressions.NewIdentifier(inline_field.InlineTimestamp_Id_Field),
					expressions.NewScalar(3),
				),
			),
			expected: []*inline_field.InlineTimestamp_builder{
				{
					Id:        0,
					Timestamp: now,
				},
				{
					Id:        1,
					Timestamp: now,
				},
			},
		},
		"using non-prime attributes": {
			initial: []*inline_field.InlineTimestamp_builder{
				{
					Id:        0,
					Timestamp: now,
				},
				{
					Id:        1,
					Timestamp: now,
				},
				{
					Id:        2,
					Timestamp: one_hour_ago,
				},
				{
					Id:        3,
					Timestamp: two_hour_ago,
				},
			},
			deleteExpression: expressions.NewOr(
				expressions.NewEquals(
					expressions.NewIdentifier(inline_field.InlineTimestamp_Timestamp_Seconds_Field),
					expressions.NewScalar(one_hour_ago.GetSeconds()),
				),
				expressions.NewEquals(
					expressions.NewIdentifier(inline_field.InlineTimestamp_Timestamp_Seconds_Field),
					expressions.NewScalar(two_hour_ago.GetSeconds()),
				),
			),
			expected: []*inline_field.InlineTimestamp_builder{
				{
					Id:        0,
					Timestamp: now,
				},
				{
					Id:        1,
					Timestamp: now,
				},
			},
		},
	}

	for testDesc, testCase := range testCases {
		for repoType, componentUnderTest := range inlineTimestampImplementationsToTest() {
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
			initial := make([]*inline_field.InlineTimestamp, 0, len(testCase.initial))
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
			expected := make([]*inline_field.InlineTimestamp, 0, len(testCase.expected))
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

func inlineTimestampImplementationsToTest() map[options.Implementation]inlineTimestampComponentUnderTest {
	return map[options.Implementation]inlineTimestampComponentUnderTest{
		options.Implementation_IMPLEMENTATION_SQLITE: sqliteInlineTimestampComponentUnderTest,
		options.Implementation_IMPLEMENTATION_PGSQL:  pgsqlInlineTimestampComponentUnderTest,
	}
}

func inlineTimestampDefaultCmpOpts() cmp.Options {
	return cmp.Options{
		cmpopts.IgnoreUnexported(inline_field.InlineTimestamp{}),
		cmpopts.IgnoreUnexported(timestamppb.Timestamp{}),
		cmpopts.SortSlices(func(x, y *inline_field.InlineTimestamp) bool {
			return x.GetId() < y.GetId()
		}),
	}
}
