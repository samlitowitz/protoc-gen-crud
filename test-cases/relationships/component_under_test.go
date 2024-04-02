package relationships_test

import (
	"testing"

	"github.com/samlitowitz/protoc-gen-crud/test-cases/relationships"
)

// saEnumComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saEnumComponentUnderTest func(t *testing.T) relationships.SAEnumRepository

// saInt32ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saInt32ComponentUnderTest func(t *testing.T) relationships.SAInt32Repository

// saInt64ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saInt64ComponentUnderTest func(t *testing.T) relationships.SAInt64Repository

// saUint32ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saUint32ComponentUnderTest func(t *testing.T) relationships.SAUint32Repository

// saUint64ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saUint64ComponentUnderTest func(t *testing.T) relationships.SAUint64Repository

// saStringComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saStringComponentUnderTest func(t *testing.T) relationships.SAStringRepository

// maAllComponentUnderTest is to be implemented to do setup and tear down for each implementation
type maAllComponentUnderTest func(t *testing.T) relationships.MAAllRepository
