package kernel

import (
	"database/connections"
	"database/connectors"
	"database/contracts"
	"database/kernel/config"
	"database/sql"
)

type ConnectionFactory struct {

}

func (c *ConnectionFactory) Make(config *config.DatabaseDriver) contracts.Connection {

	return c.newConnection(config)
}

func (c *ConnectionFactory) newConnection(config *config.DatabaseDriver) contracts.Connection {

	pdo := c.createConnector(config).Connect()

	return c.createConnection(config, pdo)
}

func (c *ConnectionFactory) createConnector(config *config.DatabaseDriver) contracts.Connector {

	switch config.Driver {
		case "mysql":
			return connectors.NewMysqlConnector(config)

		default:
			panic("Unresolved database connector driver")
	}
}

func (c *ConnectionFactory) createConnection(config *config.DatabaseDriver, pdo *sql.DB) contracts.Connection {

	switch config.Driver {
		case "mysql":
			return connections.NewMysqlConnection(pdo, config)

		default:
			panic("Unresolved database connection driver")
	}
}