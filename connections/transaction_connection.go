package connections

import (
	"database/concerns"
	"database/contracts"
	"database/query"
	"database/sql"
)

type TransactionConnection struct {
	*Connection
	*concerns.ManagesTransactions
}

type TransactionCallback = func(tc contracts.TransactionConnection) error

func NewTransactionConnection(c *Connection) contracts.TransactionConnection {

	tm := concerns.NewManagesTransactions(c.GetPDO(), c.GetGrammar())

	return &TransactionConnection{
		c,
		tm,
	}
}

func (tc *TransactionConnection) Query() contracts.QueryBuilder {

	return query.NewBuilder(tc, tc.GetGrammar())
}

func (tc *TransactionConnection) Table(table string) contracts.QueryBuilder {

	return tc.Query().From(table)
}

func (tc *TransactionConnection) Select(query string, bindings []interface{}) (*sql.Rows, error) {

	return tc.query(tc.prepareQuery(query), bindings)
}

func (tc *TransactionConnection) Insert(query string, bindings []interface{}) sql.Result {

	return tc.statement(tc.prepareQuery(query), bindings)
}

func (tc *TransactionConnection) Update(query string, bindings []interface{}) int64 {

	return tc.affectingStatement(tc.prepareQuery(query), bindings)
}

func (tc *TransactionConnection) Delete(query string, bindings []interface{}) int64 {

	return tc.affectingStatement(tc.prepareQuery(query), bindings)
}

func (tc *TransactionConnection) prepareQuery(query string) *sql.Stmt {

	var (
		err       error
		statement *sql.Stmt
	)

	if tc.TransactionLevel() > 0 {
		statement, err = tc.GetTxPDO().Prepare(query)
	} else {
		statement, err = tc.GetPDO().Prepare(query)
	}

	prepareError(err)

	return statement
}

func (tc *TransactionConnection) Transaction(args ...interface{}) contracts.TransactionConnection {

	callback, ok := args[0].(TransactionCallback)
	if !ok {
		panic("Unresolved transaction params")
	}

	tc.BeginTransaction()

	if err := callback(tc); err != nil {
		tc.RollBack(nil)
	} else {
		tc.Commit()
	}

	return tc
}
