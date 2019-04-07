package types

type SelectType interface {
	ToString() string
}

type SelectString struct {
	value string
}

func (s *SelectString) ToString() string {
	return s.value
}

func NewSelectString(value string) *SelectString {
	return &SelectString{
		value: value,
	}
}

type SelectRawString struct {
	value string
}

func (s *SelectRawString) ToString() string {
	return s.value
}

func NewSelectRawString(value string) *SelectRawString {
	return &SelectRawString{
		value: value,
	}
}