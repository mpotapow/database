package grammars

import (
	"database/contracts"
	"database/query"
	"database/query/types"
	"fmt"
	"strings"
)

type Grammar struct {
	parametrizeSymbol string
	selectComponents  map[int]interface{}
}

func NewGrammar() *Grammar {

	return &Grammar{
		parametrizeSymbol: "?",
		selectComponents:  map[int]interface{}{},
	}
}

func (g *Grammar) GetDefaultSelectComponents() map[int]interface{} {

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
		10: g.compileUnions,
		11: g.compileLock,
	}
}

func (g *Grammar) SetSelectComponents(m map[int]interface{}) {

	g.selectComponents = m
}

func (g *Grammar) SetParametrizeSymbol(s string) {

	g.parametrizeSymbol = s
}

func (g *Grammar) CompileSelect(b contracts.QueryBuilder) string {

	var queryBuilder = b.(*query.Builder)

	if len(queryBuilder.Columns) <= 0 {

		queryBuilder.Select("*")
	}

	return g.compileComponents(b, queryBuilder)
}

func (g *Grammar) compileComponents(b contracts.QueryBuilder, queryBuilder *query.Builder) string {

	var res = make([]string, 0)
	var l = len(g.selectComponents)

	for i := 0; i < l; i++ {
		var f = g.selectComponents[i]
		var part = f.(func(contracts.QueryBuilder, *query.Builder) string)(b, queryBuilder)
		res = append(res, part)
	}

	res = g.filter(res, func(s string) bool {
		return len(s) > 0
	})

	return strings.Trim(strings.Join(res, " "), " ")
}

func (g *Grammar) compileAggregate(b contracts.QueryBuilder, queryBuilder *query.Builder) string {

	if queryBuilder.Aggregate == nil {
		return ""
	}

	return fmt.Sprintf("select %s(%s) as aggregate",
		queryBuilder.Aggregate.GetFunction(),
		g.Wrap(queryBuilder.Aggregate.GetColumns()),
	)
}

func (g *Grammar) compileColumns(b contracts.QueryBuilder, queryBuilder *query.Builder) string {

	if queryBuilder.Aggregate != nil {
		return ""
	}

	res := make([]string, 0)
	for _, s := range queryBuilder.Columns {
		switch v := s.(type) {
		case *types.SelectString:
			res = append(res, g.Wrap(v.ToString()))
			break
		case *types.SelectRawString:
			res = append(res, v.ToString())
			break
		}
	}

	return "select " + strings.Join(res, ", ")
}

func (g *Grammar) compileFrom(b contracts.QueryBuilder, queryBuilder *query.Builder) string {

	return "from " + g.WrapTable(queryBuilder.Table)
}

func (g *Grammar) compileJoins(b contracts.QueryBuilder, queryBuilder *query.Builder) string {

	var res []string
	for _, j := range queryBuilder.Joins {

		qb := j.(*query.JoinClause)

		table := g.WrapTable(qb.Table)

		nestedJoins := ""
		if len(qb.Joins) != 0 {

			nestedJoins = " " + g.compileJoins(b, queryBuilder)
		}

		q := j.GetType() + " join " + table + nestedJoins + " " + g.compileWhere(qb, qb.Builder)
		res = append(res, strings.Trim(q, " "))
	}

	return strings.Join(res, " ")
}

