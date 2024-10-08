package as_timestamp_field_test

import (
	"database/sql"
	"os"
	"testing"

	as_timestamp_field "github.com/samlitowitz/protoc-gen-crud/test-cases/as-timestamp-field"

	test_cases "github.com/samlitowitz/protoc-gen-crud/test-cases"
)

func pgsqlAsTimestampComponentUnderTest(t *testing.T) as_timestamp_field.AsTimestampRepository {
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

	repo, err := as_timestamp_field.NewPgSQLAsTimestampRepository(db)
	if err != nil {
		t.Fatal("pgsql: creating repository: ", err)
	}
	return repo
}
