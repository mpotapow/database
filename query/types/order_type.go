package types

import "strings"

type OrderType interface {
	GetColumn() string
	GetDirection() string
}

type Order struct {
	column string
	direction string
}

func (o *Order) GetColumn() string {
	return o.column
}

func (o *Order) GetDirection() string {
	return o.direction
}

func NewOrder(column string, direction string) *Order {

	if strings.ToLower(direction) == "asc" {
		direction = "asc"
	} else {
		direction = "desc"
	}

	return &Order{
		column: column,
		direction: direction,
	}
}

type OrderRaw struct {
	*Order
	sql string
}

func (o *OrderRaw) ValueToString() string {
	return o.sql
}

func NewOrderRaw(sql string) *OrderRaw {

	return &OrderRaw{
		sql: sql,
		Order: NewOrder("", ""),
	}
}