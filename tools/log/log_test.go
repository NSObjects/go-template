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
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/go-xorm/core"
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
			Init()
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

func Test_loggerCore(t *testing.T) {
	tests := []struct {
		name     string
		wantCore zapcore.Core
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCore := loggerCore(); !reflect.DeepEqual(gotCore, tt.wantCore) {
				t.Errorf("loggerCore() = %v, want %v", gotCore, tt.wantCore)
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

func TestXormlog_Debug(t *testing.T) {
	type fields struct {
		showSQL bool
		level   core.LogLevel
	}
	type args struct {
		v []interface{}
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
			c := &Xormlog{
				showSQL: tt.fields.showSQL,
				level:   tt.fields.level,
			}
			c.Debug(tt.args.v...)
		})
	}
}

func TestXormlog_Debugf(t *testing.T) {
	type fields struct {
		showSQL bool
		level   core.LogLevel
	}
	type args struct {
		format string
		v      []interface{}
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
			x := Xormlog{
				showSQL: tt.fields.showSQL,
				level:   tt.fields.level,
			}
			x.Debugf(tt.args.format, tt.args.v...)
		})
	}
}

func TestXormlog_Error(t *testing.T) {
	type fields struct {
		showSQL bool
		level   core.LogLevel
	}
	type args struct {
		v []interface{}
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
			c := &Xormlog{
				showSQL: tt.fields.showSQL,
				level:   tt.fields.level,
			}
			c.Error(tt.args.v...)
		})
	}
}

func TestXormlog_Errorf(t *testing.T) {
	type fields struct {
		showSQL bool
		level   core.LogLevel
	}
	type args struct {
		format string
		v      []interface{}
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
			c := &Xormlog{
				showSQL: tt.fields.showSQL,
				level:   tt.fields.level,
			}
			c.Errorf(tt.args.format, tt.args.v...)
		})
	}
}

func TestXormlog_Info(t *testing.T) {
	type fields struct {
		showSQL bool
		level   core.LogLevel
	}
	type args struct {
		v []interface{}
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
			c := &Xormlog{
				showSQL: tt.fields.showSQL,
				level:   tt.fields.level,
			}
			c.Info(tt.args.v...)
		})
	}
}

func TestXormlog_Infof(t *testing.T) {
	type fields struct {
		showSQL bool
		level   core.LogLevel
	}
	type args struct {
		format string
		v      []interface{}
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
			c := &Xormlog{
				showSQL: tt.fields.showSQL,
				level:   tt.fields.level,
			}
			c.Infof(tt.args.format, tt.args.v...)
		})
	}
}

func TestXormlog_Warn(t *testing.T) {
	type fields struct {
		showSQL bool
		level   core.LogLevel
	}
	type args struct {
		v []interface{}
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
			c := &Xormlog{
				showSQL: tt.fields.showSQL,
				level:   tt.fields.level,
			}
			c.Warn(tt.args.v...)
		})
	}
}

func TestXormlog_Warnf(t *testing.T) {
	type fields struct {
		showSQL bool
		level   core.LogLevel
	}
	type args struct {
		format string
		v      []interface{}
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
			c := &Xormlog{
				showSQL: tt.fields.showSQL,
				level:   tt.fields.level,
			}
			c.Warnf(tt.args.format, tt.args.v...)
		})
	}
}

func TestXormlog_Level(t *testing.T) {
	type fields struct {
		showSQL bool
		level   core.LogLevel
	}
	tests := []struct {
		name   string
		fields fields
		want   core.LogLevel
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Xormlog{
				showSQL: tt.fields.showSQL,
				level:   tt.fields.level,
			}
			if got := c.Level(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Xormlog.Level() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXormlog_SetLevel(t *testing.T) {
	type fields struct {
		showSQL bool
		level   core.LogLevel
	}
	type args struct {
		l core.LogLevel
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
			c := &Xormlog{
				showSQL: tt.fields.showSQL,
				level:   tt.fields.level,
			}
			c.SetLevel(tt.args.l)
		})
	}
}

func TestXormlog_ShowSQL(t *testing.T) {
	type fields struct {
		showSQL bool
		level   core.LogLevel
	}
	type args struct {
		show []bool
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
			c := &Xormlog{
				showSQL: tt.fields.showSQL,
				level:   tt.fields.level,
			}
			c.ShowSQL(tt.args.show...)
		})
	}
}

func TestXormlog_IsShowSQL(t *testing.T) {
	type fields struct {
		showSQL bool
		level   core.LogLevel
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Xormlog{
				showSQL: tt.fields.showSQL,
				level:   tt.fields.level,
			}
			if got := c.IsShowSQL(); got != tt.want {
				t.Errorf("Xormlog.IsShowSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
