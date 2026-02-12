package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(authHandler *AuthHandler, goodsHandler *GoodsHandler, logger *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(zapLogger(logger))

	router.LoadHTMLGlob("web/templates/*.html")
	router.Static("/static", "web/static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	router.GET("/main", func(c *gin.Context) {
		c.HTML(http.StatusOK, "main.html", nil)
	})

	api := router.Group("/api")
	api.POST("/login", authHandler.Login)
	api.POST("/register", authHandler.Register)
	api.GET("/goods", goodsHandler.List)
	api.POST("/goods", goodsHandler.Add)
	api.DELETE("/goods/:id", goodsHandler.Delete)
	api.PUT("/goods/stock", goodsHandler.UpdateStock)
	api.POST("/goods/undo", goodsHandler.Undo)

	return router
}

func zapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		if raw != "" {
			path = path + "?" + raw
		}
		logger.Info("request",
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
