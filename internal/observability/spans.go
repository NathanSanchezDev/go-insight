package observability

import (
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/NathanSanchezDev/go-insight/internal/db"
	"github.com/NathanSanchezDev/go-insight/internal/models"
)

func StartSpan(traceID, parentID, service, operation string) models.Span {
	span := models.Span{
		ID:        generateUUID(),
		TraceID:   traceID,
		ParentID:  parentID,
		Service:   service,
		Operation: operation,
		StartTime: time.Now(),
	}

	err := db.StoreSpan(span)
	if err != nil {
		log.Println("Error storing span:", err)
	}

	return span
}

func EndSpan(span *models.Span) {
	span.EndTime = time.Now()
	span.Duration = span.EndTime.Sub(span.StartTime).Seconds() * 1000

	err := db.UpdateSpan(span)
	if err != nil {
		log.Println("Error updating span:", err)
	}
}

func generateUUID() string {
	return uuid.New().String()
}
