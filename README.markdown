# Protobuf Annotations

# Table of Contents

1. [Feature Support](#feature-support)
    1. [Message](#message)
        1. [Operations](#operations)
        2. [Delete Strategy](#delete-strategy)
        3. [Partial Creates/Updates](#partial-createsupdates)
        4. [Meta-Data](#meta-data)
        5. [Audit Logging](#audit-logging)
    2. [Field](#field)
        1. [Unique Identifiers](#unique-identifiers)
        2. [Auto-generate Strategy](#auto-generate-strategy)
        3. [Nullable](#nullable)
        4. [Non-scalar Field Relationship Strategy](#non-scalar-field-relationship-strategy)
    3. [References](#references)

# Feature Support

## Message

### Operations

| Implementation | Create             | Read               | Update             | Delete             |
|:---------------|:-------------------|:-------------------|:-------------------|:-------------------|
| SQLite         | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| PgSQL          | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |

### Delete Strategy

| Implementation | Hard               | Soft |
|:---------------|:-------------------|:-----|
| SQLite         | :white_check_mark: |      |
| PgSQL          | :white_check_mark: |      |

### Partial Creates/Updates

| Implementation | Field Mask         |
|:---------------|:-------------------|
| SQLite         | :white_check_mark: |
| PgSQL          | :white_check_mark: |

### Meta-Data

| Implementation | Created At | Updated At | Deleted At |
|:---------------|:-----------|:-----------|:-----------|
| SQLite         |            |            |            |
| PgSQL          |            |            |            |

### Audit Logging

| Implementation |
|:---------------|
| SQLite         |
| PgSQL          |

## Field

### Unique Identifiers

| Implementation | Single             | Composite |
|:---------------|:-------------------|:----------|
| SQLite         | :white_check_mark: | :white_check_mark: |
| PgSQL          | :white_check_mark: | :white_check_mark: |

### As Timestamp

| Implementation | google.protobuf.Timestamp |
|:---------------|-------------------|
| SQLite         | :white_check_mark: |
| PgSQL          | :white_check_mark: |

### Auto-generate Strategy

| Implementation | None | UUID | Sequential Integer |
|:---------------|:-----|:-----|:-------------------|
| SQLite         |   :white_check_mark:   |      |                    |
| PgSQL          |   :white_check_mark:   |      |                    |

### Nullable

| Implementation | Nullable                  | Not-Nullable |
|:---------------|:--------------------------|:-------------|
| SQLite         | :white_check_mark:        |              |
| PgSQL          |                          |              |

### Non-scalar Fields

| Implementation | Skip | Inline             | Relationship (see below) |
|:---------------|:-----|:-------------------|:-------------------------|
| SQLite         | :white_check_mark:     | :white_check_mark: | -                        |
| PgSQL          | :white_check_mark:     | :white_check_mark: | -                        |

#### Relationship

1. direction
1. bidirectional
2. unidirectional
2. type
1. one-to-one
2. one-to-many
3. many-to-one

##### Unidirectional

| Implementation | One-to-one          | One-to-many | Many-to-many |
|:---------------|:--------------------|:------------|:-------------|
| SQLite         | :white_check_mark:  |             |              |
| PgSQL          | :white_check_mark:  |             |              |

# References

1. https://go.dev/blog/protobuf-apiv2
2. https://golang.design/research/generic-option/
3. https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
4. https://github.com/googleapis/googleapis/blob/master/google/api/http.proto
5. https://github.com/grpc-ecosystem/grpc-gateway/tree/main/protoc-gen-openapiv2/options
6. https://go.dev/blog/intro-generics
7. https://pkg.go.dev/modernc.org/sqlite
