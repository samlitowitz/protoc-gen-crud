# Protobuf Annotations

## Message

1. operations to generate
    1. []enum: create, read, update, delete
2. implementations to generate
    1. []enum: in-memory, sqlite, pgsql
3. delete strategy
    1. enum: hard, soft
    2. default to hard delete
    3. add support for soft later
        1. deleted at support?
4. use field mask for partial updates/creates
    1. boolean
    2. default true
    3. ignored if no field mask field defined on message
5. created at
    1. boolean
    2. default false
6. updated at
    1. boolean
    2. default false

## Field

1. unique identifiers(s)
    1. []string
2. auto-generate strategy
    1. enum: none, uuid, sequential integer
    2. default none
3. nullable
    1. boolean
    2. default false
4. non-scalar field type strategy
    1. enum: inline, link via unique id, skip

### Adding field option in code

1. `options/crud.proto`
    1. `message FieldOptions`
        1. Add field option(s)
2. `internal/descriptor/types.go`
    1. `type CRUD struct`
        1. Add ability to hand field option(s)
3. `internal/descriptor/cruds.go`
    1. Parse/set field option(s)
4. `internal/descriptor/cruds_test.go`
    1. Add unit test for field option(s)

# To generate...

1. Repository Interface

```go

// source: <PROTO_SOURCE_FILE> -> `protogen.File.Desc.Path()`

// <GO_PACKAGE_NAME> -> `protogen.File.GoPackageName`
package <GO_PACKAGE_NAME>

// <MESSAGE_NAME> -> `protogen.Message.GoIdent`
type
<MESSAGE_NAME>Repository interface {
// <QUALIFIED_MESSAGE_TYPE> -> Appropriately qualified `protogen.Message.GoIdent` (current package, imported, aliased)
Create([]*<QUALIFIED_MESSAGE_TYPE>) ([]*MessageType, error)
Read() ([]*<QUALIFIED_MESSAGE_TYPE>, error)
Update([]*<QUALIFIED_MESSAGE_TYPE>) ([]*<QUALIFIED_MESSAGE_TYPE>, error)
Delete([]*<QUALIFIED_MESSAGE_TYPE>) error
}
```

1. Repository interface definitions
    1. option functions
    2. Operations
        1. Create `([]*Message, options) ([]*Message, error)`
        2. Read `([]Clause) ([]*Message, error)`
            1. Revisit or start simple
            1. Clause
                1. And
                2. Or
                3. Not
                4. Null
                5. Equal
                6. GreaterThanOrEqual
                7. LessThanOrEqual
                8. Like
                9. In
        3. Update `([]*Message, options) ([]*Message, error)`
        4. Delete `([]*Message) ([]*Message, error)`
    3. fully qualified message names
2. Repository interface implementations
    1. In Memory
        1.
    2. SQLite
    3. PgSQL
3. Auxiliary Features
    1. SQL statements to create tables
        1. SQLite
        2. PgSQL

## Clauses

### Accept expressions

1. And
2. Or
3. Not

### Accept nothing

1. Null

### Accept field expressions

1. Equal
2. GreaterThanOrEqual
3. LessThanOrEqual
4. GreaterThan
5. LessThan
6. Like
7. In

```go
package tmp

type TypID string

const (
    BookTypID TypID = "asfwer234sdfadsf" // hash of the fully qualified book type, `github.com/path/to/Book`
)


```

---

1. Generate map of field names to column names
2. Generate function to all field names belong to the type
    3. Generic function which takes []ValidFields, []Fields

```go
package tmp

import (
    "fmt"
    "strings"
)

type FieldName string

func ValidateFieldNames(validFieldNames map[FieldName]struct{}, fieldNames []FieldName) error {
    invalidFieldNames := make([]string, 0)
    for _, fieldName := range fieldNames {
        if _, ok := validFieldNames[fieldName]; ok {
            continue
        }
        invalidFieldNames = append(invalidFieldNames, string(fieldName))
    }
    if len(invalidFieldNames) == 0 {
        return nil
    }
    return fmt.Errorf("Invalid fields: %s", strings.Join(invalidFieldNames, ", "))
}
```

4. Generate type specific function which wraps generic function
1. `Read` implementation validates field names belong to the type

1. Avoid incorrect column names at compile time OR before executing the query
2. Avoid incorrect value types not matching column types at compile time

```go
package test

type ColumnName string

func EqualInt(columnName string, value int) {
    // ...
}

func InInt(columnName string, values ...int) {
    // ...
}

func LessThanInt(columnName string, value int) {
    // ...
}
```

# References

1. https://go.dev/blog/protobuf-apiv2
2. https://golang.design/research/generic-option/
3. https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
4. https://github.com/googleapis/googleapis/blob/master/google/api/http.proto
5. https://github.com/grpc-ecosystem/grpc-gateway/tree/main/protoc-gen-openapiv2/options
