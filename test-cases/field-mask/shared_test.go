package field_mask_test

import "fmt"

func mismatch(prefix, diff string) string {
	return fmt.Sprintf(
		"%s mismatch (-want +got):\n%s",
		prefix,
		diff,
	)
}
