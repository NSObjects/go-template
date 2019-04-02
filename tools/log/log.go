/*
 *
 * log.go
 * log
 *
 * Created by lin on 2018/12/10 3:18 PM
 * Copyright © 2017-2018 PYL. All rights reserved.
 *
 */

package log

import (
	"fmt"
	"go-template/configs"
	"net/http"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap/zapcore"

	"time"

	"github.com/go-xorm/core"
	"go.uber.org/zap"
)

var logger *Log

func LumberJackLogger(filePath string, maxSize int, maxBackups int, maxAge int) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    maxSize, // megabytes
		MaxBackups: maxBackups,
		MaxAge:     maxAge, //days
	}
}

func Init() {
	l := zap.New(loggerCore()).WithOptions(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	logger = NewLog(l)
}

type LogInterface interface {
	Debug(msg interface{}, fields ...zap.Field)
	Debugf(msg string, args ...interface{})
	Info(msg interface{}, fields ...zap.Field)
	Infof(msg string, args ...interface{})
	Warn(msg string, fields ...zap.Field)
	Warnf(msg string, args ...interface{})
	Error(err interface{}, fields ...zap.Field)
	ErrorSkip(skip int, err interface{}, fields ...zap.Field)
	Errorf(msg string, args ...interface{})
	Fatal(msg interface{}, fields ...zap.Field)
	Fatalf(msg string, args ...interface{})
	Panic(msg interface{}, fields ...zap.Field)
	Panicf(msg string, args ...interface{})
	DebugResponse(response *http.Response)
}

type Log struct {
	logger *zap.Logger
}

func NewLog(logger *zap.Logger) *Log {
	return &Log{
		logger: logger,
	}
}

func (l *Log) Debug(msg interface{}, fields ...zap.Field) {
	l.logger.Debug(fmt.Sprint(msg), fields...)
}

func (l *Log) Debugf(msg string, args ...interface{}) {
	logger.Debug(fmt.Sprintf(msg, args...))
}

func (l *Log) Info(msg interface{}, fields ...zap.Field) {
	l.logger.Info(fmt.Sprint(msg), fields...)
}

func (l *Log) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}

func (l *Log) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Log) Warnf(msg string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(msg, args...))
}

func (l *Log) Error(err interface{}, fields ...zap.Field) {
	l.logger.Error(fmt.Sprint(err), fields...)
}

func (*Log) ErrorSkip(skip int, err interface{}, fields ...zap.Field) {
	l := zap.New(loggerCore()).WithOptions(
		zap.AddCaller(),
		zap.AddCallerSkip(skip),
	)
	l.Error(fmt.Sprint(err),
		fields...)
	l.Sync()
}

func (l *Log) Errorf(msg string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(msg, args...))
}

func (l *Log) Fatal(msg interface{}, fields ...zap.Field) {
	l.logger.Fatal(fmt.Sprint(msg),
		fields...)
}

func (l *Log) Fatalf(msg string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(msg, args...))
}

func (l *Log) Panic(msg interface{}, fields ...zap.Field) {
	l.logger.Panic(fmt.Sprint(msg), fields...)
}

func (l *Log) Panicf(msg string, args ...interface{}) {
	l.logger.Panic(fmt.Sprintf(msg, args...))
}

func (l *Log) DebugResponse(response *http.Response) {
	bodyBuffer := make([]byte, 5000)
	var str string
	count, err := response.Body.Read(bodyBuffer)
	for ; count > 0; count, err = response.Body.Read(bodyBuffer) {
		if err != nil {
		}
		str += string(bodyBuffer[:count])
	}
	Debugf("response data : %v", str)
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func logPath(base string, level zapcore.Level) string {
	path := fmt.Sprintf("%s/%s.log", base, level.String())
	return path
}

func loggerCore() (core zapcore.Core) {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = TimeEncoder

	switch configs.System.Level {
	case configs.DebugLevel:
		consoleEncoder := zapcore.NewConsoleEncoder(config)
		consoleDebugging := zapcore.Lock(os.Stdout)
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleDebugging, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return true
			})),
		)
	case configs.PrdocutionLevel:
		topicDebugging := zapcore.AddSync(LumberJackLogger(logPath(configs.Log.Path, zapcore.DebugLevel), configs.Log.MaxSize, configs.Log.MaxBackups, configs.Log.MaxAge))
		topicErrors := zapcore.AddSync(LumberJackLogger(logPath(configs.Log.Path, zapcore.ErrorLevel), configs.Log.MaxSize, configs.Log.MaxBackups, configs.Log.MaxAge))
		jsonEncoder := zapcore.NewJSONEncoder(config)
		core = zapcore.NewTee(
			zapcore.NewCore(jsonEncoder, topicErrors, highPriority),
			zapcore.NewCore(jsonEncoder, topicDebugging, lowPriority),
		)
	}

	return core
}

func Debug(msg interface{}, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Debugf(msg string, args ...interface{}) {
	logger.Debugf(msg, args...)
}

func Info(msg interface{}, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Infof(msg string, args ...interface{}) {
	logger.Infof(msg, args...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Warnf(msg string, args ...interface{}) {
	logger.Warnf(msg, args...)
}

func Error(err interface{}, fields ...zap.Field) {
	logger.Error(err, fields...)
}

func ErrorSkip(skip int, err interface{}, fields ...zap.Field) {
	logger.ErrorSkip(skip, err, fields...)
}

func Errorf(msg string, args ...interface{}) {
	logger.Errorf(msg, args...)
}

func Fatal(msg interface{}, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Fatalf(msg string, args ...interface{}) {
	logger.Fatalf(msg, args...)
}

func Panic(msg interface{}, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

func Panicf(msg string, args ...interface{}) {
	logger.Panicf(msg, args...)
}

func DebugResponse(response *http.Response) {
	logger.DebugResponse(response)
}

type Xormlog struct {
	showSQL bool
	level   core.LogLevel
}

func (c *Xormlog) Debug(v ...interface{}) {
	logger.Debug(fmt.Sprint(v...))
}

func (Xormlog) Debugf(format string, v ...interface{}) {
	logger.Debug(fmt.Sprint(v...))
}

func (c *Xormlog) Error(v ...interface{}) {
	logger.Error(fmt.Sprint(v...))
}

func (c *Xormlog) Errorf(format string, v ...interface{}) {
	logger.Error(fmt.Sprintf(format, v...))
}

func (c *Xormlog) Info(v ...interface{}) {
	logger.Info(fmt.Sprint(v...))
}

func (c *Xormlog) Infof(format string, v ...interface{}) {
	logger.Info(fmt.Sprintf(format, v...))
}

func (c *Xormlog) Warn(v ...interface{}) {
	logger.Warn(fmt.Sprint(v...))
}

func (c *Xormlog) Warnf(format string, v ...interface{}) {
	logger.Warn(fmt.Sprintf(format, v...))
}

func (c *Xormlog) Level() core.LogLevel {
	return c.level
}

func (c *Xormlog) SetLevel(l core.LogLevel) {
	c.level = l
}

func (c *Xormlog) ShowSQL(show ...bool) {
	if len(show) == 0 {
		c.showSQL = true
	} else {
		c.showSQL = show[0]
	}
}

func (c *Xormlog) IsShowSQL() bool {
	return c.showSQL
}
