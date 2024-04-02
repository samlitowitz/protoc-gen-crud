package candidate_key_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	candidateKey "github.com/samlitowitz/protoc-gen-crud/test-cases/candidate-key"
)

func TestSAEnumRepository_Create_WithMissingPrimeAttributesFails(t *testing.T) {

}

func TestSAEnumRepository_Create_WithADuplicateCandidateKeyFails(t *testing.T) {

}

func TestSAEnumRepository_Create_WithANonDuplicateCandidateKeySucceeds(t *testing.T) {
	opts := defaultCmpOpts()

	for repoDesc, componentUnderTest := range implementationsToTest() {
		// Call setup function, inject t *testing.T, and use t.Cleanup
		repoImpl := componentUnderTest(t)
		if repoImpl == nil {
			t.Fatalf(
				"%s: no implementation provided",
				repoDesc,
			)
		}

		res, err := repoImpl.Read(context.Background(), nil)
		if err != nil {
			t.Fatalf(
				"%s: initial condition: empty repo: %s",
				repoDesc,
				err,
			)
		}
		if len(res) != 0 {
			t.Fatalf(
				"%s: initial condition: empty repo: got %d items",
				repoDesc,
				len(res),
			)
		}
		expected := []*candidateKey.SAEnum{
			{
				Id:   candidateKey.CandidateKeyEnum_ZERO,
				Data: candidateKey.CandidateKeyEnum_ZERO.String(),
			},
			{
				Id:   candidateKey.CandidateKeyEnum_ONE,
				Data: candidateKey.CandidateKeyEnum_ONE.String(),
			},
			{
				Id:   candidateKey.CandidateKeyEnum_TWO,
				Data: candidateKey.CandidateKeyEnum_TWO.String(),
			},
			{
				Id:   candidateKey.CandidateKeyEnum_THREE,
				Data: candidateKey.CandidateKeyEnum_THREE.String(),
			},
		}
		res, err = repoImpl.Create(context.Background(), expected)
		if err != nil {
			t.Fatalf(
				"%s: Create(): %s",
				repoDesc,
				err,
			)
		}

		if diff := cmp.Diff(expected, res, opts); diff != "" {
			t.Fatal(
				mismatch(
					fmt.Sprintf(
						"%s: Create():",
						repoDesc,
					),
					diff,
				),
			)
		}
	}
}

func implementationsToTest() map[string]saEnumComponentUnderTest {
	return map[string]saEnumComponentUnderTest{
		"SQLite": sqliteSAComponentUnderTest,
	}
}

func defaultCmpOpts() cmp.Options {
	return cmp.Options{
		cmpopts.IgnoreUnexported(candidateKey.SAEnum{}),
		cmpopts.SortSlices(func(x, y *candidateKey.SAEnum) bool {
			switch strings.Compare(x.GetId().String(), y.GetId().String()) {
			case -1:
				return true
			case 0:
				return true
			case 1:
				return false
			}
			panic("this should never happen")
		}),
	}
}

func mismatch(prefix, diff string) string {
	return fmt.Sprintf(
		"%s mismatch (-want +got):\n%s",
		prefix,
		diff,
	)
}
