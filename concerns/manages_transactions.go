package concerns

import (
	"database/contracts"
	"database/sql"
	"fmt"
)

type ManagesTransactions struct {
	pdo          *sql.DB
	tx           *sql.Tx
	grammar      contracts.Grammar
	transactions int
}

func NewManagesTransactions(pdo *sql.DB, grammar contracts.Grammar) *ManagesTransactions {

	return &ManagesTransactions{
		pdo:          pdo,
		grammar:      grammar,
		transactions: 0,
	}
}

func (mt *ManagesTransactions) BeginTransaction() {

	mt.createTransaction()

	mt.transactions++
}

func (mt *ManagesTransactions) createTransaction() {

	if mt.TransactionLevel() == 0 {
		tx, err := mt.pdo.Begin()
		prepareError(err)

		mt.tx = tx
	} else {
		_, err := mt.tx.Exec(mt.grammar.CompileSavepoint(
			mt.formatTxName(mt.TransactionLevel() + 1),
		))
		prepareError(err)
	}
}

func (mt *ManagesTransactions) Commit() {

	if mt.TransactionLevel() == 1 {
		err := mt.tx.Commit()
		prepareError(err)
	}

	if mt.transactions--; mt.transactions < 0 {
		mt.transactions = 0
	}
}

func (mt *ManagesTransactions) RollBack(level interface{}) {

	toLevel := 0
	if level != nil {
		toLevel = level.(int)
	} else {
		toLevel = mt.transactions - 1
	}

	if toLevel < 0 || toLevel >= mt.transactions {
		return
	}

	if toLevel == 0 {
		err := mt.tx.Rollback()
		prepareError(err)
	} else {
		_, err := mt.tx.Exec(mt.grammar.CompileSavepointRollback(
			mt.formatTxName(toLevel + 1),
		))
		prepareError(err)
	}

	mt.transactions = toLevel
}

func (mt *ManagesTransactions) TransactionLevel() int {

	return mt.transactions
}

func (mt *ManagesTransactions) GetTxPDO() *sql.Tx {

	return mt.tx
}

func (mt *ManagesTransactions) formatTxName(cnt int) string {

	return fmt.Sprintf("trans%d", cnt)
}

func prepareError(err error) {

	if err != nil {
		panic(err)
	}
}
