package relationships_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/samlitowitz/protoc-gen-crud/test-cases/relationships"

	test_cases "github.com/samlitowitz/protoc-gen-crud/test-cases"

	"github.com/samlitowitz/protoc-gen-crud/expressions"

	"github.com/google/go-cmp/cmp"
)

func pgsqlSAEnumComponentUnderTest(t *testing.T) relationships.SAEnumRepository {
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

	repo, err := relationships.NewPgSQLSAEnumRepository(db)
	if err != nil {
		t.Fatal("pgsql: creating repository: ", err)
	}
	return repo
}

func pgsqlSAInt32ComponentUnderTest(t *testing.T) relationships.SAInt32Repository {
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

	repo, err := relationships.NewPgSQLSAInt32Repository(db)
	if err != nil {
		t.Fatal("pgsql: creating repository: ", err)
	}
	return repo
}

func pgsqlSAInt64ComponentUnderTest(t *testing.T) relationships.SAInt64Repository {
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

	repo, err := relationships.NewPgSQLSAInt64Repository(db)
	if err != nil {
		t.Fatal("pgsql: creating repository: ", err)
	}
	return repo
}

func pgsqlSAUint32ComponentUnderTest(t *testing.T) relationships.SAUint32Repository {
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

	repo, err := relationships.NewPgSQLSAUint32Repository(db)
	if err != nil {
		t.Fatal("pgsql: creating repository: ", err)
	}
	return repo
}

func pgsqlSAUint64ComponentUnderTest(t *testing.T) relationships.SAUint64Repository {
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

	repo, err := relationships.NewPgSQLSAUint64Repository(db)
	if err != nil {
		t.Fatal("pgsql: creating repository: ", err)
	}
	return repo
}

func pgsqlSAStringComponentUnderTest(t *testing.T) relationships.SAStringRepository {
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

	repo, err := relationships.NewPgSQLSAStringRepository(db)
	if err != nil {
		t.Fatal("pgsql: creating repository: ", err)
	}
	return repo
}

func pgsqlMAAllComponentUnderTest(t *testing.T) relationships.MAAllRepository {
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

	repo, err := relationships.NewPgSQLMAAllRepository(db)
	if err != nil {
		t.Fatal("pgsql: creating repository: ", err)
	}
	return repo
}

func pgsqlSAInt32CreateSuccessWithReadAfterCheck(
	opts cmp.Options,
	repo relationships.SAInt32Repository,
	ctx context.Context,
	toCreate []*relationships.SAInt32,
	expectedRead []*relationships.SAInt32,
) error {
	_, err := repo.Create(ctx, toCreate)
	if err != nil {
		return fmt.Errorf("Create: %s", err)
	}
	return pgsqlSAInt32ReadCheck(opts, repo, ctx, nil, expectedRead)
}

func pgsqlSAInt32ReadCheck(
	opts cmp.Options,
	repo relationships.SAInt32Repository,
	ctx context.Context,
	expr expressions.Expression,
	expectedRead []*relationships.SAInt32,
) error {
	read, err := repo.Read(ctx, expr)
	if err != nil {
		return fmt.Errorf("Read: %s", err)
	}
	if diff := cmp.Diff(expectedRead, read, opts); diff != "" {
		return fmt.Errorf(mismatch("Read:", diff))
	}
	return nil
}

func pgsqlMAAllCreateSuccessWithReadAfterCheck(
	opts cmp.Options,
	repo relationships.MAAllRepository,
	ctx context.Context,
	toCreate []*relationships.MAAll,
	expectedRead []*relationships.MAAll,
) error {
	_, err := repo.Create(ctx, toCreate)
	if err != nil {
		return fmt.Errorf("Create: %s", err)
	}
	return pgsqlMAAllReadCheck(opts, repo, ctx, nil, expectedRead)
}

func pgsqlMAAllReadCheck(
	opts cmp.Options,
	repo relationships.MAAllRepository,
	ctx context.Context,
	expr expressions.Expression,
	expectedRead []*relationships.MAAll,
) error {
	read, err := repo.Read(ctx, expr)
	if err != nil {
		return fmt.Errorf("Read: %s", err)
	}
	if diff := cmp.Diff(expectedRead, read, opts); diff != "" {
		return fmt.Errorf(mismatch("Read:", diff))
	}
	return nil
}
