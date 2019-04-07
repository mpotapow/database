package query

import (
	"database/contracts"
	"database/query/types"
	"database/sql"
	"time"
)

type Builder struct {

	Table string

	RowLimit int
	RowOffset int

	Aggregate types.AggregateType
	Orders []types.OrderType
	Wheres []types.WhereType
	Groups []string
	Havings []string
	Columns []types.SelectType

	bindings map[string][]interface{}

	grammar contracts.Grammar
	connection contracts.Connection
}

func NewBuilder(connection contracts.Connection, grammar contracts.Grammar) contracts.QueryBuilder {

	return &Builder{
		grammar: grammar,
		connection: connection,
		bindings: map[string][]interface{}{
			"select": make([]interface{}, 0),
			"from": make([]interface{}, 0),
			"join": make([]interface{}, 0),
			"where": make([]interface{}, 0),
			"having": make([]interface{}, 0),
			"order": make([]interface{}, 0),
			"union": make([]interface{}, 0),
		},
	}
}

func (b *Builder) Select(args ...string) contracts.QueryBuilder {

	for _, arg := range args {
		b.Columns = append(b.Columns, types.NewSelectString(arg))
	}

	return b
}

func (b *Builder) SelectRaw(args ...string) contracts.QueryBuilder {

	for _, arg := range args {
		b.Columns = append(b.Columns, types.NewSelectRawString(arg))
	}

	return b
}

func (b *Builder) From(from string) contracts.QueryBuilder {

	b.Table = from

	return b
}

func (b *Builder) buildWhere(logic string, args ...interface{}) contracts.QueryBuilder {

	col, value, operator := b.prepareWhereArguments(args...)

	if value == nil {
		if operator == "=" {
			return b.WhereNull(col)
		} else {
			return b.WhereNotNull(col)
		}
	}

	whereType := b.getWhereTypeByValue(col, operator, value, logic)

	b.Wheres = append(b.Wheres, whereType)
	b.addBinding(value, "where")

	return b
}

func (b *Builder) buildWhereColumn(logic string, args ...interface{}) contracts.QueryBuilder {

	col, col2, operator := b.prepareWhereArguments(args...)
	value := col2.(string)

	whereType := types.NewWhereColumn(col, operator, value, logic)
	b.Wheres = append(b.Wheres, whereType)

	return b
}

func (b *Builder) buildWhereRaw(condition string, bindings []interface{}, logic string) contracts.QueryBuilder {

	whereType := types.NewWhereRaw(condition, logic)
	b.Wheres = append(b.Wheres, whereType)

	for _, v := range bindings {
		b.addBinding(v, "where")
	}

	return b
}

func (b *Builder) buildWhereNull(col string, operator string, logic string) contracts.QueryBuilder {

	whereType := b.getWhereTypeByValue(col, operator, nil, logic)
	b.Wheres = append(b.Wheres, whereType)

	return b
}

func (b *Builder) buildWhereIn(column string, operator string, values []interface{}, logic string) contracts.QueryBuilder {

	whereType := types.NewWhereIn(column, operator, values, logic)
	b.Wheres = append(b.Wheres, whereType)

	for _, v := range values {
		b.addBinding(v, "where")
	}

	return b
}

func (b *Builder) buildWhereBetween(column string, operator string, values []interface{}, logic string) contracts.QueryBuilder {

	whereType := types.NewWhereBetween(column, operator, values, logic)
	b.Wheres = append(b.Wheres, whereType)

	for _, v := range values {
		b.addBinding(v, "where")
	}

	return b
}

func (b *Builder) buildWhereDate(dateType string, format string, logic string, args ...interface{}) contracts.QueryBuilder {

	col, value, operator := b.prepareWhereArguments(args...)

	switch v := value.(type) {
		case string:
			value = v
			break
		case time.Time:
			value = v.Format(format)
			break
	}

	whereType := types.NewWhereDate(col, operator, value.(string), dateType, logic)
	b.Wheres = append(b.Wheres, whereType)

	b.addBinding(value, "where")

	return b
}

