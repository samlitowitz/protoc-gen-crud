package relationships_many_to_many_test

import (
	"testing"

	relationships_many_to_many "github.com/samlitowitz/protoc-gen-crud/test-cases/relationships-many-to-many"
)

// saInt32ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saInt32ComponentUnderTest func(t *testing.T) relationships_many_to_many.SAInt32Repository

// maAllComponentUnderTest is to be implemented to do setup and tear down for each implementation
type maAllComponentUnderTest func(t *testing.T) relationships_many_to_many.MAAllRepository
