package kernel

import (
	"database/contracts"
	"database/kernel/config"
)

type Manager struct {
	config *config.DatabaseConfig
	factory contracts.ConnectionFactory
	connections map[string]contracts.Connection
}

func (m *Manager) Connection(name string) contracts.Connection {

	if len(name) <= 0 {

		name = m.getDefaultDriver()
	}

	_, hasConnection := m.connections[name]

	if ! hasConnection {

		m.connections[name] = m.makeConnection(name)
	}

	return m.connections[name]
}

func (m *Manager) getDefaultDriver() string {

	return m.config.Default;
}

func (m *Manager) makeConnection(name string) contracts.Connection {

	configDriver := m.configuration(name)

	return m.factory.Make(configDriver)
}

func (m *Manager) configuration(name string) *config.DatabaseDriver {

	driverConfig, success := m.config.Connections[name]

	if ! success {
		panic("Error: not found driver for " + name)
	}

	return &driverConfig
}