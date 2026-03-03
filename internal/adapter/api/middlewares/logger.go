package middlewares

import (
	"time"

	"github.com/labstack/echo/v5"

	port "github.com/PodPloy/podploy/internal/domain/ports"
)

// LoggerMiddleware returns an Echo middleware that logs every HTTP request
// using the provided ILogger. Each log entry includes the request ID, method,
// URI, client IP, and response latency.
func LoggerMiddleware(log port.ILogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Logger().Error(err.Error())
			}

			req := c.Request()
			res := c.Response()

			reqID := req.Header.Get(echo.HeaderXRequestID)
			if reqID == "" {
				reqID = res.Header().Get(echo.HeaderXRequestID)
			}

			log.WithRequest(reqID, req.Method, req.RequestURI, "", c.RealIP(), time.Since(start))

			return err
		}
	}
}
