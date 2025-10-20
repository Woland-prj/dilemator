package logger

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Woland-prj/dilemator/pkg/logger/handlers/slogpretty"
	"github.com/Woland-prj/dilemator/pkg/logger/handlers/slogzero"
)

var (
	ErrEmptyFile          = errors.New("file path cannot be empty when multiple output is enabled")
	ErrUnknownEnvironment = errors.New("unknown environment")
)

const (
	dirPerms  = 0o755
	filePerms = 0o644

	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Interface -.
type Interface interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
}

var _ Interface = (*slog.Logger)(nil)

// New creates a new slog.Logger, that implements logger.Interface,
// if env is "local" setup colored pretty output with debug level,
// if env is "dev" setup json output with debug level,
// if env is "prod" setup json output with info level.
// When multiple logging enabled log file will be recreated on init,
// when disabled file be ignored and logged only to os.Stdout.
func New(env string, multiple bool, file string) (Interface, error) {
	var log *slog.Logger

	var writers []io.Writer

	writers = append(writers, os.Stdout)

	if multiple {
		if file == "" {
			return nil, ErrEmptyFile
		}

		if err := os.MkdirAll(filepath.Dir(file), dirPerms); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		logFile, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePerms)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		writers = append(writers, logFile)
	}

	multiWriter := io.MultiWriter(writers...)

	switch env {
	case envLocal:
		log = setupPrettySlog(multiWriter)
	case envDev:
		log = slog.New(slogzero.NewZeroStyleJSONHandler(multiWriter, slog.LevelDebug))
	case envProd:
		log = slog.New(slogzero.NewZeroStyleJSONHandler(multiWriter, slog.LevelInfo))
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownEnvironment, env)
	}

	return log, nil
}

func setupPrettySlog(w io.Writer) *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(w)

	return slog.New(handler)
}
