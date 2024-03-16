package many_to_many_relationship_test

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
	simple "github.com/samlitowitz/protoc-gen-crud/examples/many-to-many-relationship"
)

func TestSQLiteUserRepository_Create(t *testing.T) {
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
		cmpopts.IgnoreUnexported(simple.User{}),
		cmpopts.IgnoreFields(simple.User{}, "FieldMask"),
		cmpopts.IgnoreUnexported(simple.Role{}),
		cmpopts.IgnoreFields(simple.Role{}, "FieldMask"),
		cmpopts.SortSlices(func(x, y *simple.User) bool {
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
		expectedUsers []*simple.User
		expr          expressions.Expression
	}{
		"no field mask sets all fields": {
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-1",
					Password: "password-1",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					Id:       uuid.NewString(),
					Username: "username-2",
					Password: "password-2",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
			},
			nil,
		},
		"field mask only affects specifying users": {
			[]*simple.User{
				{
					FieldMask: &field_mask.FieldMask{
						Paths: []string{"id"},
					},
					Id:       "1",
					Username: "username-1",
					Password: "password-1",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					FieldMask: &field_mask.FieldMask{
						Paths: []string{"id", "username"},
					},
					Id:       uuid.NewString(),
					Username: "username-2",
					Password: "password-2",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					FieldMask: &field_mask.FieldMask{
						Paths: []string{"id", "password"},
					},
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					FieldMask: &field_mask.FieldMask{
						Paths: []string{"id", "username", "password"},
					},
					Id:       uuid.NewString(),
					Username: "username-4",
					Password: "password-4",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
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

		repo, err := simple.NewSQLiteUserRepository(db)
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

		toCreate := make([]*simple.User, 0, len(testData.expectedUsers))
		for _, user := range testData.expectedUsers {
			toCreate = append(toCreate, user)
		}

		// Check create response
		actualUsers, err := repo.Create(context.Background(), toCreate)
		if err != nil {
			t.Fatalf("%s: %s", testCase, err)
		}

		if diff := cmp.Diff(toCreate, actualUsers, opts); diff != "" {
			t.Fatalf(
				"%s: Create() mismatch (-want +got):\n%s",
				testCase,
				diff,
			)
		}

		for _, user := range testData.expectedUsers {
			if user.FieldMask == nil {
				continue
			}
			fmutils.Filter(user, user.FieldMask.GetPaths())
			user.Roles = []*simple.Role{}
		}

		// Check stored data
		actualUsers, err = repo.Read(context.Background(), testData.expr)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(testData.expectedUsers, actualUsers, opts); diff != "" {
			t.Fatalf(
				"%s: Read() mismatch (-want +got):\n%s",
				testCase,
				diff,
			)
		}
	}
}

func TestSQLiteUserRepository_Read(t *testing.T) {
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
		cmpopts.IgnoreUnexported(simple.User{}),
		cmpopts.IgnoreUnexported(simple.Role{}),
		cmpopts.SortSlices(func(x, y *simple.User) bool {
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
		unexpectedUsers []*simple.User
		expectedUsers   []*simple.User
		expr            expressions.Expression
	}{
		"no expression returns all users": {
			[]*simple.User{},
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-1",
					Password: "password-1",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					Id:       uuid.NewString(),
					Username: "username-2",
					Password: "password-2",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
			},
			nil,
		},
		"id equals expression returns matched user": {
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-2",
					Password: "password-2",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
			},
			[]*simple.User{
				{
					Id:       "1",
					Username: "username-1",
					Password: "password-1",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
			},
			expressions.NewEqual(
				expressions.NewIdentifier(simple.User_Id_Field),
				expressions.NewScalar("1"),
			),
		},
		"username equals expression returns matched user": {
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-2",
					Password: "password-2",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
			},
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-1",
					Password: "password-1",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
			},
			expressions.NewEqual(
				expressions.NewIdentifier(simple.User_Username_Field),
				expressions.NewScalar("username-1"),
			),
		},
		"password equals expression returns matched user": {
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-2",
					Password: "password-2",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
			},
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-1",
					Password: "password-1",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
			},
			expressions.NewEqual(
				expressions.NewIdentifier(simple.User_Password_Field),
				expressions.NewScalar("password-1"),
			),
		},
		"username and password equals expression returns matched user": {
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-2",
					Password: "password-2",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
			},
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-1",
					Password: "password-1",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
			},
			expressions.NewAnd(
				expressions.NewEqual(
					expressions.NewIdentifier(simple.User_Username_Field),
					expressions.NewScalar("username-1"),
				),
				expressions.NewEqual(
					expressions.NewIdentifier(simple.User_Password_Field),
					expressions.NewScalar("password-1"),
				),
			),
		},
	}

	for testCase, testData := range tests {
		err = createTable(db, origDir)
		if err != nil {
			t.Fatal(err)
		}

		repo, err := simple.NewSQLiteUserRepository(db)
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

		toCreate := make([]*simple.User, 0, len(testData.unexpectedUsers)+len(testData.expectedUsers))
		for _, user := range testData.unexpectedUsers {
			toCreate = append(toCreate, user)
		}
		for _, user := range testData.expectedUsers {
			toCreate = append(toCreate, user)
		}

		actualUsers, err := repo.Create(context.Background(), toCreate)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(toCreate, actualUsers, opts); diff != "" {
			t.Fatalf(
				"%s: Create() mismatch (-want +got):\n%s",
				testCase,
				diff,
			)
		}

		actualUsers, err = repo.Read(context.Background(), testData.expr)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(testData.expectedUsers, actualUsers, opts); diff != "" {
			t.Fatalf(
				"%s: Read() mismatch (-want +got):\n%s",
				testCase,
				diff,
			)
		}
	}
}

