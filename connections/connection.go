package connections

import (
	"database/contracts"
	"database/kernel/config"
	"database/query"
	"database/sql"
	"fmt"
)

type Connection struct {
	pdo *sql.DB
	config *config.DatabaseDriver
	queryGrammar contracts.Grammar
}

func NewConnection(pdo *sql.DB, config *config.DatabaseDriver, grammar contracts.Grammar) *Connection {

	return &Connection{
		pdo: pdo,
		config: config,
		queryGrammar: grammar,
	}
}

func (c *Connection) Query() contracts.QueryBuilder {

	return query.NewBuilder(c, c.queryGrammar)
}

func (c *Connection) Select(query string, bindings []interface{}) (*sql.Rows, error) {

	fmt.Println("== SELECT ==", query)

	statement, err := c.pdo.Prepare(query)
	prepareError(err)

	defer statement.Close()

	return statement.Query(bindings...)
}

func (c *Connection) Insert(query string, bindings []interface{}) sql.Result {

	fmt.Println("== INSERT ==", query)

	statement, err := c.pdo.Prepare(query)
	prepareError(err)

	defer statement.Close()

	res, err := statement.Exec(bindings...)
	prepareError(err)

	return res
}

func prepareError(err error) {

	if err != nil {
		panic(err)
	}
}