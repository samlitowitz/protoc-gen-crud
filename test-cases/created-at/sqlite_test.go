package created_at_test

import (
	"database/sql"
	"os"
	"testing"

	created_at "github.com/samlitowitz/protoc-gen-crud/test-cases/created-at"
)

func sqliteExecSQLFile(db *sql.DB, file string) error {
	code, err := os.ReadFile(file)
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

func sqliteCreatedAtComponentUnderTest(t *testing.T) created_at.CreatedAtRepository {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal("sqlite: ", err)
	}
	t.Cleanup(func() {
		err := db.Close()
		if err != nil {
			t.Fatal("sqlite: ", err)
		}
	})

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal("sqlite: finding working dir:", err)
	}

	err = sqliteExecSQLFile(db, origDir+string(os.PathSeparator)+"test.sqlite.sql")
	if err != nil {
		t.Fatal("sqlite: executing setup SQL: ", err)
	}

	repo, err := created_at.NewSQLiteCreatedAtRepository(db)
	if err != nil {
		t.Fatal("sqlite: creating repository: ", err)
	}
	return repo
}
