//go:generate stringer -type=Typ
package relationship

type Typ int

const (
	MANY_TO_MANY Typ = iota
	MANY_TO_ONE  Typ = iota
	ONE_TO_MANY  Typ = iota
	ONE_TO_ONE   Typ = iota
)
