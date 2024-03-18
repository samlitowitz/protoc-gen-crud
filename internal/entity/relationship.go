package entity

import "github.com/samlitowitz/protoc-gen-crud/internal/entity/relationship"

type Relationship struct {
	Type      relationship.Typ
	Direction relationship.Direction

	From *Entity
	To   *Entity
}
