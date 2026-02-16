package payment

import (
	"github.com/go-kit/kit/log"
	"time"
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

func (mw loggingMiddleware) Authorise(amount float32) (auth Authorisation, err error) {
	defer func(begin time.Time) {
		lvs := []interface{}{
			"method", "Authorise",
			"amount", amount,
			"result", auth.Authorised,
			"message", auth.Message,
			"took", time.Since(begin),
		}
		if err != nil {
			lvs = append(lvs, "err", err)
		}
		mw.logger.Log(lvs...)
	}(time.Now())
	return mw.next.Authorise(amount)
}

func (mw loggingMiddleware) Health() (health []Health) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Health",
			"result", len(health),
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Health()
}
