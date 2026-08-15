package internal_test

import (
	"errors"
	"testing"
	"time"

	"github.com/berquerant/structconfig/internal"
	"github.com/stretchr/testify/assert"
)

func TestTryParse(t *testing.T) {
	var (
		errConv     = errors.New("conv err")
		errCallback = errors.New("callback err")
	)

	for _, tc := range []struct {
		title         string
		conv          func(string) (int, error)
		callback      func(string, int) error
		input         string
		wantCallback1 string
		wantCallback2 int
		wantCalled    bool
		wantErr       bool
	}{
		{
			title: "success",
			conv:  func(_ string) (int, error) { return 1, nil },
			callback: func(_ string, _ int) error {
				return nil
			},
			input:         "a",
			wantCallback1: "a",
			wantCallback2: 1,
			wantCalled:    true,
			wantErr:       false,
		},
		{
			title: "callback failed",
			conv:  func(_ string) (int, error) { return 1, nil },
			callback: func(_ string, _ int) error {
				return errCallback
			},
			input:         "a",
			wantCallback1: "a",
			wantCallback2: 1,
			wantCalled:    true,
			wantErr:       true,
		},
		{
			title: "conv failed",
			conv:  func(_ string) (int, error) { return 0, errConv },
			callback: func(_ string, _ int) error {
				return nil
			},
			input:      "a",
			wantCalled: false,
			wantErr:    true,
		},
		{
			title: "skipped",
			conv:  func(_ string) (int, error) { return 0, internal.ErrSkipParse },
			callback: func(_ string, _ int) error {
				return nil
			},
			input:      "a",
			wantCalled: false,
			wantErr:    false,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			var (
				callback1 string
				callback2 int
				called    bool
			)
			err := internal.TryParse(
				tc.conv,
				func(s string, i int) error {
					called = true
					callback1 = s
					callback2 = i
					return tc.callback(s, i)
				},
			)(tc.input)

			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tc.wantCalled, called)
			if tc.wantCalled {
				assert.Equal(t, tc.wantCallback1, callback1)
				assert.Equal(t, tc.wantCallback2, callback2)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	for _, tc := range []struct {
		title   string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{
			title: "seconds",
			input: "10s",
			want:  10 * time.Second,
		},
		{
			title: "minutes",
			input: "5m",
			want:  5 * time.Minute,
		},
		{
			title:   "invalid",
			input:   "invalid",
			wantErr: true,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got, err := internal.ParseDuration(tc.input)
			if tc.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestParseTime(t *testing.T) {
	for _, tc := range []struct {
		title   string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			title: "RFC3339",
			input: "2026-08-15T12:00:00Z",
			want:  time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			title:   "invalid",
			input:   "2026/08/15",
			wantErr: true,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got, err := internal.ParseTime(tc.input)
			if tc.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.True(t, tc.want.Equal(got))
		})
	}
}

func TestParseSlice(t *testing.T) {
	for _, tc := range []struct {
		title   string
		input   string
		want    []int
		wantErr bool
	}{
		{
			title: "empty string",
			input: "",
			want:  []int{},
		},
		{
			title: "empty brackets",
			input: "[]",
			want:  []int{},
		},
		{
			title: "csv",
			input: "1, 2, 3",
			want:  []int{1, 2, 3},
		},
		{
			title: "brackets with elements",
			input: "[1, 2, 3]",
			want:  []int{1, 2, 3},
		},
		{
			title:   "invalid element",
			input:   "1, a, 3",
			wantErr: true,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got, err := internal.ParseSlice(tc.input, internal.ParseInt[int])
			if tc.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
