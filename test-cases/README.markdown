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

#### Candidate Key

Prime attributes may be an enum, any integer type, or a string.
All candidate key tests _**MUST**_ cover all allowed attribute types.

1. Create
   1. A new entity with no prime attributes set **MUST** have them generated
   2. A new entity with a non-duplicate candidate key **MUST** succeed
   3. A new entity with a duplicate candidate key **MUST** fail
1. Read
   1. TBD
2. Update
   3. Updating an entity with a non-existent candidate key **MUST** fail
   4. Updating an entity with an existing candidate key **MUST** succeed

##### Single Attribute

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
