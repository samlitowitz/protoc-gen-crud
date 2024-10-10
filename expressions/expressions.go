package expressions

import "time"

type FieldID string

type Expression interface {
	Operands() []Expression
}

// Logical
type logical struct {
	left  Expression
	right Expression
}

func (expr logical) Operands() []Expression {
	return []Expression{expr.left, expr.right}
}

func (expr logical) Left() Expression {
	return expr.left
}

func (expr logical) Right() Expression {
	return expr.right
}

type And struct {
	*logical
}

func NewAnd(left, right Expression) *And {
	return &And{
		&logical{
			left:  left,
			right: right,
		},
	}
}

type Or struct {
	*logical
}

func NewOr(left, right Expression) *Or {
	return &Or{
		&logical{
			left:  left,
			right: right,
		},
	}
}

type Not struct {
	expr Expression
}

func NewNot(expr Expression) *Not {
	return &Not{
		expr: expr,
	}
}

func (expr Not) Operands() []Expression {
	return []Expression{expr.expr}
}

func (expr Not) Operand() Expression {
	return expr.expr
}

// Equality

type Equal struct {
	*logical
}

func NewEqual(left, right Expression) *Equal {
	return &Equal{
		&logical{
			left:  left,
			right: right,
		},
	}
}

// Identifier
type Identifier struct {
	id FieldID
}

func NewIdentifier(id FieldID) *Identifier {
	return &Identifier{
		id: id,
	}
}

func (expr Identifier) Operands() []Expression {
	return nil
}

func (expr Identifier) ID() FieldID {
	return expr.id
}

// Scalar
type Scalar struct {
	value any
}

func NewScalar(value any) *Scalar {
	return &Scalar{
		value: value,
	}
}

func (s Scalar) Operands() []Expression {
	return nil
}

func (s Scalar) Value() any {
	return s.value
}

// Timestamp
type Timestamp time.Time

func NewTimestamp(value time.Time) Timestamp {
	return Timestamp(value)
}

func (t Timestamp) Operands() []Expression {
	return nil
}
