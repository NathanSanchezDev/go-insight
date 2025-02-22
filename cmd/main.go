package main

import (
	"github.com/NathanSanchezDev/go-insight/internal/api"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	metricsHandler := api.NewMetricsHandler()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/collect", metricsHandler.CollectMetrics)

	r.Run(":8080")
}
