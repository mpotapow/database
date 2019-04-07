package contracts

import "database/kernel/config"

type ConnectionFactory interface {

	Make(config *config.DatabaseDriver) Connection
}