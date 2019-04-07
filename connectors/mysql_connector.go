package connectors

import (
	"database/kernel/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MySqlConnector struct {
	*Connector
	config *config.DatabaseDriver
}

func NewMysqlConnector(config *config.DatabaseDriver) *MySqlConnector {
	return &MySqlConnector{
		config: config,
		Connector: NewConnector(config),
	}
}

func (m *MySqlConnector) Connect() *sql.DB {

	connection := m.Connector.CreateConnection("mysql")

	m.configureEncoding(connection)

	m.configureTimezone(connection)

	m.setModes(connection)

	return connection
}

func (m *MySqlConnector) configureEncoding(connection *sql.DB) {

	stmt, err := connection.Prepare("set names '" + m.config.Charset + "' collate " + m.config.Collation)

	prepareError(err)

	stmt.Exec()
}

func (m *MySqlConnector) configureTimezone(connection *sql.DB) {

	if len(m.config.Timezone) > 0 {

		stmt, err := connection.Prepare("set time_zone='" + m.config.Timezone + "'")

		prepareError(err)

		stmt.Exec()
	}
}

func (m *MySqlConnector) setModes(connection *sql.DB) {

	if m.config.Strict == true {

		stmt, err := connection.Prepare(m.getStrictMode())

		prepareError(err)

		stmt.Exec()
	}
}

func (m *MySqlConnector) getStrictMode() string {

	return "set session sql_mode='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'";
}

func prepareError(err error) {
	if err != nil {
		panic(err)
	}
}