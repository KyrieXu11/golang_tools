package log

import (
	rotate "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"path/filepath"
	"time"
)

var (
	levelMap = map[string]int{
		"debug": 1,
		"info":  2,
		"warn":  3,
		"error": 4,
		"panic": 5,
		"fatal": 6,
	}

	levelSlice = []string{"debug", "info", "warn", "error", "panic", "fatal"}

	maxLevel = getLevelNum(levelSlice[len(levelSlice)-1])
)

func getLevelNum(level string) int {
	return levelMap[level]
}

func getLevelString(level int) string {
	if level < 1 || level > len(levelSlice) {
		return "debug"
	}
	return levelSlice[level-1]
}

func getLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func getEncoder() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		LevelKey:       "level",
		TimeKey:        "time",
		MessageKey:     "message",
		NameKey:        "N",
		CallerKey:      "caller",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     logTime,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func logTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func getFileWriter(path, name string, rotationTime, rotationCount uint) io.Writer {
	timePattern := time.Now().Format("2006010215")
	writer, err := rotate.New(
		filepath.Join(path, name+timePattern+".log"),
		rotate.WithRotationTime(time.Duration(rotationTime)*time.Hour), // 日志切割时间间隔
		rotate.WithRotationCount(rotationCount),                        // 文件最大保存份数
	)
	if err != nil {
		panic(err)
	}
	return writer
}
