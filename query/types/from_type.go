package types

type FromType interface {
	ToString() string
}

type FromString struct {
	value string
}

func (f *FromString) ToString() string {
	return f.value
}

func NewFromString(value string) *FromString {
	return &FromString{
		value: value,
	}
}

type FromRawString struct {
	value string
}

func (f *FromRawString) ToString() string {
	return f.value
}

func NewFromRawString(value string) *FromRawString {
	return &FromRawString{
		value: value,
	}
}