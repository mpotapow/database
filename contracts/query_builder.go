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

	WhereNull(col string) QueryBuilder

	OrWhereNull(col string) QueryBuilder

	WhereNotNull(col string) QueryBuilder

	OrWhereNotNull(col string) QueryBuilder

	WhereIn(column string, values []interface{}) QueryBuilder

	OrWhereIn(column string, values []interface{}) QueryBuilder

	WhereNotIn(column string, values []interface{}) QueryBuilder

	OrWhereNotIn(column string, values []interface{}) QueryBuilder

	GroupBy(args ...string) QueryBuilder

	OrderBy(column string, direction string) QueryBuilder

	Count(column string) (*sql.Rows, error)

	Min(column string) (*sql.Rows, error)

	Max(column string) (*sql.Rows, error)

	Sum(column string) (*sql.Rows, error)

	Avg(column string) (*sql.Rows, error)
}
