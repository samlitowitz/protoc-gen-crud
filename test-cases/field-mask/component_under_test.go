package field_mask_test

import (
	"testing"

	fieldMask "github.com/samlitowitz/protoc-gen-crud/test-cases/field-mask"
)

// saInt32ComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saInt32ComponentUnderTest func(t *testing.T) fieldMask.SAInt32Repository

// maAllComponentUnderTest is to be implemented to do setup and tear down for each implementation
type maAllComponentUnderTest func(t *testing.T) fieldMask.MAAllRepository
