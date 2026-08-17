package internal

import (
	"encoding/csv"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

var (
	// Skip conv and callback but no error in TryParse.
	ErrSkipParse = errors.New("SkipParse")
)

func TryParse[S, T any](
	conv func(S) (T, error),
	callback func(S, T) error,
) func(S) error {
	return func(s S) error {
		v, err := conv(s)
		switch {
		case err == nil:
			return callback(s, v)
		case errors.Is(err, ErrSkipParse):
			return nil
		default:
			return err
		}
	}
}

func ParseInt[T constraints.Signed](s string) (T, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	return T(v), err
}

func ParseUint[T constraints.Unsigned](s string) (T, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	return T(v), err
}

func ParseFloat[T constraints.Float](s string) (T, error) {
	v, err := strconv.ParseFloat(s, 64)
	return T(v), err
}

func ParseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func ParseSlice[T any](s string, elemConv func(string) (T, error)) ([]T, error) {
	if s == "" {
		return []T{}, nil
	}
	r := csv.NewReader(strings.NewReader(s))
	r.TrimLeadingSpace = true
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	var result []T
	for _, row := range records {
		for _, item := range row {
			v, err := elemConv(item)
			if err != nil {
				return nil, err
			}
			result = append(result, v)
		}
	}
	return result, nil
}

func ParseStringSlice(s string, split bool, sep string) ([]string, error) {
	if !split {
		if s == "" {
			return []string{}, nil
		}
		return []string{s}, nil
	}
	if sep == "" || sep == "," {
		return ParseSlice(s, func(x string) (string, error) { return x, nil })
	}
	if s == "" {
		return []string{}, nil
	}
	return strings.Split(s, sep), nil
}

func IsSupportedKind(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Supported interface {
	~bool | constraints.Signed | Unsigned | ~string
}
