package candidate_key_test

import (
	"testing"

	candidateKey "github.com/samlitowitz/protoc-gen-crud/test-cases/candidate-key"
)

// saEnumComponentUnderTest is to be implemented to do setup and tear down for each implementation
type saEnumComponentUnderTest func(t *testing.T) candidateKey.SAEnumRepository