func (g *Grammar) compileWhere(b contracts.QueryBuilder, queryBuilder *query.Builder) string {

	if len(queryBuilder.Wheres) <= 0 {
		return ""
	}

	res := make([]string, 0)
	for _, w := range queryBuilder.Wheres {
		switch w.(type) {
		default:
			condition := fmt.Sprintf("%v %v %v", g.Wrap(w.GetColumn()), w.GetOperator(), g.parameterizeWhere(w, ", "))
			res = append(res, w.GetLogic()+" "+condition)
			break
		case *types.WhereIn:
			value := "(" + g.parameterizeWhere(w, ", ") + ")"
			condition := fmt.Sprintf("%v %v %v", g.Wrap(w.GetColumn()), w.GetOperator(), value)
			res = append(res, w.GetLogic()+" "+condition)
			break
		case *types.WhereColumn:
			where := w.(types.ExpressionType)
			condition := fmt.Sprintf("%v %v %v", g.Wrap(w.GetColumn()), w.GetOperator(), g.Wrap(where.ValueToString()))
			res = append(res, w.GetLogic()+" "+condition)
			break
		case *types.WhereRaw:
			where := w.(types.ExpressionType)
			res = append(res, w.GetLogic()+" "+where.ValueToString())
			break
		case *types.WhereBetween:
			value := g.parameterizeWhere(w, " and ")
			condition := fmt.Sprintf("%v %v %v", g.Wrap(w.GetColumn()), w.GetOperator(), value)
			res = append(res, w.GetLogic()+" "+condition)
			break
		case *types.WhereDate:
			where := w.(types.WhereDateType)
			condition := fmt.Sprintf("%v(%v) %v %v",
				where.GetDateType(),
				g.Wrap(w.GetColumn()),
				w.GetOperator(),
				g.parameterizeWhere(w, ", "),
			)
			res = append(res, w.GetLogic()+" "+condition)
			break
		case *types.WhereNested:
			builder := g.getQueryByWhere(w)
			str := g.compileWhere(b, builder)
			res = append(res, w.GetLogic()+" ("+str[6:]+")")
			break
		case *types.WhereSub:
			builder := g.getQueryByWhere(w)
			selectRaw := g.CompileSelect(builder)
			res = append(res, g.Wrap(w.GetColumn())+" "+w.GetOperator()+" ("+selectRaw+")")
			break
		}
	}

	conjunction := "where"

	switch b.(type) {
	case *query.JoinClause:
		conjunction = "on"
	}

	return conjunction + " " + g.removeLeadingBoolean(strings.Join(res, " "))
}

func (g *Grammar) compileGroups(b contracts.QueryBuilder, queryBuilder *query.Builder) string {

	if len(queryBuilder.Groups) <= 0 {
		return ""
	}

	return "group by " + g.columnize(queryBuilder.Groups)
}

func (g *Grammar) compileHavings(b contracts.QueryBuilder, queryBuilder *query.Builder) string {

	if len(queryBuilder.Havings) <= 0 {
		return ""
	}

	res := make([]string, 0)
	for _, w := range queryBuilder.Havings {
		switch w.(type) {
		default:
			condition := fmt.Sprintf("%v %v %v", g.Wrap(w.GetColumn()), w.GetOperator(), g.parameterizeWhere(w, ", "))
			res = append(res, w.GetLogic()+" "+condition)
			break
		case *types.WhereRaw:
			where := w.(types.ExpressionType)
			res = append(res, w.GetLogic()+" "+where.ValueToString())
			break
		}
	}

	return "having " + g.removeLeadingBoolean(strings.Join(res, " "))
}

func (g *Grammar) compileOrders(b contracts.QueryBuilder, queryBuilder *query.Builder) string {

	if len(queryBuilder.Orders) <= 0 {
		return ""
	}

	orders := make([]string, 0)
	for _, o := range queryBuilder.Orders {
		if g.isExpression(o) {
			orderExpr := o.(types.ExpressionType)
			orders = append(orders, orderExpr.ValueToString())
		} else {
			orders = append(orders, g.Wrap(o.GetColumn())+" "+o.GetDirection())
		}
	}

	return "order by " + strings.Join(orders, ", ")
}

func (g *Grammar) compileLimit(b contracts.QueryBuilder, queryBuilder *query.Builder) string {

	if queryBuilder.RowLimit <= 0 {
		return ""
	}

	return fmt.Sprintf("limit %v", queryBuilder.RowLimit)
}

func (g *Grammar) compileOffset(b contracts.QueryBuilder, queryBuilder *query.Builder) string {

	if queryBuilder.RowOffset <= 0 {
		return ""
	}

	return fmt.Sprintf("offset %v", queryBuilder.RowOffset)
}

func (g *Grammar) compileUnions(b contracts.QueryBuilder, queryBuilder *query.Builder) string {

	return ""
}

