package inline_field_test

import (
	inline_field "github.com/samlitowitz/protoc-gen-crud/test-cases/inline-field"
	"testing"
)

// inlineTimestampComponentUnderTest is to be implemented to do setup and tear down for each implementation
type inlineTimestampComponentUnderTest func(t *testing.T) inline_field.InlineTimestampRepository
