package created_at_test

import (
	"database/sql"
	"os"
	"testing"

	created_at "github.com/samlitowitz/protoc-gen-crud/test-cases/created-at"

	test_cases "github.com/samlitowitz/protoc-gen-crud/test-cases"
)

func pgsqlCreatedAtComponentUnderTest(t *testing.T) created_at.CreatedAtRepository {
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

	repo, err := created_at.NewPgSQLCreatedAtRepository(db)
	if err != nil {
		t.Fatal("pgsql: creating repository: ", err)
	}
	return repo
}
