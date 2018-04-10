package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nuclio/logger-appinsights"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/nuclio/logger"
)

func main() {
	fmt.Println("Starting test")

	// create configuration
	telemetryClientConfig := appinsights.NewTelemetryConfiguration(os.Getenv("NUCLIO_APPINSIGHTS_INSTRUMENTATION_KEY"))
	telemetryClientConfig.MaxBatchInterval = 1 * time.Second
	telemetryClientConfig.MaxBatchSize = 1024

	// create a telemetry client
	telemetryClient := appinsights.NewTelemetryClientFromConfig(telemetryClientConfig)

	// create an appinsights logger with the client
	logger, _ := appinsightslogger.NewLogger(telemetryClient, "root", logger.LevelDebug)

	// output some stuff
	logger.Error("Error message without properties")
	logger.ErrorWith("Error message with properties",
		"property1", "value1",
		"property2", 100,
		"property3", "value3")

	// create a child logger
	childLogger := logger.GetChild("ChildLogger")
	childLogger.Info("Info Message without properties")
	childLogger.InfoWith("Info Message with properties", "logger", "")

	fmt.Println("Logs sent, closing")

	// close
	if err := logger.Close(); err != nil {
		fmt.Printf("Error closing: %s", err.Error())
	}

	fmt.Println("Logs sent, closed")
}
