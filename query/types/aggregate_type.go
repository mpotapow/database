package types

type AggregateType interface {
	GetColumns() string
	GetFunction() string
}

type Aggregate struct {
	columns string
	function string
}

func (a *Aggregate) GetColumns() string {
	return a.columns
}

func (a *Aggregate) GetFunction() string {
	return a.function
}

func NewAggregate(f string, c string) *Aggregate {

	return &Aggregate{
		columns: c,
		function: f,
	}
}