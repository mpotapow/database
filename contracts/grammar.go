package contracts

type Grammar interface {

	CompileSelect(b QueryBuilder) string

	CompileInsert(b QueryBuilder, values []map[string]interface{}, columns []string) string

	CompileUpdate(b QueryBuilder, values map[string]interface{}) string

	CompileDelete(b QueryBuilder) string

	CompileTruncate(b QueryBuilder) string

	CompileSavepoint(name string) string

	CompileSavepointRollback(name string) string

	Wrap(v string) string

	PrepareBindingsForUpdate(b QueryBuilder, bindings map[string][]interface{}, values map[string]interface{}) []interface{}

	PrepareBindingsForDelete(b QueryBuilder, bindings map[string][]interface{}) []interface{}
}
