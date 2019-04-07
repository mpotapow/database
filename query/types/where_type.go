package types

import "database/contracts"

type WhereType interface {
	GetLogic() string
	GetColumn() string
	GetOperator() string
	ValueToArray() []interface{}
}

type WhereQuery interface {
	GetQuery() contracts.QueryBuilder
}

type WhereDateType interface {
	GetDateType() string
}

type WhereCallback = func(q contracts.QueryBuilder)

type Where struct {
	logic string
	column string
	operator string
}

func (w *Where) GetColumn() string {
	return w.column
}

func (w *Where) GetOperator() string {
	return w.operator
}

func (w *Where) GetLogic() string {
	return w.logic
}

func newWhere(col string, operator string, logic string) *Where {
	return &Where{
		column: col,
		logic: logic,
		operator: operator,
	}
}

type WhereColumn struct {
	*Where
	value string
}

func (w *WhereColumn) ValueToArray() []interface{} {
	return []interface{}{w.value}
}

func (w *WhereColumn) ValueToString() string {
	return w.value
}

func NewWhereColumn(col string, operator string, value string, logic string) *WhereColumn {
	return &WhereColumn{
		value: value,
		Where: newWhere(col, operator, logic),
	}
}

type WhereRaw struct {
	*Where
	value string
}

func (w *WhereRaw) ValueToArray() []interface{} {
	return []interface{}{w.value}
}

func (w *WhereRaw) ValueToString() string {
	return w.value
}

func NewWhereRaw(condition string, logic string) *WhereRaw {
	return &WhereRaw{
		value: condition,
		Where: newWhere("", "", logic),
	}
}

type WhereInt struct {
	*Where
	value int
}

func (w *WhereInt) ValueToArray() []interface{} {
	return []interface{}{w.value}
}

func NewWhereInt(col string, operator string, value int, logic string) *WhereInt {
	return &WhereInt{
		value: value,
		Where: newWhere(col, operator, logic),
	}
}

type WhereString struct {
	*Where
	value string
}

func (w *WhereString) ValueToArray() []interface{} {
	return []interface{}{w.value}
}

func NewWhereString(col string, operator string, value string, logic string) *WhereString {
	return &WhereString{
		value: value,
		Where: newWhere(col, operator, logic),
	}
}

type WhereFloat32 struct {
	*Where
	value float32
}

func (w *WhereFloat32) ValueToArray() []interface{} {
	return []interface{}{w.value}
}

func NewWhereFloat32(col string, operator string, value float32, logic string) *WhereFloat32 {
	return &WhereFloat32{
		value: value,
		Where: newWhere(col, operator, logic),
	}
}

type WhereBool struct {
	*Where
	value bool
}

func (w *WhereBool) ValueToArray() []interface{} {
	return []interface{}{w.value}
}

func NewWhereBool(col string, operator string, value bool, logic string) *WhereBool {
	return &WhereBool{
		value: value,
		Where: newWhere(col, operator, logic),
	}
}

type WhereNull struct {
	*Where
}

func (w *WhereNull) ValueToString() string {
	return "null"
}

func (w *WhereNull) ValueToArray() []interface{} {
	return []interface{}{w.ValueToString()}
}

func NewWhereNull(col string, operator string, logic string) *WhereNull {
	switch operator {
		case "=":
			operator = "is"
			break
		case "!=":
			operator = "is not"
			break
		default:
			panic("Illegal operator and value combination.")
	}
	return &WhereNull{
		Where: newWhere(col, operator, logic),
	}
}

type WhereIn struct {
	*Where
	value []interface{}
}

func (w *WhereIn) ValueToArray() []interface{} {
	return w.value
}

func NewWhereIn(col string, operator string, value []interface{}, logic string) *WhereIn {

	return &WhereIn{
		value: value,
		Where: newWhere(col, operator, logic),
	}
}

type WhereBetween struct {
	*Where
	value []interface{}
}

func (w *WhereBetween) ValueToArray() []interface{} {
	return w.value
}

func NewWhereBetween(col string, operator string, value []interface{}, logic string) *WhereBetween {

	return &WhereBetween{
		value: value,
		Where: newWhere(col, operator, logic),
	}
}

type WhereDate struct {
	*Where
	value string
	dateType string
}

func (w *WhereDate) ValueToArray() []interface{} {
	return []interface{}{w.value}
}

func (w *WhereDate) GetDateType() string {
	return w.dateType
}

func NewWhereDate(col string, operator string, value string, dateType string, logic string) *WhereDate {

	return &WhereDate{
		value: value,
		dateType: dateType,
		Where: newWhere(col, operator, logic),
	}
}

type WhereNested struct {
	*Where
	value contracts.QueryBuilder
}

func (w *WhereNested) GetQuery() contracts.QueryBuilder {
	return w.value
}

func (w *WhereNested) ValueToArray() []interface{} {
	return []interface{}{w.value}
}

func NewWhereNested(value contracts.QueryBuilder, logic string) *WhereNested {

	return &WhereNested{
		value: value,
		Where: newWhere("", "", logic),
	}
}

type WhereSub struct {
	*Where
	value contracts.QueryBuilder
}

func (w *WhereSub) GetQuery() contracts.QueryBuilder {
	return w.value
}

func (w *WhereSub) ValueToArray() []interface{} {
	return []interface{}{w.value}
}

func NewWhereSub(column string, operator string, value contracts.QueryBuilder, logic string) *WhereSub {

	return &WhereSub{
		value: value,
		Where: newWhere(column, operator, logic),
	}
}