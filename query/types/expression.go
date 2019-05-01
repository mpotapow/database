package types

type ExpressionType interface {
	ValueToString() string
}

type Expression struct {
	value string
}

func NewExpression(value string) *Expression {
	return &Expression{
		value:value,
	}
}

func (e *Expression) ValueToString() string {

	return e.value
}