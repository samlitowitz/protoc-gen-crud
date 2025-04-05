package relationships_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	test_cases "github.com/samlitowitz/protoc-gen-crud/test-cases"

	"github.com/samlitowitz/protoc-gen-crud/expressions"
	"github.com/samlitowitz/protoc-gen-crud/options"

	sqliteLib "modernc.org/sqlite/lib"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/samlitowitz/protoc-gen-crud/test-cases/relationships"
)

func TestMAAllRepository_Create_WithADuplicatePrimaryKeyFails(t *testing.T) {
	opts := maAllDefaultCmpOpts()

	for repoType, componentUnderTest := range maAllImplementationsToTest() {
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
		initialBuilder := []*relationships.MAAll_builder{
			{
				IdEnum:   relationships.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: "zero",
				Data:     "zero",
			},
			{
				IdEnum:   relationships.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: "one",
				Data:     "one",
			},
			{
				IdEnum:   relationships.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: "two",
				Data:     "two",
			},
			{
				IdEnum:   relationships.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: "three",
				Data:     "three",
			},
		}
		initial := make([]*relationships.MAAll, 0, len(initialBuilder))
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

		duplicatesBuilder := []*relationships.MAAll_builder{
			{
				IdEnum:   relationships.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: "zero",
				Data:     "zero",
			},
			{
				IdEnum:   relationships.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: "one",
				Data:     "one",
			},
			{
				IdEnum:   relationships.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: "two",
				Data:     "two",
			},
			{
				IdEnum:   relationships.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: "three",
				Data:     "three",
			},
		}
		duplicates := make([]*relationships.MAAll, 0, len(duplicatesBuilder))
		for _, builder := range duplicatesBuilder {
			duplicates = append(duplicates, builder.Build())
		}
		res, err = repoImpl.Create(context.Background(), duplicates)
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

func TestMAAllRepository_Create_WithANonDuplicatePrimaryKeySucceeds(t *testing.T) {
	opts := maAllDefaultCmpOpts()

	for repoType, componentUnderTest := range maAllImplementationsToTest() {
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
		expectedBuilder := []*relationships.MAAll_builder{
			{
				IdEnum:   relationships.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: "zero",
				Data:     "zero",
			},
			{
				IdEnum:   relationships.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: "one",
				Data:     "one",
			},
			{
				IdEnum:   relationships.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: "two",
				Data:     "two",
			},
			{
				IdEnum:   relationships.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: "three",
				Data:     "three",
			},
		}
		expected := make([]*relationships.MAAll, 0, len(expectedBuilder))
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

func TestMAAllRepository_Update_WithUnLocatablePrimaryKeyUpdatesNothing(t *testing.T) {
	for repoType, componentUnderTest := range maAllImplementationsToTest() {
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

		expectedBuilder := []*relationships.MAAll_builder{
			{
				IdEnum:   relationships.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: "zero",
				Data:     "zero",
			},
			{
				IdEnum:   relationships.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: "one",
				Data:     "one",
			},
			{
				IdEnum:   relationships.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: "two",
				Data:     "two",
			},
			{
				IdEnum:   relationships.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: "three",
				Data:     "three",
			},
		}
		expected := make([]*relationships.MAAll, 0, len(expectedBuilder))
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

func TestMAAllRepository_Update_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := maAllDefaultCmpOpts()

	for repoType, componentUnderTest := range maAllImplementationsToTest() {
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
		initialBuilder := []*relationships.MAAll_builder{
			{
				IdEnum:   relationships.PrimaryKeyEnum_ZERO,
				IdInt32:  0,
				IdInt64:  0,
				IdUint32: 0,
				IdUint64: 0,
				IdString: "zero",
				Data:     "zero",
			},
			{
				IdEnum:   relationships.PrimaryKeyEnum_ONE,
				IdInt32:  1,
				IdInt64:  1,
				IdUint32: 1,
				IdUint64: 1,
				IdString: "one",
				Data:     "one",
			},
			{
				IdEnum:   relationships.PrimaryKeyEnum_TWO,
				IdInt32:  2,
				IdInt64:  2,
				IdUint32: 2,
				IdUint64: 2,
				IdString: "two",
				Data:     "two",
			},
			{
				IdEnum:   relationships.PrimaryKeyEnum_THREE,
				IdInt32:  3,
				IdInt64:  3,
				IdUint32: 3,
				IdUint64: 3,
				IdString: "three",
				Data:     "three",
			},
		}
		initial := make([]*relationships.MAAll, 0, len(initialBuilder))
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

		expected := make([]*relationships.MAAll, 0, len(initial))
		for _, maall := range initial {
			expected = append(
				expected,
				relationships.MAAll_builder{
					IdEnum:   maall.GetIdEnum(),
					IdInt32:  maall.GetIdInt32(),
					IdInt64:  maall.GetIdInt64(),
					IdUint32: maall.GetIdUint32(),
					IdUint64: maall.GetIdUint64(),
					IdString: maall.GetIdString(),
					Data:     "UPDATED",
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

func TestMAAllRepository_Delete_WithLocatablePrimaryKeySucceeds(t *testing.T) {
	opts := maAllDefaultCmpOpts()

	testCases := map[string]struct {
		initial          []*relationships.MAAll_builder
		deleteExpression expressions.Expression
		expected         []*relationships.MAAll_builder
	}{
		"using primary key": {
			initial: []*relationships.MAAll_builder{
				{
					IdEnum:   relationships.PrimaryKeyEnum_ZERO,
					IdInt32:  0,
					IdInt64:  0,
					IdUint32: 0,
					IdUint64: 0,
					IdString: "zero",
					Data:     "zero",
				},
				{
					IdEnum:   relationships.PrimaryKeyEnum_ONE,
					IdInt32:  1,
					IdInt64:  1,
					IdUint32: 1,
					IdUint64: 1,
					IdString: "one",
					Data:     "one",
				},
				{
					IdEnum:   relationships.PrimaryKeyEnum_TWO,
					IdInt32:  2,
					IdInt64:  2,
					IdUint32: 2,
					IdUint64: 2,
					IdString: "two",
					Data:     "two",
				},
				{
					IdEnum:   relationships.PrimaryKeyEnum_THREE,
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
					expressions.NewIdentifier(relationships.Maall_IdEnum_Field),
					expressions.NewScalar(relationships.PrimaryKeyEnum_TWO),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(relationships.Maall_IdEnum_Field),
					expressions.NewScalar(relationships.PrimaryKeyEnum_THREE),
				),
			),
			expected: []*relationships.MAAll_builder{
				{
					IdEnum:   relationships.PrimaryKeyEnum_ZERO,
					IdInt32:  0,
					IdInt64:  0,
					IdUint32: 0,
					IdUint64: 0,
					IdString: "zero",
					Data:     "zero",
				},
				{
					IdEnum:   relationships.PrimaryKeyEnum_ONE,
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
			initial: []*relationships.MAAll_builder{
				{
					IdEnum:   relationships.PrimaryKeyEnum_ZERO,
					IdInt32:  0,
					IdInt64:  0,
					IdUint32: 0,
					IdUint64: 0,
					IdString: "zero",
					Data:     "zero",
				},
				{
					IdEnum:   relationships.PrimaryKeyEnum_ONE,
					IdInt32:  1,
					IdInt64:  1,
					IdUint32: 1,
					IdUint64: 1,
					IdString: "one",
					Data:     "one",
				},
				{
					IdEnum:   relationships.PrimaryKeyEnum_TWO,
					IdInt32:  2,
					IdInt64:  2,
					IdUint32: 2,
					IdUint64: 2,
					IdString: "two",
					Data:     "two",
				},
				{
					IdEnum:   relationships.PrimaryKeyEnum_THREE,
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
					expressions.NewIdentifier(relationships.Maall_Data_Field),
					expressions.NewScalar("two"),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(relationships.Maall_Data_Field),
					expressions.NewScalar("three"),
				),
			),
			expected: []*relationships.MAAll_builder{
				{
					IdEnum:   relationships.PrimaryKeyEnum_ZERO,
					IdInt32:  0,
					IdInt64:  0,
					IdUint32: 0,
					IdUint64: 0,
					IdString: "zero",
					Data:     "zero",
				},
				{
					IdEnum:   relationships.PrimaryKeyEnum_ONE,
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
		for repoType, componentUnderTest := range maAllImplementationsToTest() {
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
			initial := make([]*relationships.MAAll, 0, len(testCase.initial))
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
			expected := maAllBuild(testCase.expected)
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

func maAllBuild(in []*relationships.MAAll_builder) []*relationships.MAAll {
	out := make([]*relationships.MAAll, 0, len(in))
	for _, builder := range in {
		out = append(out, builder.Build())
	}
	return out
}

func maAllImplementationsToTest() map[options.Implementation]maAllComponentUnderTest {
	return map[options.Implementation]maAllComponentUnderTest{
		options.Implementation_IMPLEMENTATION_SQLITE: sqliteMAAllComponentUnderTest,
		options.Implementation_IMPLEMENTATION_PGSQL:  pgsqlMAAllComponentUnderTest,
	}
}

func maAllDefaultCmpOpts() cmp.Options {
	return cmp.Options{
		cmpopts.IgnoreUnexported(relationships.MAAll{}),
		cmpopts.SortSlices(func(x, y *relationships.MAAll) bool {
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
