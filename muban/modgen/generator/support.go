package generator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	modutils "github.com/NSObjects/go-template/muban/modgen/utils"
)

func (g *Generator) ensureSupportFiles() error {
	if err := g.ensureContextSupport(); err != nil {
		return err
	}
	if err := g.ensureRespSupport(); err != nil {
		return err
	}
	return nil
}

func (g *Generator) ensureContextSupport() error {
	utilsDir := filepath.Join(g.config.RepoRoot, "internal", "utils")

	has, err := packageHasSymbol(utilsDir, "BuildContext(")
	if err != nil {
		return fmt.Errorf("检查 utils 上下文支持函数失败: %w", err)
	}
	if has {
		return nil
	}

	target := filepath.Join(utilsDir, "context_trace.go")
	if _, err := os.Stat(target); err == nil {
		return nil
	}

	content := `package utils

import (
    "context"
    "time"

    "github.com/labstack/echo/v4"
)

type TraceContext struct {
    TraceID   string
    SpanID    string
    RequestID string
    UserID    string
    StartTime time.Time
}

func ExtractTraceContext(c echo.Context) *TraceContext {
    tc := &TraceContext{
        StartTime: time.Now(),
    }

    if requestID := c.Request().Header.Get("X-Request-ID"); requestID != "" {
        tc.RequestID = requestID
    } else if requestID := c.Response().Header().Get("X-Request-ID"); requestID != "" {
        tc.RequestID = requestID
    }

    if traceID := c.Request().Header.Get("X-Trace-ID"); traceID != "" {
        tc.TraceID = traceID
    }

    if spanID := c.Request().Header.Get("X-Span-ID"); spanID != "" {
        tc.SpanID = spanID
    }

    if userID := c.Get("user_id"); userID != nil {
        if uid, ok := userID.(string); ok {
            tc.UserID = uid
        }
    }

    return tc
}

func BuildContext(c echo.Context) context.Context {
    ctx := c.Request().Context()
    tc := ExtractTraceContext(c)

    ctx = context.WithValue(ctx, "trace_id", tc.TraceID)
    ctx = context.WithValue(ctx, "span_id", tc.SpanID)
    ctx = context.WithValue(ctx, "request_id", tc.RequestID)
    ctx = context.WithValue(ctx, "user_id", tc.UserID)
    ctx = context.WithValue(ctx, "start_time", tc.StartTime)

    return ctx
}

func GetTraceID(ctx context.Context) string {
    if traceID, ok := ctx.Value("trace_id").(string); ok {
        return traceID
    }
    return ""
}

func GetRequestID(ctx context.Context) string {
    if requestID, ok := ctx.Value("request_id").(string); ok {
        return requestID
    }
    return ""
}

func GetUserID(ctx context.Context) string {
    if userID, ok := ctx.Value("user_id").(string); ok {
        return userID
    }
    return ""
}

func GetStartTime(ctx context.Context) time.Time {
    if startTime, ok := ctx.Value("start_time").(time.Time); ok {
        return startTime
    }
    return time.Now()
}

func WithTraceInfo(ctx context.Context, traceID, spanID, requestID, userID string) context.Context {
    ctx = context.WithValue(ctx, "trace_id", traceID)
    ctx = context.WithValue(ctx, "span_id", spanID)
    ctx = context.WithValue(ctx, "request_id", requestID)
    ctx = context.WithValue(ctx, "user_id", userID)
    ctx = context.WithValue(ctx, "start_time", time.Now())
    return ctx
}
`

	modutils.MustWrite(target, content, false)
	return nil
}

func (g *Generator) ensureRespSupport() error {
	respDir := filepath.Join(g.config.RepoRoot, "internal", "resp")

	hasOne, err := packageHasSymbol(respDir, "OneDataResponse(")
	if err != nil {
		return fmt.Errorf("检查 resp 数据响应函数失败: %w", err)
	}
	hasOperate, err := packageHasSymbol(respDir, "OperateSuccess(")
	if err != nil {
		return fmt.Errorf("检查 resp 操作响应函数失败: %w", err)
	}
	if hasOne && hasOperate {
		return nil
	}

	target := filepath.Join(respDir, "response_helpers.go")
	if _, err := os.Stat(target); err == nil {
		return nil
	}

	content := `package resp

import (
    "net/http"

    "github.com/labstack/echo/v4"
)

type operateSuccessBody struct {
    Code int    ` + "`json:\"code\"`" + `
    Msg  string ` + "`json:\"msg\"`" + `
}

type singleDataBody struct {
    Code int         ` + "`json:\"code\"`" + `
    Msg  string      ` + "`json:\"msg\"`" + `
    Data interface{} ` + "`json:\"data\"`" + `
}

func OperateSuccess(c echo.Context) error {
    payload := operateSuccessBody{
        Code: http.StatusOK,
        Msg:  "success",
    }
    return c.JSON(http.StatusOK, payload)
}

func OneDataResponse(data interface{}, c echo.Context) error {
    payload := singleDataBody{
        Code: http.StatusOK,
        Msg:  "success",
        Data: data,
    }
    return c.JSON(http.StatusOK, payload)
}
`

	modutils.MustWrite(target, content, false)
	return nil
}

func packageHasSymbol(dir, symbol string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return false, err
		}
		if strings.Contains(string(data), symbol) {
			return true, nil
		}
	}

	return false, nil
}
