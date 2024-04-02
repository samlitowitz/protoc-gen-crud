package sqlite

import (
	"github.com/iancoleman/strcase"
)

func SQLiteQuotedIdent(s string) string {
	return "\"" + SQLiteIdent(s) + "\""
}

func SQLiteIdent(s string) string {
	return strcase.ToSnake(s)
}

func SQLiteTemplateIdent(s string) string {
	return "\"" + SQLiteIdent(s) + "\""
}
