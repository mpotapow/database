package connections

import (
	"database/contracts"
	"database/kernel/config"
	"database/query"
	"database/sql"
)

type Connection struct {
	pdo          *sql.DB
	config       *config.DatabaseDriver
	queryGrammar contracts.Grammar
}

func NewConnection(pdo *sql.DB, config *config.DatabaseDriver, grammar contracts.Grammar) *Connection {

	return &Connection{
		pdo:          pdo,
		config:       config,
		queryGrammar: grammar,
	}
}

func (c *Connection) GetPDO() *sql.DB {

	return c.pdo
}

func (c *Connection) GetGrammar() contracts.Grammar {

	return c.queryGrammar
}

func (c *Connection) Query() contracts.QueryBuilder {

	return query.NewBuilder(c, c.queryGrammar)
}

func (c *Connection) Table(table string) contracts.QueryBuilder {

	return c.Query().From(table)
}

func (c *Connection) Select(query string, bindings []interface{}) (*sql.Rows, error) {

	return c.query(c.prepareQuery(query), bindings)
}

func (c *Connection) Insert(query string, bindings []interface{}) sql.Result {

	return c.statement(c.prepareQuery(query), bindings)
}

func (c *Connection) Update(query string, bindings []interface{}) int64 {

	return c.affectingStatement(c.prepareQuery(query), bindings)
}

func (c *Connection) Delete(query string, bindings []interface{}) int64 {

	return c.affectingStatement(c.prepareQuery(query), bindings)
}

func (c *Connection) query(statement *sql.Stmt, bindings []interface{}) (*sql.Rows, error) {

	defer statement.Close()

	return statement.Query(bindings...)
}

func (c *Connection) statement(statement *sql.Stmt, bindings []interface{}) sql.Result {

	defer statement.Close()

	res, err := statement.Exec(bindings...)
	prepareError(err)

	return res
}

func (c *Connection) affectingStatement(statement *sql.Stmt, bindings []interface{}) int64 {

	defer statement.Close()

	res, err := statement.Exec(bindings...)
	prepareError(err)

	cont, err := res.RowsAffected()
	if err != nil {
		return 0
	}

	return cont
}

func (c *Connection) prepareQuery(query string) *sql.Stmt {

	statement, err := c.pdo.Prepare(query)
	prepareError(err)

	return statement
}

func (c *Connection) Transaction(args ...interface{}) contracts.TransactionConnection {

	if len(args) > 0 {
		tc := NewTransactionConnection(c)
		return tc.Transaction(args...)
	}

	return NewTransactionConnection(c)
}

func prepareError(err error) {

	if err != nil {
		panic(err)
	}
}
