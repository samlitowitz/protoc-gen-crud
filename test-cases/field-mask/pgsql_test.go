package field_mask_test

import (
	"database/sql"
	"os"
	"testing"

	test_cases "github.com/samlitowitz/protoc-gen-crud/test-cases"

	fieldMask "github.com/samlitowitz/protoc-gen-crud/test-cases/field-mask"
)

func pgsqlSAInt32ComponentUnderTest(t *testing.T) fieldMask.SAInt32Repository {
	dburl, err := test_cases.PgSQLDBURLFromEnv()
	if err != nil {
		t.Fatal("pgsql: dburl: ", err)
	}
	db, err := sql.Open("pgx", dburl)
	if err != nil {
		t.Fatal("pgsql: ", err)
	}
	t.Cleanup(func() {
		err := db.Close()
		if err != nil {
			t.Fatal("pgsql: ", err)
		}
	})

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal("pgsql: finding working dir:", err)
	}

	err = test_cases.PgSQLExecSQLFile(db, origDir+string(os.PathSeparator)+"test.pgsql.sql")
	if err != nil {
		t.Fatal("pgsql: executing setup SQL: ", err)
	}

	repo, err := fieldMask.NewPgSQLSAInt32Repository(db)
	if err != nil {
		t.Fatal("pgsql: creating repository: ", err)
	}
	return repo
}

func pgsqlMAAllComponentUnderTest(t *testing.T) fieldMask.MAAllRepository {
	dburl, err := test_cases.PgSQLDBURLFromEnv()
	if err != nil {
		t.Fatal("pgsql: dburl: ", err)
	}
	db, err := sql.Open("pgx", dburl)
	if err != nil {
		t.Fatal("pgsql: ", err)
	}
	t.Cleanup(func() {
		err := db.Close()
		if err != nil {
			t.Fatal("pgsql: ", err)
		}
	})

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal("pgsql: finding working dir:", err)
	}

	err = test_cases.PgSQLExecSQLFile(db, origDir+string(os.PathSeparator)+"test.pgsql.sql")
	if err != nil {
		t.Fatal("pgsql: executing setup SQL: ", err)
	}

	repo, err := fieldMask.NewPgSQLMAAllRepository(db)
	if err != nil {
		t.Fatal("pgsql: creating repository: ", err)
	}
	return repo
}