func (g *Grammar) compileLock(b contracts.QueryBuilder, queryBuilder *query.Builder) string {

	return ""
}

func (g *Grammar) CompileInsert(
	b contracts.QueryBuilder, values []map[string]interface{}, columns []string,
) string {

	builder := b.(*query.Builder)
	table := g.WrapTable(builder.Table)
	mock := make([]interface{}, len(values[0:1][0]))

	var params []string
	for _, _ = range values {
		params = append(params, "("+g.parameterize(mock, ", ")+")")
	}

	return "insert into " + table + "(" + g.columnize(columns) + ") values " + strings.Join(params, ", ")
}

func (g *Grammar) CompileUpdate(b contracts.QueryBuilder, values map[string]interface{}) string {

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

	return strings.Trim(q, " ")
}

func (g *Grammar) CompileDelete(b contracts.QueryBuilder) string {

	builder := b.(*query.Builder)
	table := g.WrapTable(builder.Table)

	wheres := g.compileWhere(b, builder)

	return strings.Trim("delete from " + table + " " + wheres, " ")
}

func (g *Grammar) CompileSavepoint(name string) string {

	return "SAVEPOINT " + name
}

func (g *Grammar) CompileSavepointRollback(name string) string {

	return "ROLLBACK TO SAVEPOINT " + name
}

func (g *Grammar) columnize(columns []string) string {

	for i, v := range columns {
		columns[i] = g.Wrap(v)
	}

	return strings.Join(columns, ", ")
}

func (g *Grammar) parameterize(values []interface{}, sep string) string {

	var res []string
	for _, _ = range values {
		res = append(res, g.parametrizeSymbol)
	}

	return strings.Join(res, sep)
}

func (g *Grammar) parameterizeWhere(where types.WhereType, sep string) string {

	var res []string
	if g.isExpression(where) {
		var exprArg = where.(types.ExpressionType)
		res = append(res, exprArg.ValueToString())
	} else {
		for _, _ = range where.ValueToArray() {
			res = append(res, g.parametrizeSymbol)
		}
	}

	return strings.Join(res, sep)
}

func (g *Grammar) isExpression(val interface{}) bool {

	_, ok := val.(types.ExpressionType)
	return ok
}

func (g *Grammar) getQueryByWhere(w types.WhereType) *query.Builder {

	where := w.(types.WhereQuery)
	queryBuilder := where.GetQuery()

	return queryBuilder.(*query.Builder)
}

func (g *Grammar) filter(iterator []string, f func(string) bool) []string {

	res := make([]string, 0)
	for _, v := range iterator {
		if f(v) {
			res = append(res, v)
		}
	}

	return res
}

func (g *Grammar) removeLeadingBoolean(query string) string {

	v := strings.Index(query, "and")
	if v == 0 {
		return query[4:]
	}

	v = strings.Index(query, "or")
	if v == 0 {
		return query[3:]
	}

	return query
}

func (g *Grammar) Wrap(v string) string {

	if strings.Index(v, ".") > -1 {
		var res []string
		for _, v := range strings.Split(v, ".") {
			res = append(res, g.Wrap(v))
		}

		return strings.Join(res, ".")
	}

	if v == "*" {
		return v
	}

	return "`" + v + "`"
}

func (g *Grammar) WrapTable(table interface{}) string {

	switch v := table.(type) {
	case *types.FromString:
		return g.Wrap(v.ToString())
	case *types.FromRawString:
		return v.ToString()
	}

	panic("Wrong table params for wrap")
}

func (g *Grammar) PrepareBindingsForUpdate(
	b contracts.QueryBuilder, bindings map[string][]interface{}, values map[string]interface{},
) []interface{} {

	var res []interface{}
	res = append(res, bindings["join"]...)

	for _, v := range values {
		res = append(res, v)
	}

	var queryBuilder = b.(*query.Builder)
	exceptBindings := queryBuilder.GetBindingsForSql("join", "select")

	res = append(res, exceptBindings...)

	return res
}

func (g *Grammar) PrepareBindingsForDelete(b contracts.QueryBuilder, bindings map[string][]interface{}) []interface{} {

	return b.(*query.Builder).GetBindingsForSql()
}