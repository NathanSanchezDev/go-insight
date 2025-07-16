package observability

import (
	"database/sql"
	"log"
	"time"

	"github.com/NathanSanchezDev/go-insight/core/db"
	"github.com/NathanSanchezDev/go-insight/core/models"
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
	now := time.Now()

	trace.EndTime = sql.NullTime{
		Time:  now,
		Valid: true,
	}

	trace.Duration = sql.NullFloat64{
		Float64: now.Sub(trace.StartTime).Seconds() * 1000,
		Valid:   true,
	}

	err := db.UpdateTrace(trace)
	if err != nil {
		log.Println("Error updating trace:", err)
	}
}
