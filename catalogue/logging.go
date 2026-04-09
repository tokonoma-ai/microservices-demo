package catalogue

import (
	"context"
	"strings"
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

func (mw loggingMiddleware) List(ctx context.Context, tags []string, order string, pageNum, pageSize int) (socks []Sock, err error) {
	defer func(begin time.Time) {
		lvs := []interface{}{
			"method", "List",
			"tags", strings.Join(tags, ", "),
			"order", order,
			"pageNum", pageNum,
			"pageSize", pageSize,
			"result", len(socks),
			"took", time.Since(begin),
		}
		if err != nil {
			lvs = append(lvs, "err", err)
		}
		lvs = append(lvs, TraceLogKV(ctx)...)
		mw.logger.Log(lvs...)
	}(time.Now())
	return mw.next.List(ctx, tags, order, pageNum, pageSize)
}

func (mw loggingMiddleware) Count(ctx context.Context, tags []string) (n int, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(append([]interface{}{
			"method", "Count",
			"tags", strings.Join(tags, ", "),
			"result", n,
			"err", err,
			"took", time.Since(begin),
		}, TraceLogKV(ctx)...)...)
	}(time.Now())
	return mw.next.Count(ctx, tags)
}

func (mw loggingMiddleware) Get(ctx context.Context, id string) (s Sock, err error) {
	defer func(begin time.Time) {
		lvs := []interface{}{
			"method", "Get",
			"id", id,
			"sock", s.ID,
			"took", time.Since(begin),
		}
		if err != nil {
			lvs = append(lvs, "err", err)
		}
		lvs = append(lvs, TraceLogKV(ctx)...)
		mw.logger.Log(lvs...)
	}(time.Now())
	return mw.next.Get(ctx, id)
}

func (mw loggingMiddleware) Tags(ctx context.Context) (tags []string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(append([]interface{}{
			"method", "Tags",
			"result", len(tags),
			"err", err,
			"took", time.Since(begin),
		}, TraceLogKV(ctx)...)...)
	}(time.Now())
	return mw.next.Tags(ctx)
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
