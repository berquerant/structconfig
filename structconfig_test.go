package structconfig_test

import (
	"os"
	"testing"
	"time"

	"github.com/berquerant/structconfig"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestStructConfig_AdditionalTypes(t *testing.T) {
	type Config struct {
		BoolSlice    []bool        `name:"bool_slice" default:"true,false"`
		IntSlice     []int         `name:"int_slice" default:"1,2"`
		Int32Slice   []int32       `name:"int32_slice" default:"3,4"`
		Int64Slice   []int64       `name:"int64_slice" default:"5,6"`
		UintSlice    []uint        `name:"uint_slice" default:"7,8"`
		Float32Slice []float32     `name:"float32_slice" default:"1.1,2.2"`
		Float64Slice []float64     `name:"float64_slice" default:"3.3,4.4"`
		StringSlice  []string      `name:"string_slice" default:"a,b"`
		Duration     time.Duration `name:"duration" default:"10s"`
		Time         time.Time     `name:"time" default:"2026-08-15T12:00:00Z"`
		Verbosity    int           `name:"verbose" short:"v" count:"true"`
	}

	t.Run("FromDefault", func(t *testing.T) {
		sc := structconfig.New[Config]()
		var got Config
		err := sc.FromDefault(&got)
		assert.NoError(t, err)

		assert.Equal(t, []bool{true, false}, got.BoolSlice)
		assert.Equal(t, []int{1, 2}, got.IntSlice)
		assert.Equal(t, []int32{3, 4}, got.Int32Slice)
		assert.Equal(t, []int64{5, 6}, got.Int64Slice)
		assert.Equal(t, []uint{7, 8}, got.UintSlice)
		assert.Equal(t, []float32{1.1, 2.2}, got.Float32Slice)
		assert.Equal(t, []float64{3.3, 4.4}, got.Float64Slice)
		assert.Equal(t, []string{"a", "b"}, got.StringSlice)
		assert.Equal(t, 10*time.Second, got.Duration)
		assert.True(t, time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC).Equal(got.Time))
		assert.Equal(t, 0, got.Verbosity)
	})

	t.Run("FromEnv", func(t *testing.T) {
		envs := map[string]string{
			"BOOL_SLICE":   "false,true",
			"INT_SLICE":    "10,20",
			"DURATION":     "1m",
			"TIME":         "2027-01-01T00:00:00Z",
			"VERBOSE":      "3",
		}
		for k, v := range envs {
			os.Setenv(k, v)
		}
		defer func() {
			for k := range envs {
				os.Unsetenv(k)
			}
		}()

		sc := structconfig.New[Config]()
		var got Config
		err := sc.FromEnv(&got)
		assert.NoError(t, err)

		assert.Equal(t, []bool{false, true}, got.BoolSlice)
		assert.Equal(t, []int{10, 20}, got.IntSlice)
		assert.Equal(t, time.Minute, got.Duration)
		assert.True(t, time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC).Equal(got.Time))
		assert.Equal(t, 3, got.Verbosity)
	})

	t.Run("FromFlags", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		sc := structconfig.New[Config]()
		err := sc.SetFlags(fs)
		assert.NoError(t, err)

		args := []string{
			"--int_slice", "100,200",
			"--duration", "500ms",
			"-v", "-v", "-v",
		}
		err = fs.Parse(args)
		assert.NoError(t, err)

		var got Config
		err = sc.FromFlags(&got, fs)
		assert.NoError(t, err)

		assert.Equal(t, []int{100, 200}, got.IntSlice)
		assert.Equal(t, 500*time.Millisecond, got.Duration)
		assert.Equal(t, 3, got.Verbosity)
	})

	t.Run("Merger", func(t *testing.T) {
		merger := structconfig.NewMerger[Config]()

		left := Config{
			IntSlice:  []int{10, 20},
			Duration:  10 * time.Second, // default
			Verbosity: 1,
		}
		right := Config{
			IntSlice:  []int{1, 2}, // default
			Duration:  2 * time.Minute,
			Verbosity: 0, // default
		}

		merged, err := merger.Merge(left, right)
		assert.NoError(t, err)

		assert.Equal(t, []int{10, 20}, merged.IntSlice)
		assert.Equal(t, 2*time.Minute, merged.Duration)
		assert.Equal(t, 1, merged.Verbosity)
	})
}
