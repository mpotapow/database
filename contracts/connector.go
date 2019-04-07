package contracts

import "database/sql"

type Connector interface {

	Connect() *sql.DB
}