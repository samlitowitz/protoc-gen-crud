package simple_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/samlitowitz/protoc-gen-crud/expressions"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"

	"github.com/google/uuid"
	"github.com/samlitowitz/protoc-gen-crud/examples/simple"
)

func TestSQLiteUserRepository_Create(t *testing.T) {
	opts := cmp.Options{
		cmpopts.IgnoreUnexported(simple.User{}),
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

	repo := simple.NewSQLiteUserRepository()

	users, err := repo.Read(nil)
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
		})
	}

	actualUsers, err := repo.Create(expectedUsers)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectedUsers, actualUsers, opts); diff != "" {
		t.Fatalf("Create() mismatch (-want +got):\n%s", diff)
	}
}

func TestSQLiteUserRepository_Read(t *testing.T) {
	opts := cmp.Options{
		cmpopts.IgnoreUnexported(simple.User{}),
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
				},
				{
					Id:       uuid.NewString(),
					Username: "username-2",
					Password: "password-2",
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
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
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
				},
			},
			[]*simple.User{
				{
					Id:       "1",
					Username: "username-1",
					Password: "password-1",
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
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
				},
			},
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-1",
					Password: "password-1",
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
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
				},
			},
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-1",
					Password: "password-1",
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
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
				},
			},
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-1",
					Password: "password-1",
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
		repo := simple.NewSQLiteUserRepository()

		users, err := repo.Read(nil)
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

		expectedUserCount := 3
		expectedUsers := make([]*simple.User, 0, expectedUserCount)
		for i := 0; i < expectedUserCount; i++ {
			expectedUsers = append(expectedUsers, &simple.User{
				Id:       uuid.NewString(),
				Username: fmt.Sprintf("username-%d", i),
				Password: fmt.Sprintf("password-%d", i),
			})
		}

		toCreate := make([]*simple.User, 0, len(testData.unexpectedUsers)+len(testData.expectedUsers))
		for _, user := range testData.unexpectedUsers {
			toCreate = append(toCreate, user)
		}
		for _, user := range testData.expectedUsers {
			toCreate = append(toCreate, user)
		}

		actualUsers, err := repo.Create(toCreate)
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

		actualUsers, err = repo.Read(testData.expr)
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
	opts := cmp.Options{
		cmpopts.IgnoreUnexported(simple.User{}),
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
		expectedUsers []*simple.User
		expr          expressions.Expression
	}{
		"update username": {
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-1",
					Password: "password-1",
				},
				{
					Id:       uuid.NewString(),
					Username: "username-2",
					Password: "password-2",
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
				},
			},
			[]*simple.User{
				{
					Id:       "",
					Username: "username-1-1",
					Password: "",
				},
				{
					Id:       "",
					Username: "username-2-1",
					Password: "",
				},
				{
					Id:       "",
					Username: "username-3-1",
					Password: "",
				},
			},
			nil,
		},
		"update password": {
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-1",
					Password: "password-1",
				},
				{
					Id:       uuid.NewString(),
					Username: "username-2",
					Password: "password-2",
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
				},
			},
			[]*simple.User{
				{
					Id:       "",
					Username: "",
					Password: "password-1-1",
				},
				{
					Id:       "",
					Username: "",
					Password: "password-2-1",
				},
				{
					Id:       "",
					Username: "",
					Password: "password-3-1",
				},
			},
			nil,
		},
		"update username and password": {
			[]*simple.User{
				{
					Id:       uuid.NewString(),
					Username: "username-1",
					Password: "password-1",
				},
				{
					Id:       uuid.NewString(),
					Username: "username-2",
					Password: "password-2",
				},
				{
					Id:       uuid.NewString(),
					Username: "username-3",
					Password: "password-3",
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
	}

	for testCase, testData := range tests {
		repo := simple.NewSQLiteUserRepository()

		users, err := repo.Read(nil)
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

		actualUsers, err := repo.Create(testData.createUsers)
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

		actualUsers, err = repo.Read(testData.expr)
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
			if testData.expectedUsers[i].Username == "" {
				testData.expectedUsers[i].Username = user.Username
			}
			if testData.expectedUsers[i].Password == "" {
				testData.expectedUsers[i].Password = user.Password
			}
		}

		actualUsers, err = repo.Update(testData.expectedUsers)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(testData.expectedUsers, actualUsers, opts); diff != "" {
			t.Fatalf(
				"%s: Create() mismatch (-want +got):\n%s",
				testCase,
				diff,
			)
		}

		actualUsers, err = repo.Read(testData.expr)
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
	opts := cmp.Options{
		cmpopts.IgnoreUnexported(simple.User{}),
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

	repo := simple.NewSQLiteUserRepository()

	users, err := repo.Read(nil)
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
		})
	}

	actualUsers, err := repo.Create(expectedUsers)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectedUsers, actualUsers, opts); diff != "" {
		t.Fatalf("Create() mismatch (-want +got):\n%s", diff)
	}

	err = repo.Delete([]*simple.User{expectedUsers[0]})
	if err != nil {
		t.Fatal(err)
	}

	expectedUsers = expectedUsers[1:]
	actualUsers, err = repo.Read(nil)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectedUsers, actualUsers, opts); diff != "" {
		t.Fatalf("Create() mismatch (-want +got):\n%s", diff)
	}
}
