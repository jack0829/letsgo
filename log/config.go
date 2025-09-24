package log

import (
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	DefaultConfig = Config{
		Level:      zapcore.InfoLevel.String(),
		Dir:        "./logs",
		FileName:   []string{"2006", "01", "02.log"},
		TimeFormat: time.RFC3339,
	}
)

type Config struct {
	Level      string
	Dir        string
	FileName   []string
	TimeFormat string
}
