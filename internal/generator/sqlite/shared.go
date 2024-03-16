package sqlite

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/samlitowitz/protoc-gen-crud/internal/casing"
	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
)

func SQLiteMemberField(f *descriptor.Field) string {
	if !f.HasRelationship() {
		return casing.CamelIdentifier(*f.Name)
	}
	minUIDFields := f.Relationship.CRUD.MinimalUIDFields()
	if len(minUIDFields) != 1 {
		panic(fmt.Errorf("message type must have unique identifier with exactly one field defined on field %s", f.GetName()))
	}
	return fmt.Sprintf(
		"%s.%s",
		casing.CamelIdentifier(*f.Name),
		casing.CamelIdentifier(*minUIDFields[0].Name),
	)
}

func SQLiteMemberAccessor(f *descriptor.Field) string {
	if !f.HasRelationship() {
		return fmt.Sprintf(
			"Get%s()",
			casing.CamelIdentifier(*f.Name),
		)
	}
	minUIDFields := f.Relationship.CRUD.MinimalUIDFields()
	if len(minUIDFields) != 1 {
		panic(fmt.Errorf("message type must have unique identifier with exactly one field defined on field %s", f.GetName()))
	}
	return fmt.Sprintf(
		"Get%s().Get%s()",
		casing.CamelIdentifier(*f.Name),
		casing.CamelIdentifier(*minUIDFields[0].Name),
	)
}

func SQLiteColumnNameFromFieldName(f *descriptor.Field) string {
	if !f.HasRelationship() {
		return SQLiteColumnName(*f.Name)
	}
	return SQLiteColumnName(*f.Name + "_id")
}

func SQLiteRelatesToManyTableName(f *descriptor.Field) string {
	return SQLiteTableName(
		fmt.Sprintf(
			"%s_%s",
			f.Message.CRUD.GetName(),
			f.GetName(),
		),
	)
}

func SQLiteRelatesToManyColumnName(f *descriptor.Field) string {
	return SQLiteTableName(
		fmt.Sprintf(
			"%s_%s",
			SQLiteColumnName(f.Message.CRUD.GetName()),
			SQLiteColumnNameFromFieldName(f),
		),
	)
}

func SQLiteIdent(s string) string {
	return "\"" + s + "\""
}

func SQLiteTemplateIdent(s string) string {
	return "\\\"" + s + "\\\""
}

func SQLiteTableName(s string) string {
	return strcase.ToSnake(s)
}

func SQLiteColumnName(s string) string {
	return strcase.ToSnake(s)
}
