package server

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	h "github.com/nuigcompsoc/api/internal/helpers"
	log "github.com/sirupsen/logrus"
)

/*
 * Does any random things to the context we want done before reaching the endpoint func
 */
func (s *Server) MiscMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

/*
 * This middleware logs primarily the request path, method, response status and completion latency
 */
func (s *Server) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		logFields := log.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"ip":     c.ClientIP(),
		}

		c.Next()

		logFields["latency_ns"] = time.Since(start).Nanoseconds()
		logFields["status"] = c.Writer.Status()
		log.WithFields(logFields).Info("request")
	}
}

/*
 * This middleware prints a panic in JSON to the log and redirects
 * the user to an error page with a get parameter containing the error message
 */
func RecoveryMiddlware(c *gin.Context, recovered interface{}) {
	message, ok := recovered.(string)
	if !ok {
		h.RedirectWithError(c, errors.New("unknown error"))
		return
	}
	h.RedirectWithError(c, errors.New(message))
	return
}
