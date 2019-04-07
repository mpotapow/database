package connections

import (
	"database/kernel/config"
	"database/query/grammars"
	"database/sql"
)

type MySqlConnection struct {
	*Connection
}

func NewMysqlConnection(pdo *sql.DB, config *config.DatabaseDriver) *MySqlConnection {

	return &MySqlConnection{
		Connection: NewConnection(pdo, config, grammars.NewMysqlGrammar()),
	}
}