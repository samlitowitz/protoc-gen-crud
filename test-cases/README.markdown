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

...

* Message
    * Field Mask
    * Candidate Key
      * Supported types
        * string
        * ints
        * enums
      * Single attribute
      * Multi-attribute
* Field
    * Ignore
    * Relationships
      * type X direction

## Implementation Specific

### SQLite
