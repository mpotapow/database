package contracts

import "database/sql"

type Transactable interface {

	GetTxPDO() *sql.Tx

	BeginTransaction()

	Commit()

	RollBack(level interface{})

	TransactionLevel() int
}
