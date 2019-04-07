package contracts

type Manager interface {

	Connection(name string) Connection
}