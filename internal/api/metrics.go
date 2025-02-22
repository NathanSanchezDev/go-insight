package api

import (
	"time"

	"github.com/NathanSanchezDev/go-insight/internal/models"
	"github.com/gin-gonic/gin"
)

type MetricsHandler struct {
	metrics []models.EndpointMetric
}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		metrics: make([]models.EndpointMetric, 0),
	}
}

func (h *MetricsHandler) CollectMetrics(c *gin.Context) {
	var metric models.EndpointMetric

	if err := c.ShouldBindJSON(&metric); err != nil {
		c.JSON(400, gin.H{
			"error":   "Invalid metric format",
			"details": err.Error(),
		})
		return
	}

	if metric.Timestamp.IsZero() {
		metric.Timestamp = time.Now()
	}

	h.metrics = append(h.metrics, metric)

	c.JSON(200, gin.H{
		"status":    "success",
		"metric_id": metric.RequestID,
	})
}
