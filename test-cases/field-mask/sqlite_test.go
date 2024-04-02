package field_mask_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/samlitowitz/protoc-gen-crud/expressions"

	"github.com/google/go-cmp/cmp"

	fieldMask "github.com/samlitowitz/protoc-gen-crud/test-cases/field-mask"
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

func sqliteSAEnumComponentUnderTest(t *testing.T) fieldMask.SAEnumRepository {
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

	repo, err := fieldMask.NewSQLiteSAEnumRepository(db)
	if err != nil {
		t.Fatal("sqlite: creating repository: ", err)
	}
	return repo
}

func sqliteSAInt32ComponentUnderTest(t *testing.T) fieldMask.SAInt32Repository {
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

	repo, err := fieldMask.NewSQLiteSAInt32Repository(db)
	if err != nil {
		t.Fatal("sqlite: creating repository: ", err)
	}
	return repo
}

func sqliteSAInt64ComponentUnderTest(t *testing.T) fieldMask.SAInt64Repository {
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

	repo, err := fieldMask.NewSQLiteSAInt64Repository(db)
	if err != nil {
		t.Fatal("sqlite: creating repository: ", err)
	}
	return repo
}

func sqliteSAUint32ComponentUnderTest(t *testing.T) fieldMask.SAUint32Repository {
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

	repo, err := fieldMask.NewSQLiteSAUint32Repository(db)
	if err != nil {
		t.Fatal("sqlite: creating repository: ", err)
	}
	return repo
}

func sqliteSAUint64ComponentUnderTest(t *testing.T) fieldMask.SAUint64Repository {
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

	repo, err := fieldMask.NewSQLiteSAUint64Repository(db)
	if err != nil {
		t.Fatal("sqlite: creating repository: ", err)
	}
	return repo
}

func sqliteSAStringComponentUnderTest(t *testing.T) fieldMask.SAStringRepository {
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

	repo, err := fieldMask.NewSQLiteSAStringRepository(db)
	if err != nil {
		t.Fatal("sqlite: creating repository: ", err)
	}
	return repo
}

func sqliteMAAllComponentUnderTest(t *testing.T) fieldMask.MAAllRepository {
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

	repo, err := fieldMask.NewSQLiteMAAllRepository(db)
	if err != nil {
		t.Fatal("sqlite: creating repository: ", err)
	}
	return repo
}

func sqliteSAInt32CreateSuccessWithReadAfterCheck(
	opts cmp.Options,
	repo fieldMask.SAInt32Repository,
	ctx context.Context,
	toCreate []*fieldMask.SAInt32,
	expectedRead []*fieldMask.SAInt32,
) error {
	_, err := repo.Create(ctx, toCreate)
	if err != nil {
		return fmt.Errorf("Create: %s", err)
	}
	return sqliteSAInt32ReadCheck(opts, repo, ctx, nil, expectedRead)
}

func sqliteSAInt32ReadCheck(
	opts cmp.Options,
	repo fieldMask.SAInt32Repository,
	ctx context.Context,
	expr expressions.Expression,
	expectedRead []*fieldMask.SAInt32,
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

func sqliteMAAllCreateSuccessWithReadAfterCheck(
	opts cmp.Options,
	repo fieldMask.MAAllRepository,
	ctx context.Context,
	toCreate []*fieldMask.MAAll,
	expectedRead []*fieldMask.MAAll,
) error {
	_, err := repo.Create(ctx, toCreate)
	if err != nil {
		return fmt.Errorf("Create: %s", err)
	}
	return sqliteMAAllReadCheck(opts, repo, ctx, nil, expectedRead)
}

func sqliteMAAllReadCheck(
	opts cmp.Options,
	repo fieldMask.MAAllRepository,
	ctx context.Context,
	expr expressions.Expression,
	expectedRead []*fieldMask.MAAll,
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
