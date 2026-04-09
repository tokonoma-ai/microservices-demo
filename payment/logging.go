package payment

import (
	"context"
	"time"

	"github.com/go-kit/log"
)

// LoggingMiddleware logs method calls, parameters, results, and elapsed time.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) Authorise(ctx context.Context, amount float32) (auth Authorisation, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(append([]interface{}{
			"method", "Authorise",
			"result", auth.Authorised,
			"took", time.Since(begin),
		}, TraceLogKV(ctx)...)...)
	}(time.Now())
	return mw.next.Authorise(ctx, amount)
}

func (mw loggingMiddleware) Health(ctx context.Context) (health []Health) {
	defer func(begin time.Time) {
		mw.logger.Log(append([]interface{}{
			"method", "Health",
			"result", len(health),
			"took", time.Since(begin),
		}, TraceLogKV(ctx)...)...)
	}(time.Now())
	return mw.next.Health(ctx)
}
