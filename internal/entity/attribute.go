package entity

import (
	"fmt"

	"github.com/samlitowitz/protoc-gen-crud/internal/entity/attribute"
)

type Attribute struct {
	Type attribute.Typ

	unsigned bool
}

func (att *Attribute) IsUnsigned() bool {
	return att.unsigned
}

func (att *Attribute) SetSigned() {
	att.unsigned = false
}

func (att *Attribute) SetUnsigned() error {
	switch att.Type {
	case attribute.INTEGER:
		fallthrough
	case attribute.FLOAT:
		att.unsigned = true
		return nil
	default:
		return fmt.Errorf(
			"type %s does not support unsigned",
			att.Type,
		)
	}
}
