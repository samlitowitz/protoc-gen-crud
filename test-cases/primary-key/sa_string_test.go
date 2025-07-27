package primary_key_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	test_cases "github.com/samlitowitz/protoc-gen-crud/test-cases"

	"github.com/samlitowitz/protoc-gen-crud/options"

	"github.com/samlitowitz/protoc-gen-crud/expressions"

	sqliteLib "modernc.org/sqlite/lib"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	primaryKey "github.com/samlitowitz/protoc-gen-crud/test-cases/primary-key"
)

func TestSAStringRepository_Create_WithADuplicatePrimaryKeyFails(t *testing.T) {
	opts := saStringDefaultCmpOpts()

	for repoType, componentUnderTest := range saStringImplementationsToTest() {
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
		initialBuilder := []*primaryKey.SAString_builder{
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
		initial := make([]*primaryKey.SAString, 0, len(initialBuilder))
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

		duplicatesBuilder := []*primaryKey.SAString_builder{
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
		duplicates := make([]*primaryKey.SAString, 0, len(duplicatesBuilder))
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

func TestSAStringRepository_Create_WithANonDuplicatePrimaryKeySucceeds(t *testing.T) {
	opts := saStringDefaultCmpOpts()

	for repoType, componentUnderTest := range saStringImplementationsToTest() {
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
		expectedBuilder := []*primaryKey.SAString_builder{
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
		expected := make([]*primaryKey.SAString, 0, len(expectedBuilder))
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

func TestSAStringRepository_Update_WithUnLocatablePrimaryKeyUpdatesNothing(t *testing.T) {
	for repoType, componentUnderTest := range saStringImplementationsToTest() {
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

		expectedBuilder := []*primaryKey.SAString_builder{
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
		expected := make([]*primaryKey.SAString, 0, len(expectedBuilder))
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

func TestSAStringRepository_Update_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := saStringDefaultCmpOpts()

	for repoType, componentUnderTest := range saStringImplementationsToTest() {
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
		initialBuilder := []*primaryKey.SAString_builder{
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
		initial := make([]*primaryKey.SAString, 0, len(initialBuilder))
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

		expected := make([]*primaryKey.SAString, 0, len(initial))
		for _, sastring := range initial {
			expected = append(
				expected,
				primaryKey.SAString_builder{
					Id:   sastring.GetId(),
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

func TestSAStringRepository_Delete_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := saStringDefaultCmpOpts()

	testCases := map[string]struct {
		initial          []*primaryKey.SAString_builder
		deleteExpression expressions.Expression
		expected         []*primaryKey.SAString_builder
	}{
		"using primary key": {
			initial: []*primaryKey.SAString_builder{
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
			expected: []*primaryKey.SAString_builder{
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
			initial: []*primaryKey.SAString_builder{
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
			expected: []*primaryKey.SAString_builder{
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
		for repoType, componentUnderTest := range saStringImplementationsToTest() {
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
			initial := make([]*primaryKey.SAString, 0, len(testCase.initial))
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
			expected := saStringBuild(testCase.expected)
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

func saStringBuild(in []*primaryKey.SAString_builder) []*primaryKey.SAString {
	out := make([]*primaryKey.SAString, 0, len(in))
	for _, builder := range in {
		out = append(out, builder.Build())
	}
	return out
}

func saStringImplementationsToTest() map[options.Implementation]saStringComponentUnderTest {
	return map[options.Implementation]saStringComponentUnderTest{
		options.Implementation_IMPLEMENTATION_SQLITE: sqliteSAStringComponentUnderTest,
		options.Implementation_IMPLEMENTATION_PGSQL:  pgsqlSAStringComponentUnderTest,
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
