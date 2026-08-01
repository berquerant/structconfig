package structconfig_test

import (
	"bytes"
	"log/slog"
	"os"
	"testing"

	"github.com/berquerant/structconfig"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	type Config struct {
		Host string `name:"host" default:"localhost"`
		Port int    `name:"port" default:"8080"`
	}

	t.Run("LevelInfo filters out debug logs", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))

		sc := structconfig.New[Config](structconfig.WithLogger(logger))
		var cfg Config
		err := sc.FromDefault(&cfg)
		assert.NoError(t, err)
		assert.Empty(t, buf.String(), "Info level should not log debug level structconfig messages")
	})

	t.Run("FromDefault logging at Debug level", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

		sc := structconfig.New[Config](structconfig.WithLogger(logger))
		var cfg Config
		err := sc.FromDefault(&cfg)
		assert.NoError(t, err)
		assert.Equal(t, "localhost", cfg.Host)
		assert.Equal(t, 8080, cfg.Port)

		logOutput := buf.String()
		assert.Contains(t, logOutput, "level=DEBUG")
		assert.Contains(t, logOutput, "structconfig: set field")
		assert.Contains(t, logOutput, "source=default")
		assert.Contains(t, logOutput, "field=Host")
		assert.Contains(t, logOutput, "value=localhost")
		assert.Contains(t, logOutput, "field=Port")
		assert.Contains(t, logOutput, "value=8080")
	})

	t.Run("FromEnv logging at Debug level", func(t *testing.T) {
		os.Setenv("HOST", "127.0.0.1")
		defer os.Unsetenv("HOST")

		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

		sc := structconfig.New[Config](structconfig.WithLogger(logger))
		var cfg Config
		err := sc.FromEnv(&cfg)
		assert.NoError(t, err)
		assert.Equal(t, "127.0.0.1", cfg.Host)
		assert.Equal(t, 8080, cfg.Port)

		logOutput := buf.String()
		assert.Contains(t, logOutput, "level=DEBUG")
		assert.Contains(t, logOutput, "source=env")
		assert.Contains(t, logOutput, "env_var=HOST")
		assert.Contains(t, logOutput, "value=127.0.0.1")
		assert.Contains(t, logOutput, "source=default")
		assert.Contains(t, logOutput, "field=Port")
	})

	t.Run("FromFlags logging at Debug level", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		sc := structconfig.New[Config](structconfig.WithLogger(nil))
		assert.NoError(t, sc.SetFlags(fs))
		assert.NoError(t, fs.Parse([]string{"--host", "example.com"}))

		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
		scWithLogger := structconfig.New[Config](structconfig.WithLogger(logger))

		var cfg Config
		err := scWithLogger.FromFlags(&cfg, fs)
		assert.NoError(t, err)
		assert.Equal(t, "example.com", cfg.Host)

		logOutput := buf.String()
		assert.Contains(t, logOutput, "level=DEBUG")
		assert.Contains(t, logOutput, "source=flag")
		assert.Contains(t, logOutput, "flag=host")
		assert.Contains(t, logOutput, "value=example.com")
		assert.Contains(t, logOutput, "changed=true")
	})

	t.Run("Merger logging at Debug level", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

		merger := structconfig.NewMerger[Config](structconfig.WithLogger(logger))

		left := Config{Host: "left-host", Port: 8080}
		right := Config{Host: "localhost", Port: 9000}

		merged, err := merger.Merge(left, right)
		assert.NoError(t, err)
		assert.Equal(t, "left-host", merged.Host)
		assert.Equal(t, 9000, merged.Port)

		logOutput := buf.String()
		assert.Contains(t, logOutput, "level=DEBUG")
		assert.Contains(t, logOutput, "structconfig: merged field")
		assert.Contains(t, logOutput, "field=Host")
		assert.Contains(t, logOutput, "source=left")
		assert.Contains(t, logOutput, "field=Port")
		assert.Contains(t, logOutput, "source=right")
	})
}
