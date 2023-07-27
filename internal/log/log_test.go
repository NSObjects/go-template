/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package log

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func TestLumberJackLogger(t *testing.T) {
	type args struct {
		filePath   string
		maxSize    int
		maxBackups int
		maxAge     int
	}
	tests := []struct {
		name string
		args args
		want *lumberjack.Logger
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LumberJackLogger(tt.args.filePath, tt.args.maxSize, tt.args.maxBackups, tt.args.maxAge); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LumberJackLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestNewLog(t *testing.T) {
	type args struct {
		logger *zap.Logger
	}
	tests := []struct {
		name string
		args args
		want *Log
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLog(tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLog_Debug(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		msg    interface{}
		fields []zap.Field
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
			}
			l.Debug(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestLog_Debugf(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
			}
			l.Debugf(tt.args.msg, tt.args.args...)
		})
	}
}

func TestLog_Info(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		msg    interface{}
		fields []zap.Field
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
			}
			l.Info(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestLog_Infof(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
			}
			l.Infof(tt.args.msg, tt.args.args...)
		})
	}
}

func TestLog_Warn(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		msg    string
		fields []zap.Field
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
			}
			l.Warn(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestLog_Warnf(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
			}
			l.Warnf(tt.args.msg, tt.args.args...)
		})
	}
}

func TestLog_Error(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		err    interface{}
		fields []zap.Field
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
			}
			l.Error(tt.args.err, tt.args.fields...)
		})
	}
}

func TestLog_ErrorSkip(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		skip   int
		err    interface{}
		fields []zap.Field
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
			}
			l.ErrorSkip(tt.args.skip, tt.args.err, tt.args.fields...)
		})
	}
}

func TestLog_Errorf(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
			}
			l.Errorf(tt.args.msg, tt.args.args...)
		})
	}
}

func TestLog_Fatal(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		msg    interface{}
		fields []zap.Field
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
			}
			l.Fatal(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestLog_Fatalf(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
			}
			l.Fatalf(tt.args.msg, tt.args.args...)
		})
	}
}

func TestLog_Panic(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		msg    interface{}
		fields []zap.Field
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
			}
			l.Panic(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestLog_Panicf(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
			}
			l.Panicf(tt.args.msg, tt.args.args...)
		})
	}
}

func TestLog_DebugResponse(t *testing.T) {
	type fields struct {
		logger *zap.Logger
	}
	type args struct {
		response *http.Response
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				logger: tt.fields.logger,
			}
			l.DebugResponse(tt.args.response)
		})
	}
}

func TestTimeEncoder(t *testing.T) {
	type args struct {
		t   time.Time
		enc zapcore.PrimitiveArrayEncoder
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TimeEncoder(tt.args.t, tt.args.enc)
		})
	}
}

func Test_logPath(t *testing.T) {
	type args struct {
		base  string
		level zapcore.Level
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logPath(tt.args.base, tt.args.level); got != tt.want {
				t.Errorf("logPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDebug(t *testing.T) {
	type args struct {
		msg    interface{}
		fields []zap.Field
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debug(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestDebugf(t *testing.T) {
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debugf(tt.args.msg, tt.args.args...)
		})
	}
}

func TestInfo(t *testing.T) {
	type args struct {
		msg    interface{}
		fields []zap.Field
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Info(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestInfof(t *testing.T) {
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Infof(tt.args.msg, tt.args.args...)
		})
	}
}

func TestWarn(t *testing.T) {
	type args struct {
		msg    string
		fields []zap.Field
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Warn(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestWarnf(t *testing.T) {
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Warnf(tt.args.msg, tt.args.args...)
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		err    interface{}
		fields []zap.Field
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Error(tt.args.err, tt.args.fields...)
		})
	}
}

func TestErrorSkip(t *testing.T) {
	type args struct {
		skip   int
		err    interface{}
		fields []zap.Field
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ErrorSkip(tt.args.skip, tt.args.err, tt.args.fields...)
		})
	}
}

func TestErrorf(t *testing.T) {
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Errorf(tt.args.msg, tt.args.args...)
		})
	}
}

func TestFatal(t *testing.T) {
	type args struct {
		msg    interface{}
		fields []zap.Field
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Fatal(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestFatalf(t *testing.T) {
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Fatalf(tt.args.msg, tt.args.args...)
		})
	}
}

func TestPanic(t *testing.T) {
	type args struct {
		msg    interface{}
		fields []zap.Field
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Panic(tt.args.msg, tt.args.fields...)
		})
	}
}

func TestPanicf(t *testing.T) {
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Panicf(tt.args.msg, tt.args.args...)
		})
	}
}

func TestDebugResponse(t *testing.T) {
	type args struct {
		response *http.Response
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DebugResponse(tt.args.response)
		})
	}
}
