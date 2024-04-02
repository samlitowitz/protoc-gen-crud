package primary_key_test

import (
	"testing"

	primaryKey "github.com/samlitowitz/protoc-gen-crud/test-cases/primary-key"
)

// saEnumComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saEnumComponentUnderTest func(t *testing.T) primaryKey.SAEnumRepository

// saInt32ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saInt32ComponentUnderTest func(t *testing.T) primaryKey.SAInt32Repository

// saInt64ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saInt64ComponentUnderTest func(t *testing.T) primaryKey.SAInt64Repository

// saUint32ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saUint32ComponentUnderTest func(t *testing.T) primaryKey.SAUint32Repository

// saUint64ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saUint64ComponentUnderTest func(t *testing.T) primaryKey.SAUint64Repository

// saStringComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saStringComponentUnderTest func(t *testing.T) primaryKey.SAStringRepository

// maAllComponentUnderTest is to be implemented to do setup and tear down for each implementation
type maAllComponentUnderTest func(t *testing.T) primaryKey.MAAllRepository
