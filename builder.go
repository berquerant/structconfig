package structconfig

func NewBuilder[T any](sc *StructConfig[T], merger *Merger[T]) *Builder[T] {
	return &Builder[T]{
		sc:     sc,
		merger: merger,
	}
}

type Builder[T any] struct {
	sc     *StructConfig[T]
	merger *Merger[T]
	chain  []func(*StructConfig[T]) (*T, error)
}

// Add adds a Config generator to the Builder.
func (b *Builder[T]) Add(f func(*StructConfig[T]) (*T, error)) *Builder[T] {
	b.chain = append(b.chain, f)
	return b
}

// Build generates a Config in order from the generators added earlier and override them accordingly.
func (b *Builder[T]) Build() (*T, error) {
	configList := make([]*T, len(b.chain))
	for i, c := range b.chain {
		x, err := b.newConfig(c)
		if err != nil {
			return nil, err
		}
		configList[i] = x
	}

	r, err := b.newDefault()
	if err != nil {
		return nil, err
	}
	for _, c := range configList {
		x, err := b.merger.Merge(*r, *c)
		if err != nil {
			return nil, err
		}
		r = &x
	}

	return r, nil
}

func (b *Builder[T]) newConfig(f func(*StructConfig[T]) (*T, error)) (*T, error) {
	d, err := b.newDefault()
	if err != nil {
		return nil, err
	}
	v, err := f(b.sc)
	if err != nil {
		return nil, err
	}
	r, err := b.merger.Merge(*d, *v)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (b *Builder[T]) newDefault() (*T, error) {
	var t T
	if err := b.sc.FromDefault(&t); err != nil {
		return nil, err
	}
	return &t, nil
}
