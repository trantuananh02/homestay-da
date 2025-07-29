package logger

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"os"
	"time"
)

func Trace() log.Valuer {
	return func(ctx context.Context) interface{} {
		s := ctx.Value("trace_id").(string)
		return s
	}
}

func NewContextLog(ctx context.Context) *log.Helper {
	logger := log.With(log.NewStdLogger(os.Stdout),
		"trace", Trace(),
	)
	logger = log.With(logger, "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)
	return log.NewHelper(logger).WithContext(ctx)
}

func GenerateTraceID(serviceName string) string {
	now := time.Now()
	traceID := fmt.Sprintf("%s-%d", serviceName, now.UnixNano())
	return traceID
}
