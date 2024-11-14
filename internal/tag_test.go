package internal_test

import (
	"reflect"
	"testing"

	"github.com/berquerant/structconfig/internal"
	"github.com/stretchr/testify/assert"
)

func TestTag(t *testing.T) {
	type A struct {
		T1 int `name:"t1"`
		T2 int `name:"t2" usage:"t2usage"`
		T3 int `name:"t3" default:"1"`
		T4 int `aname:"t4"`
		T5 int
		T6 int `name:"-"`
	}

	v := reflect.TypeOf(A{})

	t.Run("T1", func(t *testing.T) {
		x := internal.NewTag(v.Field(0).Tag, "")
		name, ok := x.Name()
		assert.True(t, ok)
		assert.Equal(t, "t1", name)
	})
	t.Run("T2", func(t *testing.T) {
		x := internal.NewTag(v.Field(1).Tag, "")
		name, ok := x.Name()
		assert.True(t, ok)
		assert.Equal(t, "t2", name)
		assert.Equal(t, "t2usage", x.Usage())
		_, ok = x.Default()
		assert.False(t, ok)
	})
	t.Run("T3", func(t *testing.T) {
		x := internal.NewTag(v.Field(2).Tag, "")
		name, ok := x.Name()
		assert.True(t, ok)
		assert.Equal(t, "t3", name)
		dv, ok := x.Default()
		assert.True(t, ok)
		assert.Equal(t, "1", dv)
	})
	t.Run("T4", func(t *testing.T) {
		x := internal.NewTag(v.Field(3).Tag, "a")
		name, ok := x.Name()
		assert.True(t, ok)
		assert.Equal(t, "t4", name)
	})
	t.Run("T5", func(t *testing.T) {
		x := internal.NewTag(v.Field(4).Tag, "")
		_, ok := x.Name()
		assert.False(t, ok)
	})
	t.Run("T5", func(t *testing.T) {
		x := internal.NewTag(v.Field(5).Tag, "")
		_, ok := x.Name()
		assert.False(t, ok)
	})
}
