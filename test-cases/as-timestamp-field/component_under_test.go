package as_timestamp_field_test

import (
	"testing"

	as_timestamp_field "github.com/samlitowitz/protoc-gen-crud/test-cases/as-timestamp-field"
)

// asTimestampComponentUnderTest is to be implemented to do setup and tear down for each implementation
type asTimestampComponentUnderTest func(t *testing.T) as_timestamp_field.AsTimestampRepository
