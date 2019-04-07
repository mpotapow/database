package contracts

type Grammar interface {

	CompileSelect(b QueryBuilder) string
}
