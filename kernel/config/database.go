package config

type DatabaseConfig struct {
	Default string
	Connections map[string]DatabaseDriver
}

type DatabaseDriver struct {
	Driver string
	Host string
	Port string
	Database string
	Password string
	Username string

	Strict bool
	Timezone string

	Charset string
	Collation string
}