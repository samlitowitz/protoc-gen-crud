package created_at_test

import (
	"testing"

	created_at "github.com/samlitowitz/protoc-gen-crud/test-cases/created-at"
)

// createdAtComponentUnderTest is to be implemented to do setup and tear down for each implementation
type createdAtComponentUnderTest func(t *testing.T) created_at.CreatedAtRepository
