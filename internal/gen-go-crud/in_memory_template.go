package gen_go_crud

import "text/template"

var (
	_ = template.Must(repositoryTemplate.New("repository-in-memory").Funcs(funcMap).Parse(`

// TODO: Add Comment
type InMemory{{.CRUD.Name}}Repository struct {}

// TODO: Add Comment
func NewInMemory{{.CRUD.Name}}Repository() *InMemory{{.CRUD.Name}}Repository {
	return &InMemory{{.CRUD.Name}}Repository{}
}

{{if .CRUD.Create}}
// TODO: Add Comment
func (repo *InMemory{{.CRUD.Name}}Repository) Create([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}) ([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, error) {

}
{{end}}

{{if .CRUD.Read}}
// TODO: Add Comment
// Read is incomplete and it should be considered unstable
func (repo *InMemory{{.CRUD.Name}}Repository) Read() ([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, error) {

}
{{end}}

{{if .CRUD.Update}}
// TODO: Add Comment
func (repo *InMemory{{.CRUD.Name}}Repository) Update([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}) ([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}, error) {

}
{{end}}

{{if .CRUD.Delete}}
// TODO: Add Comment
func (repo *InMemory{{.CRUD.Name}}Repository) Delete([]*{{.CRUD.GoType .CRUD.File.GoPkg.Path}}) error {

}
{{end}}
`))
)
