package gen_go_crud

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/samlitowitz/protoc-gen-crud/internal/casing"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"

	"github.com/iancoleman/strcase"
)

func init() {
	strcase.ConfigureAcronym("UID", "uid")
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

func SQLiteColumnIdentifier(s string) string {
	return SQLiteIdent(SQLiteColumnName(s))
}

type sqlite struct{}

func (sqlite *sqlite) UpdateBinds(crud *descriptor.CRUD) string {
	bindsStrs := make([]string, 0, len(crud.Fields))
	for _, def := range crud.Fields {
		bindsStrs = append(
			bindsStrs,
			fmt.Sprintf(
				"%s.%s",
				crud.CamelCaseName(),
				casing.CamelIdentifier(def.GetName()),
			),
		)
	}
	_, uidFields := getMinimalUID(crud)
	for _, def := range uidFields {
		bindsStrs = append(
			bindsStrs,
			fmt.Sprintf(
				"%s.%s",
				crud.CamelCaseName(),
				casing.CamelIdentifier(def.GetName()),
			),
		)
	}
	return strings.Join(bindsStrs, ", ")
}

func getMinimalUID(crud *descriptor.CRUD) (string, []*descriptor.Field) {
	var minUIDName string
	var minUIDFields []*descriptor.Field

	for name, fields := range crud.UniqueIdentifiers {
		if minUIDFields == nil {
			minUIDName = name
			minUIDFields = fields
			continue
		}
		if len(minUIDFields) > len(fields) {
			minUIDName = name
			minUIDFields = fields
		}
	}
	return minUIDName, minUIDFields
}

func (sqlite *sqlite) TableName(crud *descriptor.CRUD) string {
	return SQLiteTableName(crud.GetName())
}

func (sqlite *sqlite) TableIdentifierName(crud *descriptor.CRUD) string {
	return strconv.Quote(SQLiteIdent(SQLiteTableName(crud.GetName())))
}

func (sqlite *sqlite) ColumnName(field *descriptor.Field) string {
	return SQLiteColumnName(field.GetName())
}

var (
	sqliteFuncMap template.FuncMap = map[string]interface{}{
		"sqliteIdent":      SQLiteTemplateIdent,
		"sqliteTableName":  SQLiteTableName,
		"sqliteColumnName": SQLiteColumnName,
	}

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
	{{if ne (len .CRUD.UniqueIdentifiers) 0}}if len(toCreate) == 0 {
		return nil, nil
	}
	binds := []any{}
	bindsStrs := []string{}
	for _, val := range toCreate {
		{{- range $field := .CRUD.DataFields}}
		binds = append(binds, val.{{ camelIdentifier $field.GetName }})
		{{- end}}
		bindsStrs = append(bindsStrs, "(
			{{- range $i, $field := .CRUD.DataFields -}}
			{{if $i}},{{end}}?
			{{- end -}})")
	}
	stmt, err := repo.db.Prepare(
		fmt.Sprintf(
			"INSERT INTO {{sqliteIdent (sqliteTableName .CRUD.GetName)}} (
			{{- range $i, $field := .CRUD.DataFields -}}
			{{if $i}},{{end}}{{sqliteIdent (sqliteColumnName $field.GetName)}}
			{{- end -}}) VALUES \n %s",
			strings.Join(bindsStrs, ",\n"),
		),
	)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(ctx, binds...)
	if err != nil {
		return nil, err
	}
	return toCreate, nil
	{{else}}
	panic("cannot create: no fields defined")
	{{end}}
}
{{end}}

{{if .CRUD.Read}}
// Read returns a set of {{.CRUD.Name}}s matching the provided criteria
// Read is incomplete and it should be considered unstable
func (repo *SQLite{{.CRUD.Name}}Repository) Read(ctx context.Context, expr expressions.Expression) ([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, error) {
	query := "SELECT {{ range $i, $field := .CRUD.DataFields -}}
		{{if $i}},{{end}}{{sqliteIdent (sqliteColumnName $field.GetName)}}
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
		{{$.CRUD.CamelCaseName}} := &{{.CRUD.GoType .CRUD.File.GoPkg.Path}}{}
		if err = rows.Scan(
		{{- range $i, $field := .CRUD.DataFields -}}
		{{if $i}},{{end}} &{{$.CRUD.CamelCaseName}}.{{camelIdentifier $field.GetName}}
		{{- end -}}
		); err != nil {
			return nil, err
		}
		found = append(found, {{$.CRUD.CamelCaseName}})
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
	{{if eq (len .CRUD.MinimalUIDFields) 0}}
	panic("cannot update: no unique identifiers defined")
	{{else if eq (len .CRUD.NonMinimalUIDDataFields) 0}}
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
		"UPDATE {{sqliteIdent (sqliteTableName .CRUD.GetName)}} SET {{ range $i, $field := .CRUD.NonMinimalUIDDataFields -}}
		{{if $i}},{{end}}{{sqliteIdent (sqliteColumnName $field.GetName)}} = ?
		{{- end }} WHERE {{ range $i, $field := .CRUD.MinimalUIDFields -}}
		{{if $i}},{{end}}{{sqliteIdent (sqliteColumnName $field.GetName)}} = ?
		{{- end }}",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	for _, {{$.CRUD.CamelCaseName}} := range toUpdate {
		_, err = stmt.ExecContext(ctx, {{ range $i, $field := .CRUD.NonMinimalUIDDataFields -}}
		{{if $i}},{{end}}{{$.CRUD.CamelCaseName}}.{{camelIdentifier $field.GetName}}
		{{- end }},{{ range $i, $field := .CRUD.MinimalUIDFields -}}
		{{if $i}},{{end}}{{$.CRUD.CamelCaseName}}.{{camelIdentifier $field.GetName}}
		{{- end }})
		if err != nil {
			return nil, err
		}
	}
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
var sqlite{{.CRUD.Name}}FieldMetaData = map[expressions.FieldID]struct{
	tableName string
	columnName string
}{
{{- range $name, $data := .FieldByFieldConstants}}
	{{$name}}: {
		tableName: "{{$.SQLite.TableName $.CRUD}}",
		columnName: "{{$.SQLite.ColumnName $data.Def}}",
	},
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
			metaData, ok := sqlite{{.CRUD.Name}}FieldMetaData[expr.ID()]
			if !ok {
				return "", nil, fmt.Errorf("missing meta-data: field id: %s", expr.ID())
			}
			return fmt.Sprintf(
				"\"%s\".\"%s\"",
				metaData.tableName,
				metaData.columnName,
			), nil, nil
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
`))
)
