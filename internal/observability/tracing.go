package observability

import (
	"log"
	"time"

	"github.com/NathanSanchezDev/go-insight/internal/db"
	"github.com/NathanSanchezDev/go-insight/internal/models"
)

func StartTrace(serviceName string) models.Trace {
	trace := models.Trace{
		ID:          GenerateUUID(),
		ServiceName: serviceName,
		StartTime:   time.Now(),
	}

	err := db.StoreTrace(trace)
	if err != nil {
		log.Println("Error storing trace:", err)
	}

	return trace
}

func EndTrace(trace *models.Trace) {
	trace.EndTime = time.Now()
	trace.Duration = trace.EndTime.Sub(trace.StartTime).Seconds() * 1000

	err := db.UpdateTrace(trace)
	if err != nil {
		log.Println("Error updating trace:", err)
	}
}
