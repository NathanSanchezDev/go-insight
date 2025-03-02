package main

import (
	"fmt"
	"log"

	"github.com/NathanSanchezDev/go-insight/internal/api"
	"github.com/NathanSanchezDev/go-insight/internal/db"
)

func main() {
	db.InitDB()

	api.InsertLog("auth-service", "INFO", "User login successful")

	logs, err := api.GetLogs()
	if err != nil {
		log.Fatal("❌ Failed to fetch logs:", err)
	}

	fmt.Println("\n📜 Retrieved Logs:")
	for _, logEntry := range logs {
		fmt.Printf("[%s] %s - %s: %s\n", logEntry.Timestamp, logEntry.ServiceName, logEntry.LogLevel, logEntry.Message)
	}
}
