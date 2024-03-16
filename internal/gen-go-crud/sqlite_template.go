package gen_go_crud

import (
	"fmt"
	"text/template"

	"github.com/samlitowitz/protoc-gen-crud/internal/casing"

	"github.com/iancoleman/strcase"
	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
)

func init() {
	strcase.ConfigureAcronym("UID", "uid")
}

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
		return *f.Name
	}
	return *f.Name + "_id"
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

type sqlite struct{}

var (
	sqliteFuncMap template.FuncMap = map[string]interface{}{
		"sqliteIdent":                   SQLiteTemplateIdent,
		"sqliteTableName":               SQLiteTableName,
		"sqliteColumnName":              SQLiteColumnName,
		"sqliteColumnNameFromFieldName": SQLiteColumnNameFromFieldName,
		"sqliteMemberAccessor":          SQLiteMemberAccessor,
		"sqliteMemberField":             SQLiteMemberField,
		"toLowerCamel":                  strcase.ToLowerCamel,
	}

	// TODO: Add support for many to many relationships, account for field type, i.e. array

	_ = template.Must(repositoryTemplate.New("repository-sqlite").Funcs(funcMap).Funcs(sqliteFuncMap).Parse(`
// InMemory{{.CRUD.Name}}Repository is an in memory implementation of the {{.CRUD.Name}}Repository interface.
type SQLite{{.CRUD.Name}}Repository struct {
	db *sql.DB
}

// NewInMemory creates a new InMemory{{.CRUD.Name}}Repository to be used.
func NewSQLite{{.CRUD.Name}}Repository(db *sql.DB) (*SQLite{{.CRUD.Name}}Repository, error) {
	_, ok := db.Driver().(*sqlite.Driver)
	if !ok {
		return nil, fmt.Errorf("invalid driver, must be of type *modernc.org/sqlite.Driver")
	}
	return &SQLite{{.CRUD.Name}}Repository{
		db: db,
	}, nil
}

{{if .CRUD.Create}}
// Create creates new {{.CRUD.Name}}s.
// Successfully created {{.CRUD.Name}}s are returned along with any errors that may have occurred.
func (repo *SQLite{{.CRUD.Name}}Repository) Create(ctx context.Context, toCreate []*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}) ([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, error) {
	{{if eq (len .CRUD.UniqueIdentifiers) 0 -}}
	panic("cannot create: no fields defined")
	{{else -}}
	if len(toCreate) == 0 {
		return nil, nil
	}
	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	{{if eq .CRUD.FieldMaskFieldName "" -}}
	binds := []any{}
	bindsStrs := []string{}
	for _, {{toLowerCamel $.CRUD.GetName}} := range toCreate {
		{{- range $field := .CRUD.DataFields}}
		binds = append(binds, {{toLowerCamel $.CRUD.GetName}}.{{sqliteMemberAccessor $field}})
		{{- end}}
		bindsStrs = append(bindsStrs, "(
			{{- range $i, $field := .CRUD.DataFields -}}
			{{if $i}},{{end}}?
			{{- end -}})")
	}
	_, err = tx.ExecContext(
		ctx,
		fmt.Sprintf(
			"INSERT INTO {{sqliteIdent (sqliteTableName .CRUD.GetName)}} (
			{{- range $i, $field := .CRUD.DataFields -}}
			{{if $i}},{{end}}{{sqliteIdent (sqliteColumnName $field.GetName)}}
			{{- end -}}) VALUES \n %s",
			strings.Join(bindsStrs, ",\n"),
		),
		binds...
	)
	if err != nil {
		return nil, err
	}
	{{- else -}}
	noMaskBinds := []any{}
	noMaskBindsStrs := []string{}
	for _, {{toLowerCamel $.CRUD.GetName}} := range toCreate {
		if {{toLowerCamel $.CRUD.GetName}}.{{camelIdentifier $.CRUD.FieldMaskFieldName}} == nil {
			{{- range $field := .CRUD.DataFields}}
			noMaskBinds = append(noMaskBinds, {{toLowerCamel $.CRUD.GetName}}.{{sqliteMemberAccessor $field}})
			{{- end}}
			noMaskBindsStrs = append(noMaskBindsStrs, "(
			{{- range $i, $field := .CRUD.DataFields -}}
			{{if $i}},{{end}}?
			{{- end -}})")
			continue
		}
		valuesByColName, err := sqlite{{.CRUD.Name}}GetCreateValuesByColumnName({{toLowerCamel $.CRUD.GetName}}, {{toLowerCamel $.CRUD.GetName}}.{{camelIdentifier $.CRUD.FieldMaskFieldName}})
		if err != nil {
			return nil, err
		}
		if len(valuesByColName) == 0 {
			continue
		}
		var binds []any
		var cols []string
		var params []string
		for colName, value := range valuesByColName {
			cols = append(cols, "\"" + colName + "\"")
			params = append(params, "?")
			binds = append(binds, value)
		}
		_, err = tx.ExecContext(
			ctx,
			fmt.Sprintf(
				"INSERT INTO {{sqliteIdent (sqliteTableName .CRUD.GetName)}} (%s) VALUES \n (%s)",
				strings.Join(cols, ", "),
				strings.Join(params, ", "),
			),
			binds...
		)
		if err != nil {
			return nil, err
		}
	}
	if len(noMaskBinds) > 0 {
		_, err = tx.ExecContext(
			ctx,
			fmt.Sprintf(
				"INSERT INTO {{sqliteIdent (sqliteTableName .CRUD.GetName)}} (
				{{- range $i, $field := .CRUD.DataFields -}}
				{{if $i}},{{end}}{{sqliteIdent (sqliteColumnName (sqliteColumnNameFromFieldName $field))}}
				{{- end -}}) VALUES \n %s",
				strings.Join(noMaskBindsStrs, ",\n"),
			),
			noMaskBinds...
		)
		if err != nil {
			return nil, err
		}
	}
	{{ end }}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return toCreate, nil
	{{- end }}
}
{{end}}

{{if .CRUD.Read}}
// Read returns a set of {{.CRUD.Name}}s matching the provided criteria
// Read is incomplete and it should be considered unstable
func (repo *SQLite{{.CRUD.Name}}Repository) Read(ctx context.Context, expr expressions.Expression) ([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, error) {
	query := "SELECT {{ range $i, $field := .CRUD.DataFields -}}
		{{if $i}},{{end}}{{sqliteIdent (sqliteColumnName (sqliteColumnNameFromFieldName $field))}}
		{{- end }} FROM {{sqliteIdent (sqliteTableName .CRUD.GetName)}}"
	clauses, binds, err := whereClauseFromExpressionFor{{.CRUD.Name}}(expr)
	if err != nil {
		return nil, err
	}
	if clauses != "" {
		query += "\nWHERE\n" + clauses
	}
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.QueryContext(ctx, binds...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var found []*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}
	for rows.Next() {
		{{toLowerCamel $.CRUD.GetName}} := &{{.CRUD.GoType .CRUD.File.GoPkg.Path}}{
			{{range $i, $field := .CRUD.DataFields -}}
			{{if $field.HasRelationship}}{{camelIdentifier $field.GetName}}: &{{$field.FieldMessage.GoType $.CRUD.File.GoPkg.Path}}{},{{end}}
			{{- end}}
		}
		if err = rows.Scan(
		{{- range $i, $field := .CRUD.DataFields -}}
		{{if $i}},{{end}} &{{toLowerCamel $.CRUD.GetName}}.{{sqliteMemberField $field}}
		{{- end -}}
		); err != nil {
			return nil, err
		}
		found = append(found, {{toLowerCamel $.CRUD.GetName}})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return found, nil
}
{{end}}

{{if .CRUD.Update}}
// Update modifies existing {{.CRUD.Name}}s based on the defined unique identifiers.
func (repo *SQLite{{.CRUD.Name}}Repository) Update(ctx context.Context, toUpdate []*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}) ([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, error) {
	{{if eq (len .CRUD.MinimalUIDFields) 0 -}}
	panic("cannot update: no unique identifiers defined")
	{{else if eq (len .CRUD.NonMinimalUIDDataFields) 0 -}}
	panic("cannot update: all fields are part of the minimal unique identifier: use create and delete instead")
	{{else}}
	if len(toUpdate) == 0 {
		return nil, nil
	}
	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		"UPDATE {{sqliteIdent (sqliteTableName .CRUD.GetName)}} SET {{range $i, $field := .CRUD.NonMinimalUIDDataFields -}}
		{{if $i}},{{end}}{{sqliteIdent (sqliteColumnName (sqliteColumnNameFromFieldName $field))}} = ?
		{{- end }} WHERE {{ range $i, $field := .CRUD.MinimalUIDFields -}}
		{{if $i}},{{end}}{{sqliteIdent (sqliteColumnName $field.GetName)}} = ?
		{{- end }}",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	{{ if eq .CRUD.FieldMaskFieldName ""}}
	for _, {{toLowerCamel $.CRUD.GetName}} := range toUpdate {
		_, err = stmt.ExecContext(ctx, {{ range $i, $field := .CRUD.NonMinimalUIDDataFields -}}
		{{if $i}},{{end}}{{toLowerCamel $.CRUD.GetName}}.{{sqliteMemberAccessor $field}}
		{{- end }},{{ range $i, $field := .CRUD.MinimalUIDFields -}}
		{{if $i}},{{end}}{{toLowerCamel $.CRUD.GetName}}.{{sqliteMemberAccessor $field}}
		{{- end }})
		if err != nil {
			return nil, err
		}
	}
	{{else}}
	for _, {{toLowerCamel $.CRUD.GetName}} := range toUpdate {
		if {{toLowerCamel $.CRUD.GetName}}.{{camelIdentifier $.CRUD.FieldMaskFieldName}} == nil {
			_, err = stmt.ExecContext(ctx, {{ range $i, $field := .CRUD.NonMinimalUIDDataFields -}}
			{{if $i}},{{end}}{{toLowerCamel $.CRUD.GetName}}.{{sqliteMemberAccessor $field}}
			{{- end }},{{ range $i, $field := .CRUD.MinimalUIDFields -}}
			{{if $i}},{{end}}{{toLowerCamel $.CRUD.GetName}}.{{sqliteMemberAccessor $field}}
			{{- end }})
			if err != nil {
				return nil, err
			}
			continue
		}
		valuesByColName, err := sqlite{{.CRUD.Name}}GetUpdateValuesByColumnName({{toLowerCamel $.CRUD.GetName}}, {{toLowerCamel $.CRUD.GetName}}.{{camelIdentifier $.CRUD.FieldMaskFieldName}})
		if err != nil {
			return nil, err
		}
		if len(valuesByColName) == 0 {
			continue
		}
		var binds []any
		var setStmts []string
		for colName, value := range valuesByColName {
			setStmts = append(setStmts, fmt.Sprintf("\"%s\" = ?", colName))
			binds = append(binds, value)
		}
		_, err = tx.ExecContext(
			ctx,
			fmt.Sprintf(
				"UPDATE {{sqliteIdent (sqliteTableName .CRUD.GetName)}} SET %s WHERE {{ range $i, $field := .CRUD.MinimalUIDFields -}}
				{{if $i}},{{end}}{{sqliteIdent (sqliteColumnName $field.GetName)}} = ?
				{{- end }}",
				strings.Join(setStmts, ", "),
			),
			append(
				binds,
				{{ range $i, $field := .CRUD.MinimalUIDFields -}}
				{{if $i}},{{end}}{{toLowerCamel $.CRUD.GetName}}.Get{{camelIdentifier $field.GetName}}()
				{{- end }},
			)...
		)
		if err != nil {
			return nil, err
		}
	}
	{{ end }}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return toUpdate, nil
	{{end}}
}
{{end}}

{{if .CRUD.Delete}}
// Delete deletes {{.CRUD.Name}}s based on the defined unique identifiers
func (repo *SQLite{{.CRUD.Name}}Repository) Delete(ctx context.Context, expr expressions.Expression) error {
	query := "DELETE FROM {{sqliteIdent (sqliteTableName .CRUD.GetName)}}"
	clauses, binds, err := whereClauseFromExpressionFor{{.CRUD.Name}}(expr)
	if err != nil {
		return err
	}
	if clauses != "" {
		query += "\nWHERE\n" + clauses
	}
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, binds...)
	if err != nil {
		return err
	}
	return nil
}
{{end}}

{{if or .CRUD.Read .CRUD.Delete}}
var sqlite{{.CRUD.Name}}ColumnNameByFieldID = map[expressions.FieldID]string{
{{- range $name, $data := .FieldByFieldConstants}}
	{{$name}}: "{{sqliteColumnName $data.Def.GetName}}",
{{- end}}
}

func whereClauseFromExpressionFor{{.CRUD.Name}}(expr expressions.Expression) (string, []any, error) {
	if expr == nil {
		return "", nil, nil
	}
	switch expr := expr.(type) {
		case *expressions.And:
			left, leftBinds, err := whereClauseFromExpressionFor{{.CRUD.Name}}(expr.Left())
			if err != nil {
				return "", nil, err
			}
			right, rightBinds, err := whereClauseFromExpressionFor{{.CRUD.Name}}(expr.Right())
			if err != nil {
				return "", nil, err
			}
			return fmt.Sprintf("%s AND %s", left, right), append(leftBinds, rightBinds...), nil

		case *expressions.Or:
				left, leftBinds, err := whereClauseFromExpressionFor{{.CRUD.Name}}(expr.Left())
			if err != nil {
				return "", nil, err
			}
			right, rightBinds, err := whereClauseFromExpressionFor{{.CRUD.Name}}(expr.Right())
			if err != nil {
				return "", nil, err
			}
			return fmt.Sprintf("%s OR %s", left, right), append(leftBinds, rightBinds...), nil
		case *expressions.Not:
			operand, binds, err := whereClauseFromExpressionFor{{.CRUD.Name}}(expr.Operand())
			if err != nil {
				return "", nil, err
			}
			return fmt.Sprintf("NOT %s", operand), binds, nil

		case *expressions.Equal:
				left, leftBinds, err := whereClauseFromExpressionFor{{.CRUD.Name}}(expr.Left())
			if err != nil {
				return "", nil, err
			}
			right, rightBinds, err := whereClauseFromExpressionFor{{.CRUD.Name}}(expr.Right())
			if err != nil {
				return "", nil, err
			}
			return fmt.Sprintf("%s = %s", left, right), append(leftBinds, rightBinds...), nil

		case *expressions.Identifier:
			{{if .FieldByFieldConstants}}if _, ok := valid{{.CRUD.Name}}Fields[expr.ID()]; !ok {
				return "", nil, fmt.Errorf("invalid field id: %s", expr.ID())
			}
			colName, ok := sqlite{{.CRUD.Name}}ColumnNameByFieldID[expr.ID()]
			if !ok {
				return "", nil, fmt.Errorf("missing meta-data: field id: %s", expr.ID())
			}
			return fmt.Sprintf("{{sqliteIdent (sqliteTableName .CRUD.GetName)}}.\"%s\"",colName), nil, nil
			{{else}}
			return "", nil, fmt.Errorf("identifiers not supported")
			{{end}}
		case *expressions.Scalar:
			return "?", []any{expr.Value()}, nil
		default:
			return "", nil, fmt.Errorf("unknown expression")
	}
}
{{end}}
{{if and (or .CRUD.Update .CRUD.Create) (ne .CRUD.FieldMaskFieldName "")}}
func sqlite{{.CRUD.Name}}GetCreateValuesByColumnName(def *{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, fieldMask *fieldmaskpb.FieldMask) (map[string]any, error) {
	if fieldMask == nil {
		return nil, fmt.Errorf("no field mask provided")
	}
	{{toLowerCamel $.CRUD.GetName}} := &{{.CRUD.GoType .CRUD.File.GoPkg.Path}}{}
	valuesByColumnName := make(map[string]any, 0)
	nestedMask := fmutils.NestedMaskFromPaths(fieldMask.Paths)
	{{ range $i, $field := .CRUD.DataFields -}}
	if _, ok := nestedMask["{{$field.GetName}}"]; ok {
		valuesByColumnName["{{sqliteColumnName (sqliteColumnNameFromFieldName $field)}}"] = def.{{sqliteMemberAccessor $field}}
	} else {
		valuesByColumnName["{{sqliteColumnName (sqliteColumnNameFromFieldName $field)}}"] = {{toLowerCamel $.CRUD.GetName}}.{{sqliteMemberAccessor $field}}
	}
	{{end -}}
	return valuesByColumnName, nil
}
func sqlite{{.CRUD.Name}}GetUpdateValuesByColumnName(def *{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, fieldMask *fieldmaskpb.FieldMask) (map[string]any, error) {
	if fieldMask == nil {
		return nil, fmt.Errorf("no field mask provided")
	}
	valuesByColumnName := make(map[string]any, 0)
	nestedMask := fmutils.NestedMaskFromPaths(fieldMask.Paths)
	{{ range $i, $field := .CRUD.DataFields -}}
	if _, ok := nestedMask["{{$field.GetName}}"]; ok {
		valuesByColumnName["{{sqliteColumnName (sqliteColumnNameFromFieldName $field)}}"] = def.{{sqliteMemberAccessor $field}}
	}
	{{end -}}
	return valuesByColumnName, nil
}
{{end}}
`))
)
