package kernel

import (
	"database/contracts"
	"database/kernel/config"
)

func NewDatabaseConfig() *config.DatabaseConfig {

	return new(config.DatabaseConfig)
}

func NewDatabaseManager(c *config.DatabaseConfig) contracts.Manager {

	return &Manager{
		config: c,
		factory: newConnectionFactory(),
		connections: make(map[string]contracts.Connection),
	}
}

func newConnectionFactory() contracts.ConnectionFactory {

	return &ConnectionFactory{}
}
