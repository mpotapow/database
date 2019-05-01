package query

import (
	"database/contracts"
	"database/query/types"
)

type JoinClause struct {
	contracts.QueryBuilder
	parent *Builder
	joinType string
}

func NewJoinClause(builder *Builder, joinType string, table interface{}) contracts.JoinQueryBuilder {

	JoinClause := &JoinClause{
		builder.newQuery(),
		builder,
		joinType,
	}

	switch v := table.(type) {
		case string:
			JoinClause.From(v)
		case types.ExpressionType:
			qb := JoinClause.QueryBuilder.(*Builder)
			qb.Table = types.NewFromRawString(v.ValueToString())
	}

	return JoinClause
}

func (j *JoinClause) GetType() string {

	return j.joinType
}

func (j *JoinClause) GetQueryBuilder() contracts.QueryBuilder {

	return j.QueryBuilder
}

func (j *JoinClause) buildOn(args []interface{}, where string) contracts.JoinQueryBuilder {

	qb := j.GetQueryBuilder().(*Builder)
	if qb.isCallback(args[0]) {

		qb.whereNested(args[0].(types.WhereCallback), where)
	} else {

		qb.buildWhereColumn(args, where)
	}

	return j
}

func (j *JoinClause) On(args ...interface{}) contracts.JoinQueryBuilder {

	return j.buildOn(args, "and")
}

func (j *JoinClause) OrOn(args ...interface{}) contracts.JoinQueryBuilder {

	return j.buildOn(args, "or")
}

func (j *JoinClause) forSubQuery() contracts.QueryBuilder {

	return j.parent.newQuery()
}

func (j *JoinClause) newQuery() contracts.JoinQueryBuilder {

	qb := j.GetQueryBuilder().(*Builder)
	return NewJoinClause(j.parent, j.joinType, qb.Table.ToString())
}