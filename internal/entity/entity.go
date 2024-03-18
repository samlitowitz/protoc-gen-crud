package entity

type Entity struct {
	Attributes    []*Attribute
	Relationships []*Relationship

	MetaData map[string]any
}
