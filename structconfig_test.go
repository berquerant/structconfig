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
		Float64Slice    []float64     `name:"float64_slice" default:"3.3,4.4"`
		StringSlice     []string      `name:"string_slice" default:"a,b" split:"true"`
		RawStringSlice  []string      `name:"raw_string_slice" default:"a,b"`
		CustomSepSlice  []string      `name:"custom_sep_slice" default:"x;y" split:"true" sep:";"`
		Duration        time.Duration `name:"duration" default:"10s"`
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
		assert.Equal(t, []string{"a,b"}, got.RawStringSlice)
		assert.Equal(t, []string{"x", "y"}, got.CustomSepSlice)
		assert.Equal(t, 10*time.Second, got.Duration)
		assert.True(t, time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC).Equal(got.Time))
		assert.Equal(t, 0, got.Verbosity)
	})

	t.Run("FromEnv", func(t *testing.T) {
		envs := map[string]string{
			"BOOL_SLICE":       "false,true",
			"INT_SLICE":        "10,20",
			"INT32_SLICE":      "30,40",
			"INT64_SLICE":      "50,60",
			"UINT_SLICE":       "70,80",
			"FLOAT32_SLICE":    "10.5,20.5",
			"FLOAT64_SLICE":    "30.5,40.5",
			"STRING_SLICE":     "x,y,z",
			"RAW_STRING_SLICE": "x,y,z",
			"CUSTOM_SEP_SLICE": "x;y;z",
			"DURATION":         "1m",
			"TIME":             "2027-01-01T00:00:00Z",
			"VERBOSE":          "3",
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
		assert.Equal(t, []int32{30, 40}, got.Int32Slice)
		assert.Equal(t, []int64{50, 60}, got.Int64Slice)
		assert.Equal(t, []uint{70, 80}, got.UintSlice)
		assert.Equal(t, []float32{10.5, 20.5}, got.Float32Slice)
		assert.Equal(t, []float64{30.5, 40.5}, got.Float64Slice)
		assert.Equal(t, []string{"x", "y", "z"}, got.StringSlice)
		assert.Equal(t, []string{"x,y,z"}, got.RawStringSlice)
		assert.Equal(t, []string{"x", "y", "z"}, got.CustomSepSlice)
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
			"--bool_slice", "false,true",
			"--int_slice", "100,200",
			"--int32_slice", "300,400",
			"--int64_slice", "500,600",
			"--uint_slice", "700,800",
			"--float32_slice", "1.5,2.5",
			"--float64_slice", "3.5,4.5",
			"--string_slice", "foo,bar",
			"--string_slice", "baz",
			"--raw_string_slice", "foo,bar",
			"--raw_string_slice", "baz",
			"--custom_sep_slice", "foo;bar",
			"--custom_sep_slice", "baz",
			"--duration", "500ms",
			"--time", "2028-12-31T23:59:59Z",
			"-v", "-v", "-v",
		}
		err = fs.Parse(args)
		assert.NoError(t, err)

		var got Config
		err = sc.FromFlags(&got, fs)
		assert.NoError(t, err)

		assert.Equal(t, []bool{false, true}, got.BoolSlice)
		assert.Equal(t, []int{100, 200}, got.IntSlice)
		assert.Equal(t, []int32{300, 400}, got.Int32Slice)
		assert.Equal(t, []int64{500, 600}, got.Int64Slice)
		assert.Equal(t, []uint{700, 800}, got.UintSlice)
		assert.Equal(t, []float32{1.5, 2.5}, got.Float32Slice)
		assert.Equal(t, []float64{3.5, 4.5}, got.Float64Slice)
		assert.Equal(t, []string{"foo", "bar", "baz"}, got.StringSlice)
		assert.Equal(t, []string{"foo,bar", "baz"}, got.RawStringSlice)
		assert.Equal(t, []string{"foo", "bar", "baz"}, got.CustomSepSlice)
		assert.Equal(t, 500*time.Millisecond, got.Duration)
		assert.True(t, time.Date(2028, 12, 31, 23, 59, 59, 0, time.UTC).Equal(got.Time))
		assert.Equal(t, 3, got.Verbosity)
	})

	t.Run("Merger", func(t *testing.T) {
		merger := structconfig.NewMerger[Config]()

		left := Config{
			BoolSlice:      []bool{true},
			IntSlice:       []int{10, 20},
			Int32Slice:     []int32{30},
			Int64Slice:     []int64{5, 6}, // default
			UintSlice:      []uint{70},
			Float32Slice:   []float32{1.1, 2.2}, // default
			Float64Slice:   []float64{30.5},
			StringSlice:    []string{"a", "b"}, // default
			RawStringSlice: []string{"a,b"},    // default
			CustomSepSlice: []string{"x", "y"}, // default
			Duration:       10 * time.Second,   // default
			Time:           time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC),
			Verbosity:      1,
		}
		right := Config{
			BoolSlice:      []bool{true, false}, // default
			IntSlice:       []int{1, 2},         // default
			Int32Slice:     []int32{3, 4},       // default
			Int64Slice:     []int64{500, 600},
			UintSlice:      []uint{7, 8}, // default
			Float32Slice:   []float32{10.5},
			Float64Slice:   []float64{3.3, 4.4}, // default
			StringSlice:    []string{"win"},
			RawStringSlice: []string{"win_raw"},
			CustomSepSlice: []string{"win_sep"},
			Duration:       2 * time.Minute,
			Time:           time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC), // default
			Verbosity:      0,                                             // default
		}

		merged, err := merger.Merge(left, right)
		assert.NoError(t, err)

		assert.Equal(t, []bool{true}, merged.BoolSlice)
		assert.Equal(t, []int{10, 20}, merged.IntSlice)
		assert.Equal(t, []int32{30}, merged.Int32Slice)
		assert.Equal(t, []int64{500, 600}, merged.Int64Slice)
		assert.Equal(t, []uint{70}, merged.UintSlice)
		assert.Equal(t, []float32{10.5}, merged.Float32Slice)
		assert.Equal(t, []float64{30.5}, merged.Float64Slice)
		assert.Equal(t, []string{"win"}, merged.StringSlice)
		assert.Equal(t, []string{"win_raw"}, merged.RawStringSlice)
		assert.Equal(t, []string{"win_sep"}, merged.CustomSepSlice)
		assert.Equal(t, 2*time.Minute, merged.Duration)
		assert.True(t, time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC).Equal(merged.Time))
		assert.Equal(t, 1, merged.Verbosity)
	})
}