func TestSQLiteUserRepository_Update(t *testing.T) {
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
		cmpopts.IgnoreUnexported(simple.User{}),
		cmpopts.IgnoreFields(simple.User{}, "FieldMask"),
		cmpopts.IgnoreUnexported(simple.Role{}),
		cmpopts.IgnoreFields(simple.Role{}, "FieldMask"),
		cmpopts.SortSlices(func(x, y *simple.User) bool {
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
		createUsers   []*simple.User
		updates       []*simple.User
		expectedUsers []*simple.User
		expr          expressions.Expression
	}{
		"update username and password": {
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-1",
					Password: "password-1",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					Id:       uuid.NewString(),
					Username: "username-2",
					Password: "password-2",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
			},
			[]*simple.User{
				{
					Id:       "",
					Username: "username-1-1",
					Password: "password-1-1",
				},
				{
					Id:       "",
					Username: "username-2-1",
					Password: "password-2-1",
				},
				{
					Id:       "",
					Username: "username-3-1",
					Password: "password-3-1",
				},
			},
			[]*simple.User{
				{
					Id:       "",
					Username: "username-1-1",
					Password: "password-1-1",
				},
				{
					Id:       "",
					Username: "username-2-1",
					Password: "password-2-1",
				},
				{
					Id:       "",
					Username: "username-3-1",
					Password: "password-3-1",
				},
			},
			nil,
		},
		"update with field mask only affects specifying users": {
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-1",
					Password: "password-1",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					Id:       uuid.NewString(),
					Username: "username-2",
					Password: "password-2",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
				{
					Id:       uuid.NewString(),
					Username: "username-4",
					Password: "password-4",
					Roles: []*simple.Role{
						{
							Id: uuid.NewString(),
						},
					},
				},
			},
			[]*simple.User{
				{
					Id:       "",
					Username: "username-1-1",
					Password: "password-1-1",
				},
				{
					FieldMask: &field_mask.FieldMask{
						Paths: []string{"username"},
					},
					Id:       "",
					Username: "username-2-1",
					Password: "password-2-1",
				},
				{
					FieldMask: &field_mask.FieldMask{
						Paths: []string{"password"},
					},
					Id:       "",
					Username: "username-3-1",
					Password: "password-3-1",
				},
				{
					FieldMask: &field_mask.FieldMask{
						Paths: []string{"id", "username", "password"},
					},
					Id:       "",
					Username: "username-4-1",
					Password: "password-4-1",
				},
			},
			[]*simple.User{
				{
					Id:       "",
					Username: "username-1-1",
					Password: "password-1-1",
				},
				{
					Id:       "",
					Username: "username-2-1",
					Password: "password-2",
				},
				{
					Id:       "",
					Username: "username-3",
					Password: "password-3-1",
				},
				{
					Id:       "",
					Username: "username-4-1",
					Password: "password-4-1",
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

		repo, err := simple.NewSQLiteUserRepository(db)
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

		actualUsers, err := repo.Create(context.Background(), testData.createUsers)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(testData.createUsers, actualUsers, opts); diff != "" {
			t.Fatalf(
				"%s: Create() mismatch (-want +got):\n%s",
				testCase,
				diff,
			)
		}

		actualUsers, err = repo.Read(context.Background(), testData.expr)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(testData.createUsers, actualUsers, opts); diff != "" {
			t.Fatalf(
				"%s: Read() mismatch (-want +got):\n%s",
				testCase,
				diff,
			)
		}

		for i, user := range testData.createUsers {
			if testData.expectedUsers[i].Id == "" {
				testData.expectedUsers[i].Id = user.Id
			}
			if testData.expectedUsers[i].Roles == nil {
				testData.expectedUsers[i].Roles = []*simple.Role{}
			}
			if testData.expectedUsers[i].Profile.Id == "" {
				testData.expectedUsers[i].Profile.Id = user.Profile.GetId()
			}
			if testData.updates[i].Id == "" {
				testData.updates[i].Id = user.Id
			}
			if testData.updates[i].Roles == nil {
				testData.updates[i].Roles = []*simple.Role{}
			}
			if testData.updates[i].Profile.Id == "" {
				testData.updates[i].Profile.Id = user.Profile.Id
			}
		}

		_, err = repo.Update(context.Background(), testData.updates)
		if err != nil {
			t.Fatal(err)
		}

		actualUsers, err = repo.Read(context.Background(), testData.expr)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(testData.expectedUsers, actualUsers, opts); diff != "" {
			t.Fatalf(
				"%s: Read() mismatch (-want +got):\n%s",
				testCase,
				diff,
			)
		}
	}
}

func TestSQLiteUserRepository_Delete(t *testing.T) {
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
		cmpopts.IgnoreUnexported(simple.User{}),
		cmpopts.IgnoreUnexported(simple.Role{}),
		cmpopts.SortSlices(func(x, y *simple.User) bool {
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

	repo, err := simple.NewSQLiteUserRepository(db)
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
	expectedUsers := make([]*simple.User, 0, expectedUserCount)
	for i := 0; i < expectedUserCount; i++ {
		expectedUsers = append(expectedUsers, &simple.User{
			Id:       uuid.NewString(),
			Username: fmt.Sprintf("username-%d", i),
			Password: fmt.Sprintf("password-%d", i),
			Profile: &simple.Role{
				Id: uuid.NewString(),
			},
		})
	}

	actualUsers, err := repo.Create(context.Background(), expectedUsers)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectedUsers, actualUsers, opts); diff != "" {
		t.Fatalf("Create() mismatch (-want +got):\n%s", diff)
	}

	expr := expressions.NewEqual(
		expressions.NewIdentifier(simple.User_Id_Field),
		expressions.NewScalar(expectedUsers[0].Id),
	)
	err = repo.Delete(context.Background(), expr)
	if err != nil {
		t.Fatal(err)
	}

	expectedUsers = expectedUsers[1:]
	actualUsers, err = repo.Read(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectedUsers, actualUsers, opts); diff != "" {
		t.Fatalf("Create() mismatch (-want +got):\n%s", diff)
	}
}

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
		cmpopts.IgnoreUnexported(simple.Role{}),
		cmpopts.IgnoreFields(simple.Role{}, "FieldMask"),
		cmpopts.SortSlices(func(x, y *simple.Role) bool {
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
		expectedProfiles []*simple.Role
		expr             expressions.Expression
	}{
		"no field mask sets all fields": {
			[]*simple.Role{
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
			[]*simple.Role{
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

		repo, err := simple.NewSQLiteRoleRepository(db)
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

		toCreate := make([]*simple.Role, 0, len(testData.expectedProfiles))
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
		cmpopts.IgnoreUnexported(simple.Role{}),
		cmpopts.SortSlices(func(x, y *simple.Role) bool {
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
		unexpectedProfiles []*simple.Role
		expectedProfiles   []*simple.Role
		expr               expressions.Expression
	}{
		"no expression returns all users": {
			[]*simple.Role{},
			[]*simple.Role{
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
			[]*simple.Role{
				{
					Id:   uuid.NewString(),
					Name: "name-2",
				},
				{
					Id:   uuid.NewString(),
					Name: "name-3",
				},
			},
			[]*simple.Role{
				{
					Id:   "1",
					Name: "name-1",
				},
			},
			expressions.NewEqual(
				expressions.NewIdentifier(simple.Role_Id_Field),
				expressions.NewScalar("1"),
			),
		},
		"name equals expression returns matched user": {
			[]*simple.Role{
				{
					Id:   uuid.NewString(),
					Name: "name-2",
				},
				{
					Id:   uuid.NewString(),
					Name: "name-3",
				},
			},
			[]*simple.Role{
				{
					Id:   uuid.NewString(),
					Name: "name-1",
				},
			},
			expressions.NewEqual(
				expressions.NewIdentifier(simple.Role_Name_Field),
				expressions.NewScalar("name-1"),
			),
		},
	}

	for testCase, testData := range tests {
		err = createTable(db, origDir)
		if err != nil {
			t.Fatal(err)
		}

		repo, err := simple.NewSQLiteRoleRepository(db)
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

		toCreate := make([]*simple.Role, 0, len(testData.unexpectedProfiles)+len(testData.expectedProfiles))
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
		cmpopts.IgnoreUnexported(simple.Role{}),
		cmpopts.IgnoreFields(simple.Role{}, "FieldMask"),
		cmpopts.SortSlices(func(x, y *simple.Role) bool {
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
		createProfiles   []*simple.Role
		updates          []*simple.Role
		expectedProfiles []*simple.Role
		expr             expressions.Expression
	}{
		"update name": {
			[]*simple.Role{
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
			[]*simple.Role{
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
			[]*simple.Role{
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
			[]*simple.Role{
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
			[]*simple.Role{
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
			[]*simple.Role{
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

		repo, err := simple.NewSQLiteRoleRepository(db)
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
		cmpopts.IgnoreUnexported(simple.Role{}),
		cmpopts.SortSlices(func(x, y *simple.Role) bool {
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

	repo, err := simple.NewSQLiteRoleRepository(db)
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
	expectedProfiles := make([]*simple.Role, 0, expectedUserCount)
	for i := 0; i < expectedUserCount; i++ {
		expectedProfiles = append(expectedProfiles, &simple.Role{
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
		expressions.NewIdentifier(simple.Role_Id_Field),
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
