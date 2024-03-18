//go:generate stringer -type=Direction
package relationship

type Direction int

const (
	BIDIRECTIONAL  Direction = iota
	UNIDIRECTIONAL Direction = iota
)