func (b *Builder) Where(args ...interface{}) contracts.QueryBuilder {

	return b.buildWhere("and", args...)
}

func (b *Builder) OrWhere(args ...interface{}) contracts.QueryBuilder {

	return b.buildWhere("or", args...)
}

func (b *Builder) WhereColumn(args ...interface{}) contracts.QueryBuilder {

	return b.buildWhereColumn("and", args...)
}

func (b *Builder) OrWhereColumn(args ...interface{}) contracts.QueryBuilder {

	return b.buildWhereColumn("or", args...)
}

func (b *Builder) WhereRaw(condition string, bindings ...interface{}) contracts.QueryBuilder {

	return b.buildWhereRaw(condition, bindings, "and")
}

func (b *Builder) OrWhereRaw(condition string, bindings ...interface{}) contracts.QueryBuilder {

	return b.buildWhereRaw(condition, bindings, "or")
}

func (b *Builder) WhereNull(col string) contracts.QueryBuilder {

	return b.buildWhereNull(col, "=", "and")
}

func (b *Builder) OrWhereNull(col string) contracts.QueryBuilder {

	return b.buildWhereNull(col, "=", "or")
}

func (b *Builder) WhereNotNull(col string) contracts.QueryBuilder {

	return b.buildWhereNull(col, "!=", "and")
}

func (b *Builder) OrWhereNotNull(col string) contracts.QueryBuilder {

	return b.buildWhereNull(col, "!=", "or")
}

func (b *Builder) WhereIn(column string, values []interface{}) contracts.QueryBuilder {

	return b.buildWhereIn(column, "in", values, "and")
}

func (b *Builder) OrWhereIn(column string, values []interface{}) contracts.QueryBuilder {

	return b.buildWhereIn(column, "in", values, "or")
}

func (b *Builder) WhereNotIn(column string, values []interface{}) contracts.QueryBuilder {

	return b.buildWhereIn(column, "not in", values, "and")
}

func (b *Builder) OrWhereNotIn(column string, values []interface{}) contracts.QueryBuilder {

	return b.buildWhereIn(column, "not in", values, "or")
}

func (b *Builder) WhereBetween(column string, from interface{}, to interface{}) contracts.QueryBuilder {

	return b.buildWhereBetween(column, "between", []interface{}{from, to}, "and")
}

func (b *Builder) OrWhereBetween(column string, from interface{}, to interface{}) contracts.QueryBuilder {

	return b.buildWhereBetween(column, "between", []interface{}{from, to}, "or")
}

func (b *Builder) WhereNotBetween(column string, from interface{}, to interface{}) contracts.QueryBuilder {

	return b.buildWhereBetween(column, "not between", []interface{}{from, to}, "and")
}

func (b *Builder) OrWhereNotBetween(column string, from interface{}, to interface{}) contracts.QueryBuilder {

	return b.buildWhereBetween(column, "not between", []interface{}{from, to}, "or")
}

func (b *Builder) WhereDate(args ...interface{}) contracts.QueryBuilder {

	return b.buildWhereDate("date",  "2006-01-02", "and", args...)
}

func (b *Builder) OrWhereDate(args ...interface{}) contracts.QueryBuilder {

	return b.buildWhereDate("date", "2006-01-02", "or", args...)
}

func (b *Builder) WhereTime(args ...interface{}) contracts.QueryBuilder {

	return b.buildWhereDate("time", "15:04:05", "and", args...)
}

func (b *Builder) OrWhereTime(args ...interface{}) contracts.QueryBuilder {

	return b.buildWhereDate("time", "15:04:05", "or", args...)
}

func (b *Builder) WhereDay(args ...interface{}) contracts.QueryBuilder {

	return b.buildWhereDate("day", "02", "and", args...)
}

