package internal

import (
	"log/slog"
	"reflect"
)

// NewMerger returns a new Merger.
//
// anyCallback parses "default" tag value and set it.
// anyEqual reports true if left equals right when kind of arguments are not supported.
// prefix adds a prefix to "default" tag name.
func NewMerger[T any](
	anyCallback func(StructField, string, func() reflect.Value) error,
	anyEqual func(left, right any) (bool, error),
	prefix string,
	logger *slog.Logger,
) *Merger[T] {
	return &Merger[T]{
		anyCallback: anyCallback,
		anyEqual:    anyEqual,
		prefix:      prefix,
		logger:      logger,
	}
}

type Merger[T any] struct {
	anyCallback func(StructField, string, func() reflect.Value) error
	anyEqual    func(left, right any) (bool, error)
	prefix      string
	logger      *slog.Logger
}

func (m Merger[T]) newReceptor(ptr *T) (*PairsReceptor, error) {
	return DefaultReceptor(ptr, m.anyCallback, nil)
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
	if lType != rType {
		return false, nil
	}
	if IsSupportedKind(lType.Kind()) {
		return left == right, nil
	}
	if lType.Kind() == reflect.Slice && IsSupportedKind(lType.Elem().Kind()) {
		return reflect.DeepEqual(left, right), nil
	}
	if lType == durationTyp || lType == timeTyp {
		return reflect.DeepEqual(left, right), nil
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

func (m Merger[T]) selectNonDefaultValue(defaultVal, candidateVal reflect.Value) (reflect.Value, bool, error) {
	ok, err := m.equal(defaultVal.Interface(), candidateVal.Interface())
	if err != nil {
		return reflect.Value{}, false, err
	}
	if !ok {
		return candidateVal, true, nil
	}
	return reflect.Value{}, false, nil
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
		if _, ok := f.Tag().Name(); !ok {
			// ignore the field without 'name' tag
			continue
		}

		name := f.Name()
		fv := vv.Elem().FieldByName(name)

		sources := []struct {
			name  string
			val   reflect.Value
		}{
			{"right", rValue.FieldByName(name)},
			{"left", lValue.FieldByName(name)},
		}

		for _, src := range sources {
			newVal, selected, err := m.selectNonDefaultValue(fv, src.val)
			if err != nil {
				return v, err
			}
			if selected {
				fv.Set(newVal)
				if m.logger != nil {
					m.logger.Debug(
						"structconfig: merged field",
						slog.String("field", name),
						slog.String("source", src.name),
						slog.Any("value", newVal.Interface()),
					)
				}
				break
			}
		}
	}
	return v, nil
}
