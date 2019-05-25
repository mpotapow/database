package grammars

import (
	"database/contracts"
	"database/query"
	"database/query/types"
	"strings"
)

type MysqlGrammar struct {
	*Grammar
}

func NewMysqlGrammar() *MysqlGrammar {

	var mg = &MysqlGrammar{
		Grammar: NewGrammar(),
	}

	mg.Grammar.SetParametrizeSymbol("?")
	mg.Grammar.SetSelectComponents(mg.GetMysqlSelectComponents())

	return mg
}

func (g *MysqlGrammar) GetMysqlSelectComponents() map[int]interface{} {

	return map[int]interface{}{
		0:  g.compileAggregate,
		1:  g.compileColumns,
		2:  g.compileFrom,
		3:  g.compileJoins,
		4:  g.compileWhere,
		5:  g.compileGroups,
		6:  g.compileHavings,
		7:  g.compileOrders,
		8:  g.compileLimit,
		9:  g.compileOffset,
		10: g.compileLock,
	}
}

func (g *MysqlGrammar) CompileSelect(b contracts.QueryBuilder) string {

	sql := g.Grammar.CompileSelect(b)

	queryBuilder := b.(*query.Builder)

	if len(queryBuilder.Unions) > 0 {

		sql = "(" + sql + ") " + g.compileUnions(b, queryBuilder)
	}

	return sql
}

func (g *MysqlGrammar) compileUnions(b contracts.QueryBuilder, queryBuilder *query.Builder) string {

	sql := ""
	for _, v := range queryBuilder.Unions {
		sql += g.compileUnion(v)
	}

	if len(queryBuilder.UnionOrders) > 0 {
		orderTmp := queryBuilder.Orders
		queryBuilder.Orders = queryBuilder.UnionOrders

		sql += " " + g.compileOrders(b, queryBuilder)

		queryBuilder.Orders = orderTmp
	}

	if queryBuilder.UnionLimit > 0 {
		limitTmp := queryBuilder.RowLimit
		queryBuilder.RowLimit = queryBuilder.UnionLimit

		sql += " " + g.compileLimit(b, queryBuilder)

		queryBuilder.RowLimit = limitTmp
	}

	if queryBuilder.UnionOffset > 0 {
		offsetTmp := queryBuilder.RowOffset
		queryBuilder.RowOffset = queryBuilder.UnionOffset

		sql += " " + g.compileOffset(b, queryBuilder)

		queryBuilder.RowOffset = offsetTmp
	}

	return strings.TrimLeft(sql, " ")
}

func (g *MysqlGrammar) compileUnion(union types.UnionType) string {

	conjunction := " union "
	if union.IsAll() {
		conjunction = " union all "
	}

	return conjunction + "(" + union.GetValue().ToSql() + ")"
}

func (g *MysqlGrammar) CompileUpdate(b contracts.QueryBuilder, values map[string]interface{}) string {

	builder := b.(*query.Builder)
	table := g.WrapTable(builder.Table)

	joins := ""
	if len(builder.Joins) > 0 {
		joins = " " + g.compileJoins(b, builder)
	}

	var columns []string
	for col, _ := range values {
		columns = append(columns, g.Wrap(col)+" = "+g.parametrizeSymbol)
	}

	wheres := g.compileWhere(b, builder)

	q := "update " + table + joins + " set " + strings.Join(columns, ", ") + " " + wheres

	if len(builder.Orders) > 0 {
		q += " " + g.compileOrders(b, builder)
	}

	if builder.RowLimit > 0 {
		q += " " + g.compileLimit(b, builder)
	}

	return strings.Trim(q, " ")
}

func (g *MysqlGrammar) CompileDelete(b contracts.QueryBuilder) string {

	builder := b.(*query.Builder)
	table := g.WrapTable(builder.Table)

	wheres := g.compileWhere(b, builder)

	if len(builder.Joins) > 0 {
		return g.compileDeleteWithJoins(b, builder, table, wheres)
	} else {
		return g.compileDeleteWithoutJoins(b, builder, table, wheres)
	}
}

func (g *MysqlGrammar) PrepareBindingsForDelete(b contracts.QueryBuilder, bindings map[string][]interface{}) []interface{} {

	var res []interface{}
	res = append(res, bindings["join"]...)

	var queryBuilder = b.(*query.Builder)
	exceptBindings := queryBuilder.GetBindingsForSql("join", "select")

	res = append(res, exceptBindings...)

	return res
}

func (g *MysqlGrammar) compileDeleteWithJoins(
	b contracts.QueryBuilder, builder *query.Builder, table string, wheres string,
) string {

	joins := " " + g.compileJoins(b, builder)

	return strings.Trim("delete from "+table+joins+" "+wheres, " ")
}

func (g *MysqlGrammar) compileDeleteWithoutJoins(
	b contracts.QueryBuilder, builder *query.Builder, table string, wheres string,
) string {

	q := strings.Trim("delete from "+table+" "+wheres, " ")

	if len(builder.Orders) > 0 {
		q += " " + g.compileOrders(b, builder)
	}

	if builder.RowLimit > 0 {
		q += " " + g.compileLimit(b, builder)
	}

	return q
}
