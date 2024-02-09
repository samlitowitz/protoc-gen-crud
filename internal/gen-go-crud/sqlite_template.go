package gen_go_crud

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"

	"github.com/iancoleman/strcase"
)

func init() {
	strcase.ConfigureAcronym("UID", "uid")
}

func SQLiteIdent(s string) string {
	return "\"" + s + "\""
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

func (sqlite *sqlite) CreateQuery(crud *descriptor.CRUD) string {
	if len(crud.Fields) == 0 {
		return ""
	}
	query := `INSERT INTO %s (%s) VALUES%s%%s`
	cols := make([]string, 0, len(crud.Fields))
	for _, def := range crud.Fields {
		cols = append(cols, SQLiteColumnIdentifier(def.GetName()))
	}
	return strconv.Quote(fmt.Sprintf(
		query,
		SQLiteIdent(SQLiteTableName(crud.GetName())),
		strings.Join(cols, ", "),
		"\n",
	))
}

var (
	_ = template.Must(repositoryTemplate.New("repository-sqlite").Funcs(funcMap).Parse(`
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
	if len(toCreate) == 0 {
		return nil, nil
	}
	binds := []any{}
	bindsStrs := []string{}
	for i, val := range toCreate {
		{{- range $constName, $data := .FieldByFieldConstants}}
		binds = append(
			binds,
			sql.Named(
				fmt.Sprintf(":{{toLower $constName}}_%d", i),
				val.{{ camelIdentifier $data.Def.GetName }},
			),
		)
		bindsStrs = append(bindsStrs, fmt.Sprintf(":{{toLower $constName}}_%d", i))
		{{- end}}
	}
	query := fmt.Sprintf(
		{{ .SQLite.CreateQuery .CRUD }},
		strings.Join(bindsStrs, ",\n"),
	)
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(binds...)
	if err != nil {
		return nil, err
	}
	return toCreate, nil
}
{{end}}

{{if .CRUD.Read}}
// Read returns a set of {{.CRUD.Name}}s matching the provided criteria
// Read is incomplete and it should be considered unstable
func (repo *SQLite{{.CRUD.Name}}Repository) Read(ctx context.Context, expr expressions.Expression) ([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, error) {
	panic("not implemented")
}
{{end}}

{{if .CRUD.Update}}
// Update modifies existing {{.CRUD.Name}}s based on the defined unique identifiers.
func (repo *SQLite{{.CRUD.Name}}Repository) Update(ctx context.Context, toUpdate []*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}) ([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, error) {
	panic("not implemented")
}
{{end}}

{{if .CRUD.Delete}}
// Delete deletes {{.CRUD.Name}}s based on the defined unique identifiers
func (repo *SQLite{{.CRUD.Name}}Repository) Delete(ctx context.Context, expr expressions.Expression) error {
	panic("not implemented")
}
{{end}}
`))
)
