package types

import "database/contracts"

type UnionType interface {
	IsAll() bool
	GetValue() contracts.QueryBuilder
}

type Union struct {
	all bool
	value contracts.QueryBuilder
}

func (u *Union) IsAll() bool {
	return u.all
}

func (u *Union) GetValue() contracts.QueryBuilder {
	return u.value
}

func NewUnion(value contracts.QueryBuilder, all bool) *Union {

	return &Union{
		all: all,
		value: value,
	}
}