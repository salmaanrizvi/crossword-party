package config

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

// ConfigureLogger initializes and returns a logging client
func ConfigureLogger(cfg *Config) *zap.Config {
	var zapCfg zap.Config
	if cfg.Env == Production {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}

	zapCfg.EncoderConfig.TimeKey = "timestamp"
	zapCfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	err := zapCfg.Level.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		fmt.Println("Error setting log level", err)
		zapCfg.Level.SetLevel(zapcore.InfoLevel)
	}

	l, _ := zapCfg.Build()
	logger = l.Sugar().Named(fmt.Sprintf("%s-%s", "crossword-party", cfg.AppVersion))
	return &zapCfg
}

// Logger is the only acceptable way to access the logger
func Logger() *zap.SugaredLogger {
	return logger
}

// GetGinLoggerMiddleware returns a gin.HandlerFunc (middleware) that logs requests using zap
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
func GetGinLoggerMiddleware(logger *zap.SugaredLogger) func(*gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		latency := time.Now().Sub(start)
		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			logger.Infow(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("latency", fmt.Sprintf("%v", latency)),
			)
		}
	}
}

// GetRecoveryLoggerMiddleware returns a gin.HandlerFunc (middleware)
// that recovers from any panics and logs requests using zap
func GetRecoveryLoggerMiddleware(logger *zap.SugaredLogger, printStackTrace bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Errorw(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if printStackTrace {
					logger.Errorw("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Errorw("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}
