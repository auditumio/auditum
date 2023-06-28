package types

type BoolValue struct {
	Bool  bool
	Valid bool
}

func (v BoolValue) True() bool {
	return v.Valid && v.Bool
}

func (v BoolValue) False() bool {
	return v.Valid && !v.Bool
}
