package contracts

import "database/sql"

type QueryBuilder interface {

	ToSql() string

	Get() (*sql.Rows, error)

	Limit(n int) QueryBuilder

	Offset(n int) QueryBuilder

	From(from string) QueryBuilder

	Select(args ...string) QueryBuilder

	SelectRaw(args ...string) QueryBuilder

	Where(param ...interface{}) QueryBuilder

	OrWhere(param ...interface{}) QueryBuilder

	WhereColumn(args ...interface{}) QueryBuilder

	OrWhereColumn(args ...interface{}) QueryBuilder

	WhereRaw(condition string, bindings ...interface{}) QueryBuilder

	OrWhereRaw(condition string, bindings ...interface{}) QueryBuilder

	WhereNull(col string) QueryBuilder

	OrWhereNull(col string) QueryBuilder

	WhereNotNull(col string) QueryBuilder

	OrWhereNotNull(col string) QueryBuilder

	WhereIn(column string, values ...interface{}) QueryBuilder

	OrWhereIn(column string, values ...interface{}) QueryBuilder

	WhereNotIn(column string, values ...interface{}) QueryBuilder

	OrWhereNotIn(column string, values ...interface{}) QueryBuilder

	WhereBetween(column string, from interface{}, to interface{}) QueryBuilder

	OrWhereBetween(column string, from interface{}, to interface{}) QueryBuilder

	WhereNotBetween(column string, from interface{}, to interface{}) QueryBuilder

	OrWhereNotBetween(column string, from interface{}, to interface{}) QueryBuilder

	WhereDate(args ...interface{}) QueryBuilder

	OrWhereDate(args ...interface{}) QueryBuilder

	WhereTime(args ...interface{}) QueryBuilder

	OrWhereTime(args ...interface{}) QueryBuilder

	WhereDay(args ...interface{}) QueryBuilder

	OrWhereDay(args ...interface{}) QueryBuilder

	WhereMonth(args ...interface{}) QueryBuilder

	OrWhereMonth(args ...interface{}) QueryBuilder

	WhereYear(args ...interface{}) QueryBuilder

	OrWhereYear(args ...interface{}) QueryBuilder

	GroupBy(args ...string) QueryBuilder

	OrderBy(column string, direction string) QueryBuilder

	OrderByDesc(column string) QueryBuilder

	OrderByRaw(sql string, bindings ...interface{}) QueryBuilder

	Count(column string) (*sql.Rows, error)

	Min(column string) (*sql.Rows, error)

	Max(column string) (*sql.Rows, error)

	Sum(column string) (*sql.Rows, error)

	Avg(column string) (*sql.Rows, error)
}
