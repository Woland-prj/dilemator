package slogzero

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strings"
	"time"
)

type ZeroStyleJSONHandler struct {
	w io.Writer
	l slog.Level
}

func NewZeroStyleJSONHandler(w io.Writer, level slog.Level) *ZeroStyleJSONHandler {
	return &ZeroStyleJSONHandler{w: w, l: level}
}

func (h *ZeroStyleJSONHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.l
}

//nolint:gocritic // Should implement interface slog.Handler
func (h *ZeroStyleJSONHandler) Handle(_ context.Context, r slog.Record) error {
	// Определяем caller
	var caller string

	if r.PC != 0 {
		fn := runtime.FuncForPC(r.PC)
		if fn != nil {
			file, line := fn.FileLine(r.PC)
			caller = fmt.Sprintf("%s:%d", file, line)
		}
	}

	// Приводим уровень к нижнему регистру
	level := strings.ToLower(r.Level.String())

	// Преобразуем время в UTC без миллисекунд
	ts := r.Time.UTC().Truncate(time.Second).Format(time.RFC3339)

	// Собираем сообщение
	entry := map[string]any{
		"level":   level,
		"time":    ts,
		"caller":  caller,
		"message": r.Message,
	}

	// Включаем атрибуты (если есть)
	r.Attrs(func(a slog.Attr) bool {
		entry[a.Key] = a.Value.Any()

		return true
	})

	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = h.w.Write(append(b, '\n'))

	return err
}

func (h *ZeroStyleJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Можно просто вернуть тот же handler, если не требуется хранить контекст
	return h
}

func (h *ZeroStyleJSONHandler) WithGroup(name string) slog.Handler {
	return h
}
