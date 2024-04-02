All test cases must include the following

1. A `generate.go` file which runs `protoc` with the appropriate flags and arguments when the `go generate`
   from `protoc-gen-crud`'s root directory
2. A set of `proto` files from which the code to be tested is generated
3. A set of `_test.go` files which tests the generated code
4. A `.gitignore` file which ignores all files except itself and the files explicitly mentioned in the previous points

- TODO: create template/generator for new test case base directory?

# Test Cases

1. With no implementations set, only the repository interface should be generated
2. With only the `UNKNOWN` implementation set, only the repository interface should be generated

## Implementation Agnostic

These test cases must be covered for _EVERY_ implementation.

### Message

#### Primary Key

Prime attributes may be an enum, any integer type, or a string.
All primary key tests _**MUST**_ cover all allowed attribute types.

1. Create
    1. A new entity with a duplicate primary key set **MUST** fail
    2. A new entity with a non-duplicate primary key **MUST** succeed
    3. [FEATURE] A new entity with no prime attributes set **MUST** have them generated
2. Read
    1. No applicable requirements
3. Update
    1. Updating an entity with an un-locatable primary key **MUST** not update any entity
    2. Updating an entity with a locatable primary key **MUST** succeed
4. Delete
    1. Calling delete with no expression **MUST** not delete any entity
    2. Calling delete with an expression **MUST** delete any entity matching that expression
        1. Primary key
        2. Non-prime attributes

#### Field Mask

1. Create
    1. A new entity with any prime attribute excluded by field mask **MUST** fail
    2. A new entity with any or all non-prime attributes excluded by field mask **MUST** use empty values for those attributes
    3. A new entity with no field mask used **MUST** succeed
2. Read
    1. No applicable requirements
3. Update
    1. Updating an entity with any prime attribute excluded by field mask **MUST** fail
    2. Updating an entity with any or all non-prime attributes excluded by field mask **MUST** not modify those values
    3. Updating an entity with no field mask used **MUST** modify all non-prime attributes
4. Delete
    1. No applicable requirements


### Field

#### Ignore
Ignored fields are to be excluded from all generated CRUD code

TODO: How to test?
1. Ignored fields are not supported by any CRUD operations
2. Including in ignored field in the primary key **MUST** fail to compile

#### Relationships
1. manual/managed/??? - TODO: need correct nomenclature for this

TODO: Need to write test cases/expected behaviors

1. Always generate join types and CRUD code
1. Managed 

...

* Message
    * Field Mask
        * Partial mask of candidate keys prime attributes
    * Candidate Key
        * Supported types
            * string
            * ints
            * enums
        * Single attribute
        * Multi-attribute
* Field
    * Ignore
    * No Generate (do not generate candidate key attributes, this is ignored on all other fields)
    * Relationships
        * type X direction

## Implementation Specific

### SQLite

# References

1. Wikipedia contributors. (2024, March 13). Candidate key. Wikipedia. https://en.wikipedia.org/wiki/Candidate_key
