package internal

import "reflect"

// Type is the metadata of struct.
type Type struct {
	typ    reflect.Type // Struct
	prefix string
}

// NewType constructs [Type] from struct value.
// prefix is for [Tag].
func NewType(v any, prefix string) (*Type, error) {
	t := reflect.TypeOf(v)
	if x := t.Kind(); x != reflect.Struct {
		return nil, JoinErrors(ErrNotStruct, Errorf("cannot accept type %s", t.Name()))
	}
	return &Type{
		typ:    t,
		prefix: prefix,
	}, nil
}

// Fields returns the metadata of all fields of the struct.
func (t Type) Fields() []StructField {
	xs := make([]StructField, t.typ.NumField())
	for i := range t.typ.NumField() {
		x := t.typ.Field(i)
		xs[i] = NewStructField(
			x.Name,
			x.Type.Kind(),
			NewTag(x.Tag, t.prefix),
		)
	}
	return xs
}

// Name returns the name of the struct.
func (t Type) Name() string {
	return t.typ.Name()
}

// Accept [Call] r on each of [Type.Fields].
func (t Type) Accept(r Receptor) error {
	for _, f := range t.Fields() {
		if err := Call(r, f); err != nil {
			return err
		}
	}
	return nil
}
