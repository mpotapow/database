package connectors

import (
	"database/kernel/config"
	"database/sql"
)

type Connector struct {
	config *config.DatabaseDriver
}

func NewConnector(config *config.DatabaseDriver) *Connector {

	return &Connector{
		config: config,
	}
}

func (c *Connector) CreateConnection(dsn string) *sql.DB {

	connectParams := c.config.Host + ":" + c.config.Port
	authParams := c.config.Username + ":" + c.config.Password

	connection, err := sql.Open(dsn, authParams + "@tcp(" + connectParams + ")/" + c.config.Database)

	if err != nil {
		panic("No connection to server. Error: " + err.Error())
	}

	return connection
}
