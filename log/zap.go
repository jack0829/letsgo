package log

import (
	ZAP "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

var (
	zap         *ZAP.Logger
	sugar       *ZAP.SugaredLogger
	writeSyncer zapcore.WriteSyncer
)

func W() io.Writer {
	return writeSyncer
}

func Z() *ZAP.Logger {
	return zap
}

func L() *ZAP.SugaredLogger {
	return sugar
}

func Debugf(tpl string, args ...interface{}) {
	sugar.Debugf(tpl, args...)
}

func Infof(tpl string, args ...interface{}) {
	sugar.Infof(tpl, args...)
}

func Warnf(tpl string, args ...interface{}) {
	sugar.Warnf(tpl, args...)
}

func Errorf(tpl string, args ...interface{}) {
	sugar.Errorf(tpl, args...)
}

func DPanicf(tpl string, args ...interface{}) {
	sugar.DPanicf(tpl, args...)
}

func Panicf(tpl string, args ...interface{}) {
	sugar.Panicf(tpl, args...)
}

func Fatalf(tpl string, args ...interface{}) {
	sugar.Fatalf(tpl, args...)
}

func Init(opts ...Option) (err error) {

	c := DefaultConfig
	for _, opt := range opts {
		opt(&c)
	}

	level, err := zapcore.ParseLevel(c.Level)
	if err != nil {
		return
	}

	logger, err := NewLogger(c.Dir, c.FileName...)
	if err != nil {
		return
	}

	writeSyncer = zapcore.AddSync(logger)

	encoder := ZAP.NewProductionEncoderConfig()
	if c.TimeFormat != "" {
		encoder.EncodeTime = zapcore.TimeEncoderOfLayout(c.TimeFormat)
	}

	zap = ZAP.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoder),
			writeSyncer,
			level,
		),
		ZAP.AddCaller(),
	)

	sugar = zap.Sugar()
	return
}
