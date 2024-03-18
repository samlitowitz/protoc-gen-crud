//go:generate stringer -type=Typ
package attribute

type Typ int

const (
	UNKNOWN Typ = iota

	BOOL    Typ = iota
	BYTE    Typ = iota
	INTEGER Typ = iota
	FLOAT   Typ = iota
	STRING  Typ = iota
)
