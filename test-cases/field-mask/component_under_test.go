package field_mask_test

import (
	"testing"

	fieldMask "github.com/samlitowitz/protoc-gen-crud/test-cases/field-mask"
)

// saEnumComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saEnumComponentUnderTest func(t *testing.T) fieldMask.SAEnumRepository

// saInt32ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saInt32ComponentUnderTest func(t *testing.T) fieldMask.SAInt32Repository

// saInt64ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saInt64ComponentUnderTest func(t *testing.T) fieldMask.SAInt64Repository

// saUint32ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saUint32ComponentUnderTest func(t *testing.T) fieldMask.SAUint32Repository

// saUint64ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saUint64ComponentUnderTest func(t *testing.T) fieldMask.SAUint64Repository

// saStringComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saStringComponentUnderTest func(t *testing.T) fieldMask.SAStringRepository

// maAllComponentUnderTest is to be implemented to do setup and tear down for each implementation
type maAllComponentUnderTest func(t *testing.T) fieldMask.MAAllRepository
