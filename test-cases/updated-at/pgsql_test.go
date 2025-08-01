package updated_at_test

import (
	"database/sql"
	"os"
	"testing"

	updated_at "github.com/samlitowitz/protoc-gen-crud/test-cases/updated-at"

	test_cases "github.com/samlitowitz/protoc-gen-crud/test-cases"
)

func pgsqlUpdatedAtComponentUnderTest(t *testing.T) updated_at.UpdatedAtRepository {
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

	repo, err := updated_at.NewPgSQLUpdatedAtRepository(db)
	if err != nil {
		t.Fatal("pgsql: creating repository: ", err)
	}
	return repo
}
