/*
Copyright 2018 The Nuclio Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
