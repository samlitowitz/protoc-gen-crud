package relationships_test

import (
	"testing"

	"github.com/samlitowitz/protoc-gen-crud/test-cases/relationships"
)

// saInt32ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saInt32ComponentUnderTest func(t *testing.T) relationships.SAInt32Repository

// maAllComponentUnderTest is to be implemented to do setup and tear down for each implementation
type maAllComponentUnderTest func(t *testing.T) relationships.MAAllRepository
