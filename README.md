# Protobuf Annotations

# Table of Contents

1. [Feature Support](#feature-support)
    1. [Message](#message)
        1. [Operations](#operations)
        2. [Delete Strategy](#delete-strategy)
        3. [Partial Updates/Creates](#partial-updatescreates)
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

### Delete Strategy

| Implementation | Hard               | Soft |
|:---------------|:-------------------|:-----|
| SQLite         | :white_check_mark: |      |

### Partial Updates/Creates

| Implementation | Field Mask |
|:---------------|:-----------|
| SQLite         |            |

### Meta-Data

| Implementation | Created At | Updated At | Deleted At |
|:---------------|:-----------|:-----------|:-----------|
| SQLite         |            |            |            |

### Audit Logging

| Implementation |
|:---------------|
| SQLite         |

## Field

### Unique Identifiers

| Implementation | Single             | Composite |
|:---------------|:-------------------|:----------|
| SQLite         | :white_check_mark: |           |:white_check_mark:

### Auto-generate Strategy

| Implementation | None | UUID | Sequential Integer |
|:---------------|:-----|:-----|:-------------------|
| SQLite         |      |      |                    |

### Nullable

| Implementation | Nullable | Not-Nullable |
|:---------------|:---------|:-------------|
| SQLite         |          |              |

### Non-scalar Field Relationship Strategy

| Implementation | Inline | Link via Unique ID | Skip |
|:---------------|:-------|:-------------------|:-----|
| SQLite         |        |                    |      |

# References

1. https://go.dev/blog/protobuf-apiv2
2. https://golang.design/research/generic-option/
3. https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
4. https://github.com/googleapis/googleapis/blob/master/google/api/http.proto
5. https://github.com/grpc-ecosystem/grpc-gateway/tree/main/protoc-gen-openapiv2/options
6. https://go.dev/blog/intro-generics
7. https://pkg.go.dev/modernc.org/sqlite
