package as_timestamp_field_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	as_timestamp_field "github.com/samlitowitz/protoc-gen-crud/test-cases/as-timestamp-field"

	"google.golang.org/protobuf/types/known/timestamppb"

	test_cases "github.com/samlitowitz/protoc-gen-crud/test-cases"

	"github.com/samlitowitz/protoc-gen-crud/options"

	"github.com/samlitowitz/protoc-gen-crud/expressions"

	sqliteLib "modernc.org/sqlite/lib"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestAsTimestampRepository_Create_WithADuplicatePrimaryKeyFails(t *testing.T) {
	opts := asTimestampDefaultCmpOpts()

	for repoType, componentUnderTest := range asTimestampImplementationsToTest() {
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
		initial := []*as_timestamp_field.AsTimestamp{
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

		duplicates := []*as_timestamp_field.AsTimestamp{
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
		res, err = repoImpl.Create(context.Background(), duplicates)
		if err == nil {
			t.Fatalf("%s: Create(): expected error", repoDesc)
		}

		test_cases.AssertSQLErrorCode(
			t,
			repoType,
			map[options.Implementation]any{
				options.Implementation_PGSQL:  "23505",
				options.Implementation_SQLITE: sqliteLib.SQLITE_CONSTRAINT_PRIMARYKEY,
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

func TestAsTimestamp_DescriptorRepository_Create_WithANonDuplicatePrimaryKeySucceeds(t *testing.T) {
	opts := asTimestampDefaultCmpOpts()

	for repoType, componentUnderTest := range asTimestampImplementationsToTest() {
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
		expected := []*as_timestamp_field.AsTimestamp{
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

func TestAsTimestamp_DescriptorRepository_Update_WithUnLocatablePrimaryKeyUpdatesNothing(t *testing.T) {
	for repoType, componentUnderTest := range asTimestampImplementationsToTest() {
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

		expected := []*as_timestamp_field.AsTimestamp{
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

func TestAsTimestamp_DescriptorRepository_Update_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := asTimestampDefaultCmpOpts()

	for repoType, componentUnderTest := range asTimestampImplementationsToTest() {
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
		initial := []*as_timestamp_field.AsTimestamp{
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

		expected := make([]*as_timestamp_field.AsTimestamp, 0, len(initial))
		for _, inlineTimestamp := range initial {
			expected = append(
				expected,
				&as_timestamp_field.AsTimestamp{
					Id:        inlineTimestamp.GetId(),
					Timestamp: timestamppb.New(time.Now().Add(time.Hour * 5)),
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

func TestAsTimestamp_DescriptorRepository_Delete_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := asTimestampDefaultCmpOpts()

	now := timestamppb.New(time.Now())
	one_hour_ago := timestamppb.New(time.Now().Add(time.Hour * -1))
	two_hour_ago := timestamppb.New(time.Now().Add(time.Hour * -2))

	testCases := map[string]struct {
		initial          []*as_timestamp_field.AsTimestamp
		deleteExpression expressions.Expression
		expected         []*as_timestamp_field.AsTimestamp
	}{
		"using primary key": {
			initial: []*as_timestamp_field.AsTimestamp{
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
				expressions.NewEqual(
					expressions.NewIdentifier(as_timestamp_field.AsTimestamp_Id_Field),
					expressions.NewScalar(2),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(as_timestamp_field.AsTimestamp_Id_Field),
					expressions.NewScalar(3),
				),
			),
			expected: []*as_timestamp_field.AsTimestamp{
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
			initial: []*as_timestamp_field.AsTimestamp{
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
				expressions.NewEqual(
					expressions.NewIdentifier(as_timestamp_field.AsTimestamp_Timestamp_Field),
					expressions.NewScalar(one_hour_ago.GetSeconds()),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(as_timestamp_field.AsTimestamp_Timestamp_Field),
					expressions.NewScalar(two_hour_ago.GetSeconds()),
				),
			),
			expected: []*as_timestamp_field.AsTimestamp{
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
		for repoType, componentUnderTest := range asTimestampImplementationsToTest() {
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
			if diff := cmp.Diff(testCase.expected, res, opts); diff != "" {
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

func asTimestampImplementationsToTest() map[options.Implementation]asTimestampComponentUnderTest {
	return map[options.Implementation]asTimestampComponentUnderTest{
		options.Implementation_SQLITE: sqliteAsTimestampComponentUnderTest,
		options.Implementation_PGSQL:  pgsqlAsTimestampComponentUnderTest,
	}
}

func asTimestampDefaultCmpOpts() cmp.Options {
	return cmp.Options{
		cmpopts.IgnoreUnexported(as_timestamp_field.AsTimestamp{}),
		cmpopts.IgnoreUnexported(timestamppb.Timestamp{}),
		cmpopts.SortSlices(func(x, y *as_timestamp_field.AsTimestamp) bool {
			return x.GetId() < y.GetId()
		}),
	}
}
