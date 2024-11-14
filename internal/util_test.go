package internal_test

import (
	"errors"
	"testing"

	"github.com/berquerant/structconfig/internal"
	"github.com/stretchr/testify/assert"
)

func TestTryParse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var (
			callback1 string
			callback2 int
		)
		assert.Nil(t, internal.TryParse(
			func(_ string) (int, error) { return 1, nil },
			func(s string, i int) error {
				callback1 = s
				callback2 = i
				return nil
			},
		)("a"))
		assert.Equal(t, "a", callback1)
		assert.Equal(t, 1, callback2)
	})
	t.Run("callback failed", func(t *testing.T) {
		var (
			callback1 string
			callback2 int
		)
		assert.NotNil(t, internal.TryParse(
			func(_ string) (int, error) { return 1, nil },
			func(s string, i int) error {
				callback1 = s
				callback2 = i
				return errors.New("err")
			},
		)("a"))
		assert.Equal(t, "a", callback1)
		assert.Equal(t, 1, callback2)
	})
	t.Run("conv failed", func(t *testing.T) {
		var (
			callbacked bool
		)
		assert.NotNil(t, internal.TryParse(
			func(_ string) (int, error) { return 0, errors.New("err") },
			func(_ string, _ int) error {
				callbacked = true
				return nil
			},
		)("a"))
		assert.False(t, callbacked)
	})
	t.Run("skipped", func(t *testing.T) {
		var (
			callbacked bool
		)
		assert.Nil(t, internal.TryParse(
			func(_ string) (int, error) { return 0, internal.ErrSkipParse },
			func(_ string, _ int) error {
				callbacked = true
				return nil
			},
		)("a"))
		assert.False(t, callbacked)
	})
}
