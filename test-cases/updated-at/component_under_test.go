package updated_at_test

import (
	"testing"

	updated_at "github.com/samlitowitz/protoc-gen-crud/test-cases/updated-at"
)

// updatedAtComponentUnderTest is to be implemented to do setup and tear down for each implementation
type updatedAtComponentUnderTest func(t *testing.T) updated_at.UpdatedAtRepository
