package relationships_many_to_many_test

import (
	"database/sql"
	"os"
	"testing"

	relationships_many_to_many "github.com/samlitowitz/protoc-gen-crud/test-cases/relationships-many-to-many"
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

func sqliteSAInt32ComponentUnderTest(t *testing.T) relationships_many_to_many.SAInt32Repository {
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

	repo, err := relationships_many_to_many.NewSQLiteSAInt32Repository(db)
	if err != nil {
		t.Fatal("sqlite: creating repository: ", err)
	}
	return repo
}

func sqliteMAAllComponentUnderTest(t *testing.T) relationships_many_to_many.MAAllRepository {
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

	repo, err := relationships_many_to_many.NewSQLiteMAAllRepository(db)
	if err != nil {
		t.Fatal("sqlite: creating repository: ", err)
	}
	return repo
}
