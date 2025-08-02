package relationships_one_to_one_test

import (
	"testing"

	relationships_one_to_one "github.com/samlitowitz/protoc-gen-crud/test-cases/relationships-one-to-one"
)

// saInt32ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saInt32ComponentUnderTest func(t *testing.T) relationships_one_to_one.SAInt32Repository

// maAllComponentUnderTest is to be implemented to do setup and tear down for each implementation
type maAllComponentUnderTest func(t *testing.T) relationships_one_to_one.MAAllRepository
