package grammars

import (
	"database/contracts"
	"database/query"
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
	mg.Grammar.SetSelectComponents(mg.Grammar.GetDefaultSelectComponents())

	return mg
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
