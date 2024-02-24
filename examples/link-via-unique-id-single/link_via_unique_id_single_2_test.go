package link_via_unique_id_single_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/mennanov/fmutils"
	"github.com/samlitowitz/protoc-gen-crud/expressions"
	"google.golang.org/genproto/protobuf/field_mask"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"

	"github.com/google/uuid"
	simple "github.com/samlitowitz/protoc-gen-crud/examples/link-via-unique-id-single"
)

func TestSQLiteProfileRepository_Create(t *testing.T) {
	// REFURL: https://github.com/golang/go/blob/988b718f4130ab5b3ce5a5774e1a58e83c92a163/src/path/filepath/path_test.go#L600
	// -- START -- //
	if runtime.GOOS == "ios" {
		restore := chtmpdir(t)
		defer restore()
	}

	tmpDir := t.TempDir()

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal("finding working dir:", err)
	}
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatal("entering temp dir:", err)
	}
	defer os.Chdir(origDir)
	// -- END -- //

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	opts := cmp.Options{
		cmpopts.IgnoreUnexported(simple.Profile{}),
		cmpopts.IgnoreFields(simple.Profile{}, "FieldMask"),
		cmpopts.SortSlices(func(x, y *simple.Profile) bool {
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

	tests := map[string]struct {
		expectedProfiles []*simple.Profile
		expr             expressions.Expression
	}{
		"no field mask sets all fields": {
			[]*simple.Profile{
				{
					Id:   uuid.NewString(),
					Name: "name-1",
				},
				{
					Id:   uuid.NewString(),
					Name: "name-2",
				},
				{
					Id:   uuid.NewString(),
					Name: "name-3",
				},
			},
			nil,
		},
		"field mask only affects specifying users": {
			[]*simple.Profile{
				{
					FieldMask: &field_mask.FieldMask{
						Paths: []string{"id"},
					},
					Id:   "1",
					Name: "name-1",
				},
				{
					FieldMask: &field_mask.FieldMask{
						Paths: []string{"id", "name"},
					},
					Id:   uuid.NewString(),
					Name: "name-2",
				},
				{
					FieldMask: &field_mask.FieldMask{
						Paths: []string{"id", "password"},
					},
					Id:   uuid.NewString(),
					Name: "name-3",
				},
				{
					FieldMask: &field_mask.FieldMask{
						Paths: []string{"id", "name", "password"},
					},
					Id:   uuid.NewString(),
					Name: "name-4",
				},
			},
			nil,
		},
	}

	for testCase, testData := range tests {
		err = createTable(db, origDir)
		if err != nil {
			t.Fatal(err)
		}

		repo, err := simple.NewSQLiteProfileRepository(db)
		if err != nil {
			t.Fatal(err)
		}

		users, err := repo.Read(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(users) != 0 {
			t.Fatalf(
				"%s: expected 0 users, got %d",
				testCase,
				len(users),
			)
		}

		toCreate := make([]*simple.Profile, 0, len(testData.expectedProfiles))
		for _, user := range testData.expectedProfiles {
			toCreate = append(toCreate, user)
		}

		// Check create response
		actualProfiles, err := repo.Create(context.Background(), toCreate)
		if err != nil {
			t.Fatalf("%s: %s", testCase, err)
		}

		if diff := cmp.Diff(toCreate, actualProfiles, opts); diff != "" {
			t.Fatalf(
				"%s: Create() mismatch (-want +got):\n%s",
				testCase,
				diff,
			)
		}

		for _, user := range testData.expectedProfiles {
			if user.FieldMask == nil {
				continue
			}
			fmutils.Filter(user, user.FieldMask.GetPaths())
		}

		// Check stored data
		actualProfiles, err = repo.Read(context.Background(), testData.expr)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(testData.expectedProfiles, actualProfiles, opts); diff != "" {
			t.Fatalf(
				"%s: Read() mismatch (-want +got):\n%s",
				testCase,
				diff,
			)
		}
	}
}

func TestSQLiteProfileRepository_Read(t *testing.T) {
	// REFURL: https://github.com/golang/go/blob/988b718f4130ab5b3ce5a5774e1a58e83c92a163/src/path/filepath/path_test.go#L600
	// -- START -- //
	if runtime.GOOS == "ios" {
		restore := chtmpdir(t)
		defer restore()
	}

	tmpDir := t.TempDir()

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal("finding working dir:", err)
	}
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatal("entering temp dir:", err)
	}
	defer os.Chdir(origDir)
	// -- END -- //

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	opts := cmp.Options{
		cmpopts.IgnoreUnexported(simple.Profile{}),
		cmpopts.SortSlices(func(x, y *simple.Profile) bool {
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

	tests := map[string]struct {
		unexpectedProfiles []*simple.Profile
		expectedProfiles   []*simple.Profile
		expr               expressions.Expression
	}{
		"no expression returns all users": {
			[]*simple.Profile{},
			[]*simple.Profile{
				{
					Id:   uuid.NewString(),
					Name: "name-1",
				},
				{
					Id:   uuid.NewString(),
					Name: "name-2",
				},
				{
					Id:   uuid.NewString(),
					Name: "name-3",
				},
			},
			nil,
		},
		"id equals expression returns matched user": {
			[]*simple.Profile{
				{
					Id:   uuid.NewString(),
					Name: "name-2",
				},
				{
					Id:   uuid.NewString(),
					Name: "name-3",
				},
			},
			[]*simple.Profile{
				{
					Id:   "1",
					Name: "name-1",
				},
			},
			expressions.NewEqual(
				expressions.NewIdentifier(simple.Profile_Id_Field),
				expressions.NewScalar("1"),
			),
		},
		"name equals expression returns matched user": {
			[]*simple.Profile{
				{
					Id:   uuid.NewString(),
					Name: "name-2",
				},
				{
					Id:   uuid.NewString(),
					Name: "name-3",
				},
			},
			[]*simple.Profile{
				{
					Id:   uuid.NewString(),
					Name: "name-1",
				},
			},
			expressions.NewEqual(
				expressions.NewIdentifier(simple.Profile_Name_Field),
				expressions.NewScalar("name-1"),
			),
		},
	}

	for testCase, testData := range tests {
		err = createTable(db, origDir)
		if err != nil {
			t.Fatal(err)
		}

		repo, err := simple.NewSQLiteProfileRepository(db)
		if err != nil {
			t.Fatal(err)
		}

		users, err := repo.Read(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(users) != 0 {
			t.Fatalf(
				"%s: expected 0 users, got %d",
				testCase,
				len(users),
			)
		}

		toCreate := make([]*simple.Profile, 0, len(testData.unexpectedProfiles)+len(testData.expectedProfiles))
		for _, user := range testData.unexpectedProfiles {
			toCreate = append(toCreate, user)
		}
		for _, user := range testData.expectedProfiles {
			toCreate = append(toCreate, user)
		}

		actualProfiles, err := repo.Create(context.Background(), toCreate)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(toCreate, actualProfiles, opts); diff != "" {
			t.Fatalf(
				"%s: Create() mismatch (-want +got):\n%s",
				testCase,
				diff,
			)
		}

		actualProfiles, err = repo.Read(context.Background(), testData.expr)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(testData.expectedProfiles, actualProfiles, opts); diff != "" {
			t.Fatalf(
				"%s: Read() mismatch (-want +got):\n%s",
				testCase,
				diff,
			)
		}
	}
}

func TestSQLiteProfileRepository_Update(t *testing.T) {
	// REFURL: https://github.com/golang/go/blob/988b718f4130ab5b3ce5a5774e1a58e83c92a163/src/path/filepath/path_test.go#L600
	// -- START -- //
	if runtime.GOOS == "ios" {
		restore := chtmpdir(t)
		defer restore()
	}

	tmpDir := t.TempDir()

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal("finding working dir:", err)
	}
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatal("entering temp dir:", err)
	}
	defer os.Chdir(origDir)
	// -- END -- //

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	opts := cmp.Options{
		cmpopts.IgnoreUnexported(simple.Profile{}),
		cmpopts.IgnoreFields(simple.Profile{}, "FieldMask"),
		cmpopts.SortSlices(func(x, y *simple.Profile) bool {
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

	tests := map[string]struct {
		createProfiles   []*simple.Profile
		updates          []*simple.Profile
		expectedProfiles []*simple.Profile
		expr             expressions.Expression
	}{
		"update name": {
			[]*simple.Profile{
				{
					Id:   uuid.NewString(),
					Name: "name-1",
				},
				{
					Id:   uuid.NewString(),
					Name: "name-2",
				},
				{
					Id:   uuid.NewString(),
					Name: "name-3",
				},
			},
			[]*simple.Profile{
				{
					Id:   "",
					Name: "name-1-1",
				},
				{
					Id:   "",
					Name: "name-2-1",
				},
				{
					Id:   "",
					Name: "name-3-1",
				},
			},
			[]*simple.Profile{
				{
					Id:   "",
					Name: "name-1-1",
				},
				{
					Id:   "",
					Name: "name-2-1",
				},
				{
					Id:   "",
					Name: "name-3-1",
				},
			},
			nil,
		},
		"update with field mask only affects specifying users": {
			[]*simple.Profile{
				{
					Id:   uuid.NewString(),
					Name: "name-1",
				},
				{
					Id:   uuid.NewString(),
					Name: "name-2",
				},
				{
					Id:   uuid.NewString(),
					Name: "name-3",
				},
				{
					Id:   uuid.NewString(),
					Name: "name-4",
				},
			},
			[]*simple.Profile{
				{
					Id:   "",
					Name: "name-1-1",
				},
				{
					FieldMask: &field_mask.FieldMask{
						Paths: []string{"name"},
					},
					Id:   "",
					Name: "name-2-1",
				},
				{
					FieldMask: &field_mask.FieldMask{
						Paths: []string{"password"},
					},
					Id:   "",
					Name: "name-3-1",
				},
				{
					FieldMask: &field_mask.FieldMask{
						Paths: []string{"id", "name", "password"},
					},
					Id:   "",
					Name: "name-4-1",
				},
			},
			[]*simple.Profile{
				{
					Id:   "",
					Name: "name-1-1",
				},
				{
					Id:   "",
					Name: "name-2-1",
				},
				{
					Id:   "",
					Name: "name-3",
				},
				{
					Id:   "",
					Name: "name-4-1",
				},
			},
			nil,
		},
	}

	for testCase, testData := range tests {
		err = createTable(db, origDir)
		if err != nil {
			t.Fatal(err)
		}

		repo, err := simple.NewSQLiteProfileRepository(db)
		if err != nil {
			t.Fatal(err)
		}

		users, err := repo.Read(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(users) != 0 {
			t.Fatalf(
				"%s: expected 0 users, got %d",
				testCase,
				len(users),
			)
		}

		actualProfiles, err := repo.Create(context.Background(), testData.createProfiles)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(testData.createProfiles, actualProfiles, opts); diff != "" {
			t.Fatalf(
				"%s: Create() mismatch (-want +got):\n%s",
				testCase,
				diff,
			)
		}

		actualProfiles, err = repo.Read(context.Background(), testData.expr)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(testData.createProfiles, actualProfiles, opts); diff != "" {
			t.Fatalf(
				"%s: Read() mismatch (-want +got):\n%s",
				testCase,
				diff,
			)
		}

		for i, user := range testData.createProfiles {
			if testData.expectedProfiles[i].Id == "" {
				testData.expectedProfiles[i].Id = user.Id
			}
			if testData.updates[i].Id == "" {
				testData.updates[i].Id = user.Id
			}
		}

		_, err = repo.Update(context.Background(), testData.updates)
		if err != nil {
			t.Fatal(err)
		}

		actualProfiles, err = repo.Read(context.Background(), testData.expr)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(testData.expectedProfiles, actualProfiles, opts); diff != "" {
			t.Fatalf(
				"%s: Read() mismatch (-want +got):\n%s",
				testCase,
				diff,
			)
		}
	}
}

func TestSQLiteProfileRepository_Delete(t *testing.T) {
	// REFURL: https://github.com/golang/go/blob/988b718f4130ab5b3ce5a5774e1a58e83c92a163/src/path/filepath/path_test.go#L600
	// -- START -- //
	if runtime.GOOS == "ios" {
		restore := chtmpdir(t)
		defer restore()
	}

	tmpDir := t.TempDir()

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal("finding working dir:", err)
	}
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatal("entering temp dir:", err)
	}
	defer os.Chdir(origDir)
	// -- END -- //

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = createTable(db, origDir)
	if err != nil {
		t.Fatal(err)
	}

	opts := cmp.Options{
		cmpopts.IgnoreUnexported(simple.Profile{}),
		cmpopts.SortSlices(func(x, y *simple.Profile) bool {
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

	repo, err := simple.NewSQLiteProfileRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	users, err := repo.Read(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 0 {
		t.Fatalf("expected 0 users, got %d", len(users))
	}

	expectedUserCount := 3
	expectedProfiles := make([]*simple.Profile, 0, expectedUserCount)
	for i := 0; i < expectedUserCount; i++ {
		expectedProfiles = append(expectedProfiles, &simple.Profile{
			Id:   uuid.NewString(),
			Name: fmt.Sprintf("name-%d", i),
		})
	}

	actualProfiles, err := repo.Create(context.Background(), expectedProfiles)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectedProfiles, actualProfiles, opts); diff != "" {
		t.Fatalf("Create() mismatch (-want +got):\n%s", diff)
	}

	expr := expressions.NewEqual(
		expressions.NewIdentifier(simple.Profile_Id_Field),
		expressions.NewScalar(expectedProfiles[0].Id),
	)
	err = repo.Delete(context.Background(), expr)
	if err != nil {
		t.Fatal(err)
	}

	expectedProfiles = expectedProfiles[1:]
	actualProfiles, err = repo.Read(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectedProfiles, actualProfiles, opts); diff != "" {
		t.Fatalf("Create() mismatch (-want +got):\n%s", diff)
	}
}

func createTable(db *sql.DB, dir string) error {
	code, err := os.ReadFile(path.Join(dir, "user.sqlite.sql"))
	if err != nil {
		return err
	}
	stmt, err := db.Prepare(string(code))
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

// REFURL: https://github.com/golang/go/blob/988b718f4130ab5b3ce5a5774e1a58e83c92a163/src/path/filepath/path_test.go#L553
func chtmpdir(t *testing.T) (restore func()) {
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	d, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	if err := os.Chdir(d); err != nil {
		t.Fatalf("chtmpdir: %v", err)
	}
	return func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatalf("chtmpdir: %v", err)
		}
		os.RemoveAll(d)
	}
}
