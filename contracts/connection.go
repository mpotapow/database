package contracts

import "database/sql"

type TransactionConnection interface {
	Connection
	Transactable
}

type Connection interface {

	GetPDO() *sql.DB

	GetGrammar() Grammar

	Query() QueryBuilder

	Table(table string) QueryBuilder

	Select(query string, bindings []interface{}) (*sql.Rows, error)

	Insert(query string, bindings []interface{}) sql.Result

	Update(query string, bindings []interface{}) int64

	Delete(query string, bindings []interface{}) int64

	Transaction(args ...interface{}) TransactionConnection

	Statement(sql string, bindings []interface{}) sql.Result
}
