package contracts

type Grammar interface {

	CompileSelect(b QueryBuilder) string

	CompileInsert(b QueryBuilder, values []map[string]interface{}, columns []string) string

	Wrap(v string) string
}
