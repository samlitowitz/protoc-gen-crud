package test_cases

import (
	"testing"

	"modernc.org/sqlite"
	sqliteLib "modernc.org/sqlite/lib"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/samlitowitz/protoc-gen-crud/options"
)

func AssertSQLErrorCode(t *testing.T, typ options.Implementation, lut map[options.Implementation]any, err error, prefix string) {
	if _, ok := lut[typ]; !ok {
		t.Fatal(prefix, "missing LUT entry: ", typ.String())
	}
	switch typ {
	case options.Implementation_PGSQL:
		sqlErr, ok := err.(*pgconn.PgError)
		if !ok {
			t.Fatalf("%sexpected *pgconn.PgError, got %T", prefix, err)
		}
		expectedCode, ok := lut[typ].(string)
		if !ok {
			t.Fatal(prefix, "expected LUT value to be of type string")
		}
		if sqlErr.Code != expectedCode {
			t.Fatalf(
				"%sexpected duplicate error code, got %s: %s",
				prefix,
				sqlErr.Code,
				sqlErr.Detail,
			)
		}
	case options.Implementation_SQLITE:
		sqlErr, ok := err.(*sqlite.Error)
		if !ok {
			t.Fatalf("%sexpected *sqlite.Error, got %T", prefix, err)
		}
		if sqlErr.Code() != sqliteLib.SQLITE_CONSTRAINT_PRIMARYKEY {
			t.Fatalf(prefix, "expected duplicate error code, got %d", sqlErr.Code())
		}
	default:
		t.Fatal(prefix, "unhandled implementation: ", typ.String())
	}
}
