package server

import (
	"github.com/gin-gonic/gin"
)

// Returns the routes associated with /v1
func (s *Server) v1Router(r *gin.RouterGroup) {
	r.Use(s.MiscMiddleware())
	r.Use(s.LoggingMiddleware())
	r.GET("/", s.RootGet)
	r.GET("ping", s.MiscV1PingGet)
	r.GET("brew", s.MiscV1BrewGet)
	r.GET("events", s.EventsV1Get)

	// EVENTS route
	e := r.Group("/events")
	e.GET("upcoming", s.EventsV1UpcomingGet)
	e.GET("upcoming/:id", s.EventsV1UpcomingSocIDGet)
	e.GET("past", s.EventsV1PastGet)
	e.GET("past/:id", s.EventsV1PastSocIDGet)
}

// SetupRouter function will perform all route operations
func SetupRouter() *gin.Engine {

	// We are now relying on our own logging middleware to log all paths accessed to stdout
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.CustomRecovery(RecoveryMiddlware))

	r.Use(func(c *gin.Context) {
		// add header Access-Control-Allow-Origin
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	})

	return r
}
