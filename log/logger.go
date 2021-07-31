package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger struct {
	conf   *Config
	logger *zap.SugaredLogger
}

type Config struct {
	logLevel    string          // 默认日志记录级别
	stackTrace  string          // 记录堆栈的级别
	atomicLevel zap.AtomicLevel // 用于动态更改日志记录级别
	projectName string          // 项目名称
	callerSkip  int             // CallerSkip次数
	jsonFormat  bool            // 输出json格式
	consoleOut  bool            // 是否输出到console
	fileOut     *fileOut
}

type fileOut struct {
	enable        bool   // 是否将日志输出到文件
	path          string // 日志保存路径
	name          string // 日志保存的名称，不些随机生成
	rotationTime  uint   // 日志切割时间间隔(小时)
	rotationCount uint   // 文件最大保存份数
}

func defaultConfig() *Config {
	return &Config{
		logLevel:    "info",
		stackTrace:  "panic",
		atomicLevel: zap.NewAtomicLevel(),
		projectName: "",
		callerSkip:  1,
		jsonFormat:  true,
		consoleOut:  true,
		fileOut: &fileOut{
			enable:        true,
			path:          "logs",
			name:          "log",
			rotationTime:  24,
			rotationCount: 7,
		},
	}
}

func NewLogger() *Logger {
	l := &Logger{
		conf: defaultConfig(),
	}
	l.setDefaultConf()
	return l
}

func (l *Logger) setConf(conf *Config) {
	var cores []zapcore.Core

	var encoder zapcore.Encoder

	if conf.jsonFormat {
		encoder = zapcore.NewJSONEncoder(getEncoder())
	} else {
		encoder = zapcore.NewConsoleEncoder(getEncoder())
	}

	conf.atomicLevel.SetLevel(getLevel(conf.logLevel))

	if conf.consoleOut {
		writer := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(encoder, writer, conf.atomicLevel)
		cores = append(cores, core)
	}

	if conf.fileOut.enable {
		for i := 1; i <= maxLevel; i++ {
			level := getLevelString(i)
			fileWriter := getFileWriter(
				conf.fileOut.path,
				fmt.Sprintf("%s-%s-", conf.fileOut.name, level),
				conf.fileOut.rotationTime,
				conf.fileOut.rotationCount,
			)
			writer := zapcore.AddSync(fileWriter)
			core := zapcore.NewCore(encoder, writer, getLevel(level))
			cores = append(cores, core)
		}
	}

	combinedCore := zapcore.NewTee(cores...)

	logger := zap.New(
		combinedCore,
		zap.AddCallerSkip(conf.callerSkip),
		zap.AddStacktrace(getLevel(conf.stackTrace)),
		zap.AddCaller(),
	)

	if conf.projectName != "" {
		logger = logger.Named(conf.projectName)
	}

	defer l.Sync()

	l.logger = logger.Sugar()
}

func (l *Logger) setDefaultConf() {
	conf := l.conf
	l.setConf(conf)
}

func (l *Logger) SetConfig(conf *Config) {
	l.conf = conf
	l.setConf(conf)
}

func (l *Logger) Sync() {
	l.logger.Sync()
}

func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.logger.Panicf(template, args...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.logger.Panicw(msg, keysAndValues...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}
