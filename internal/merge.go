package internal

import "reflect"

// NewMerger returns a new Merger.
//
// anyCallback parses "default" tag value and set it.
// anyEqual reports true if left equals right when kind of arguments are not supported.
// prefix adds a prefix to "default" tag name.
func NewMerger[T any](
	anyCallback func(StructField, string, func() reflect.Value) error,
	anyEqual func(left, right any) (bool, error),
	prefix string,
) *Merger[T] {
	return &Merger[T]{
		anyCallback: anyCallback,
		anyEqual:    anyEqual,
		prefix:      prefix,
	}
}

type Merger[T any] struct {
	anyCallback func(StructField, string, func() reflect.Value) error
	anyEqual    func(left, right any) (bool, error)
	prefix      string
}

func (m Merger[T]) newReceptor(ptr *T) (*PairsReceptor, error) {
	return DefaultReceptor(ptr, m.anyCallback)
}

func (m Merger[T]) getType() (*Type, error) {
	var value T
	return NewType(value, m.prefix)
}

func (m Merger[T]) equal(left, right any) (bool, error) {
	if (left == nil) != (right == nil) {
		return false, nil
	}
	if left == nil { // right == nil
		return true, nil
	}

	lType, rType := reflect.TypeOf(left), reflect.TypeOf(right)
	if lType.Kind() != rType.Kind() {
		return false, nil
	}
	if IsSupportedKind(lType.Kind()) {
		return left == right, nil
	}
	if eq := m.anyEqual; eq != nil {
		return eq(left, right)
	}
	return false, nil
}

func (m Merger[T]) defaultValue() (T, error) {
	var value T
	typ, err := m.getType()
	if err != nil {
		return value, err
	}

	r, err := m.newReceptor(&value)
	if err != nil {
		return value, err
	}

	if err := typ.Accept(r); err != nil {
		return value, err
	}

	return value, nil
}

// Merge values based on the 'default' tag values.
// For each field, if the right value is not the default, use it; if not, use the left value.
// If that is also the default, set the default value. Return this instance.
func (m Merger[T]) Merge(left, right T) (T, error) {
	v, err := m.defaultValue()
	if err != nil {
		return v, err
	}

	typ, err := m.getType()
	if err != nil {
		return v, err
	}

	lValue, rValue := reflect.ValueOf(left), reflect.ValueOf(right)
	vv := reflect.ValueOf(&v)
	for _, f := range typ.Fields() {
		name := f.Name()
		fv := vv.Elem().FieldByName(name)

		{
			rv := rValue.FieldByName(name)
			ok, err := m.equal(fv.Interface(), rv.Interface())
			if err != nil {
				return v, err
			}
			if !ok {
				// overwrite by not default value from right
				fv.Set(rv)
				continue
			}
		}
		{
			lv := lValue.FieldByName(name)
			ok, err := m.equal(fv.Interface(), lv.Interface())
			if err != nil {
				return v, err
			}
			if !ok {
				// overwrite by not default value from left
				fv.Set(lv)
				continue
			}
		}
	}
	return v, nil
}
