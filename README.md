# Protobuf Annotations

## Message

1. operations to generate
    1. []enum: create, read, update, delete
2. delete strategy
    1. enum: hard, soft
    2. default to hard delete
    3. add support for soft later
        1. delete at support?
3. use field mask for partial updates/creates if present
    1. boolean
    2. default true
4. create at
    1. boolean
    2. default false
5. updated at
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


1. unique identifier(s)
    1. auto-generated?
        1. In-memory needs to specify methodology (uuid, sequential?)
1. Every value is set/updated every time OR field mask support for partial updates/creates?
2. Nullable?
    1. Default no
    1. Implement later, shouldn't be using null anyway
3. Hard and soft delete
    1. Default to hard, add support for soft later

# To generate...

1. Repository interface definitions
    1. option functions
    2. Operations
        1. Create `([]*Message, options) ([]*[]Message, error)`
        2. Read `([]Clause) ([]*Message, error)`
            1. Clause
                1. IsEqual
                2. IsNotEqual
                3. And
                4. Or
                5. CLAUSE
        3. Update `([]*Message, options) ([]*Message, error)`
        4. Delete `([]*Message) ([]*Message, error)`
2. Repository interface implementations

3. SQL statements to create tables
4. DB's to support, in order
    1. SQLite
    2. PgSQL

# Try?

1. Earthly, https://github.com/earthly/earthly

# References

1. https://go.dev/blog/protobuf-apiv2
2. https://golang.design/research/generic-option/
3. https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
4. https://github.com/googleapis/googleapis/blob/master/google/api/http.proto
5. https://github.com/grpc-ecosystem/grpc-gateway/tree/main/protoc-gen-openapiv2/options