func (b *Builder) OrWhereDay(args ...interface{}) contracts.QueryBuilder {

	return b.buildWhereDate("day", "02", "or", args...)
}

func (b *Builder) WhereMonth(args ...interface{}) contracts.QueryBuilder {

	return b.buildWhereDate("month", "01", "and", args...)
}

func (b *Builder) OrWhereMonth(args ...interface{}) contracts.QueryBuilder {

	return b.buildWhereDate("month", "01", "or", args...)
}

func (b *Builder) WhereYear(args ...interface{}) contracts.QueryBuilder {

	return b.buildWhereDate("year", "2006", "and", args...)
}

func (b *Builder) OrWhereYear(args ...interface{}) contracts.QueryBuilder {

	return b.buildWhereDate("year", "2006", "or", args...)
}

func (b *Builder) GroupBy(args ...string) contracts.QueryBuilder {

	b.Groups = append(b.Groups, args...)

	return b
}

func (b *Builder) OrderBy(column string, direction string) contracts.QueryBuilder {

	b.Orders = append(b.Orders, types.NewOrder(column, direction))

	return b
}

func (b *Builder) Limit(n int) contracts.QueryBuilder {

	b.RowLimit = n

	return b
}

func (b *Builder) Offset(n int) contracts.QueryBuilder {

	b.RowOffset = n

	return b
}

func (b *Builder) Get() (*sql.Rows, error) {

	return b.runSelect()
}

func (b *Builder) Count(column string) (*sql.Rows, error) {

	return b.aggregate("count", column)
}

func (b *Builder) Min(column string) (*sql.Rows, error) {

	return b.aggregate("min", column)
}

func (b *Builder) Max(column string) (*sql.Rows, error) {

	return b.aggregate("max", column)
}

func (b *Builder) Sum(column string) (*sql.Rows, error) {

	return b.aggregate("sum", column)
}

func (b *Builder) Avg(column string) (*sql.Rows, error) {

	return b.aggregate("avg", column)
}

func (b *Builder) aggregate(function string, column string) (*sql.Rows, error) {

	var clone = *b
	clone.Columns = []types.SelectType{types.NewSelectString(column)}
	clone.setAggregate(function, column)

	return clone.Get()
}

func (b *Builder) setAggregate(function string, column string) contracts.QueryBuilder {

	b.Aggregate = types.NewAggregate(function, column)

	return b
}

func (b *Builder) runSelect() (*sql.Rows, error) {

	return b.connection.Select(b.ToSql(), b.getBindingsForSql())
}

func (b *Builder) ToSql() string {

	return b.grammar.CompileSelect(b)
}

func (b *Builder) prepareWhereArguments(args ...interface{}) (string, interface{}, string) {

	if len(args) == 2 {
		return args[0].(string), args[1], "="
	} else {
		return args[0].(string), args[2], args[1].(string)
	}
}

func (b *Builder) getWhereTypeByValue(col string, operator string, value interface{}, logic string) types.WhereType {
	switch v := value.(type) {
		case string:
			return types.NewWhereString(col, operator, v, logic)
		case int:
			return types.NewWhereInt(col, operator, v, logic)
		case float32:
			return types.NewWhereFloat32(col, operator, v, logic)
		case bool:
			return types.NewWhereBool(col, operator, v, logic)
		case nil:
			return types.NewWhereNull(col, operator, logic)
		default:
			panic("Illegal where type")
	}
}

func (b *Builder) addBinding(value interface{}, bindingType string) {

	if value != nil {
		b.bindings[bindingType] = append(b.bindings[bindingType], value)
	}
}

func (b *Builder) getBindingsForSql() []interface{} {

	res := make([]interface{}, 0)
	bindingIterator := []string{"select", "from", "join", "where", "having", "order", "union"}

	for _, t := range bindingIterator {
		for _, v := range b.bindings[t] {
			res = append(res, v)
		}
	}

	return res
}