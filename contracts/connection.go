package contracts

import "database/sql"

type Connection interface {

	Query() QueryBuilder

	Select(query string, bindings []interface{}) (*sql.Rows, error)
}
