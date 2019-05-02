package contracts

import "database/sql"

type Connection interface {

	Query() QueryBuilder

	Select(query string, bindings []interface{}) (*sql.Rows, error)

	Insert(query string, bindings []interface{}) sql.Result

	Update(query string, bindings []interface{}) int64

	Delete(query string, bindings []interface{}) int64
}
