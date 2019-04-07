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
	selectComponents map[int]interface{}
}

func NewGrammar() *Grammar {

	return &Grammar{
		parametrizeSymbol: "?",
		selectComponents: map[int]interface{}{},
	}
}

func (g *Grammar) GetDefaultSelectComponents() map[int]interface{} {

	return map[int]interface{}{
		0: g.compileAggregate,
		1: g.compileColumns,
		2: g.compileFrom,
		3: g.compileJoins,
		4: g.compileWhere,
		5: g.compileGroups,
		6: g.compileHavings,
		7: g.compileOrders,
		8: g.compileLimit,
		9: g.compileOffset,
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

	return g.compileComponents(queryBuilder)
}

func (g *Grammar) compileComponents(queryBuilder *query.Builder) string {

	var res = make([]string, 0)
	var l = len(g.selectComponents)

	for i := 0; i < l; i++ {
		var f = g.selectComponents[i]
		var part = f.(func(*query.Builder) string)(queryBuilder)
		res = append(res, part)
	}

	res = g.filter(res, func(s string) bool {
		return len(s) > 0
	})

	return strings.Trim(strings.Join(res, " "), " ")
}

func (g *Grammar) compileAggregate(queryBuilder *query.Builder) string {

	if queryBuilder.Aggregate == nil {
		return ""
	}

	return fmt.Sprintf("select %s(%s) as aggregate",
		queryBuilder.Aggregate.GetFunction(),
		g.wrap(queryBuilder.Aggregate.GetColumns()),
	)
}

func (g *Grammar) compileColumns(queryBuilder *query.Builder) string {

	if queryBuilder.Aggregate != nil {
		return ""
	}

	res := make([]string, 0)
	for _, s := range queryBuilder.Columns {
		switch v := s.(type) {
			case *types.SelectString:
				res = append(res, g.wrap(v.ToString()))
				break
			case *types.SelectRawString:
				res = append(res, v.ToString())
				break
		}
	}

	return "select " + strings.Join(res, ", ")
}

func (g *Grammar) compileFrom(queryBuilder *query.Builder) string {

	return "from " + g.wrap(queryBuilder.Table)
}

func (g *Grammar) compileJoins(queryBuilder *query.Builder) string {

	return ""
}

func (g *Grammar) compileWhere(queryBuilder *query.Builder) string {

	if len(queryBuilder.Wheres) <= 0 {
		return ""
	}

	res := make([]string, 0)
	for _, w := range queryBuilder.Wheres {
		switch w.(type) {
			default:
				condition := fmt.Sprintf("%v %v %v", g.wrap(w.GetColumn()), w.GetOperator(), g.parameterize(w, ", "))
				res = append(res, w.GetLogic() + " " + condition)
				break
			case *types.WhereIn:
				value := "(" + g.parameterize(w, ", ") + ")"
				condition := fmt.Sprintf("%v %v %v", g.wrap(w.GetColumn()), w.GetOperator(), value)
				res = append(res, w.GetLogic() + " " + condition)
				break
			case *types.WhereColumn:
				where := w.(types.WhereExpressionType)
				condition := fmt.Sprintf("%v %v %v", g.wrap(w.GetColumn()), w.GetOperator(), g.wrap(where.ValueToString()))
				res = append(res, w.GetLogic() + " " + condition)
				break
			case *types.WhereRaw:
				where := w.(types.WhereExpressionType)
				res = append(res, w.GetLogic() + " " + where.ValueToString())
				break
			case *types.WhereBetween:
				value := g.parameterize(w, " and ")
				condition := fmt.Sprintf("%v %v %v", g.wrap(w.GetColumn()), w.GetOperator(), value)
				res = append(res, w.GetLogic() + " " + condition)
				break
			case *types.WhereDate:
				where := w.(types.WhereDateType)
				condition := fmt.Sprintf("%v(%v) %v %v",
					where.GetDateType(),
					g.wrap(w.GetColumn()),
					w.GetOperator(),
					g.parameterize(w, ", "),
				)
				res = append(res, w.GetLogic() + " " + condition)
				break
		}
	}

	return "where " + g.removeLeadingBoolean(strings.Join(res, " "))
}

func (g *Grammar) compileGroups(queryBuilder *query.Builder) string {

	if len(queryBuilder.Groups) <= 0 {
		return ""
	}

	return "group by " + g.columnize(queryBuilder.Groups)
}

func (g *Grammar) compileHavings(queryBuilder *query.Builder) string {

	return ""
}

func (g *Grammar) compileOrders(queryBuilder *query.Builder) string {

	if len(queryBuilder.Orders) <= 0 {
		return ""
	}

	orders := make([]string, 0)
	for _, o := range queryBuilder.Orders {
		orders = append(orders, g.wrap(o.GetColumn()) + " " + o.GetDirection())
	}

	return "order by " + strings.Join(orders, ", ")
}

func (g *Grammar) compileLimit(queryBuilder *query.Builder) string {

	if queryBuilder.RowLimit <= 0 {
		return ""
	}

	return fmt.Sprintf("limit %v", queryBuilder.RowLimit)
}

func (g *Grammar) compileOffset(queryBuilder *query.Builder) string {

	if queryBuilder.RowOffset <= 0 {
		return ""
	}

	return fmt.Sprintf("offset %v", queryBuilder.RowOffset)
}

func (g *Grammar) compileUnions(queryBuilder *query.Builder) string {

	return ""
}

func (g *Grammar) compileLock(queryBuilder *query.Builder) string {

	return ""
}

func (g *Grammar) columnize(columns []string) string  {

	for i, v := range columns {
		columns[i] = g.wrap(v)
	}

	return strings.Join(columns, ", ")
}

func (g *Grammar) parameterize(where types.WhereType, sep string) string {

	res := make([]string, 0)
	if g.isExpression(where) {
		var exprArg = where.(types.WhereExpressionType)
		res = append(res, exprArg.ValueToString())
	} else {
		for _, _ = range where.ValueToArray() {
			res = append(res, g.parametrizeSymbol)
		}
	}

	return strings.Join(res, sep)
}

func (g *Grammar) isExpression(where types.WhereType) bool {

	 _, ok := where.(types.WhereExpressionType)
	 return ok
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

func (g *Grammar) wrap(v string) string {

	if v == "*" {
		return v
	}

	return "`" + v + "`"
}
