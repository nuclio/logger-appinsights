package main

import (
	"fmt"
	"time"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/Microsoft/ApplicationInsights-Go/appinsights/contracts"
)

type AppInsightsLogger struct {
	client *appinsights.TelemetryClient
	name   string
}

func (logger *AppInsightsLogger) Close() {
	select {
	case <-(*logger.client).Channel().Close(10 * time.Second):
	case <-time.After(30 * time.Second):
	}
}

func toString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

func (logger *AppInsightsLogger) emitUnstructured(severity contracts.SeverityLevel, format interface{}, vars ...interface{}) {
	message := fmt.Sprintf(toString(format), vars...)
	trace := appinsights.NewTraceTelemetry(message, severity)
	(*logger.client).Track(trace)
}

func (logger *AppInsightsLogger) emitStructured(severity contracts.SeverityLevel, message interface{}, vars ...interface{}) {
	trace := appinsights.NewTraceTelemetry(toString(message), severity)
	// set properties
	for i := 0; i < len(vars); i++ {
		key := toString(vars[i])
		i++
		value := toString(vars[i])

		trace.Properties[key] = value
	}
	(*logger.client).Track(trace)
}

// implemenet https://github.com/nuclio/logger/blob/master/logger.go interface

func (logger *AppInsightsLogger) Error(format interface{}, vars ...interface{}) {
	logger.emitUnstructured(appinsights.Error, format, vars...)
}

func (logger *AppInsightsLogger) Warn(format interface{}, vars ...interface{}) {
	logger.emitUnstructured(appinsights.Warning, format, vars...)
}

func (logger *AppInsightsLogger) Info(format interface{}, vars ...interface{}) {
	logger.emitUnstructured(appinsights.Information, format, vars...)
}

func (logger *AppInsightsLogger) Debug(format interface{}, vars ...interface{}) {
	// debug will use the *Verbose* severity level
	logger.emitUnstructured(appinsights.Verbose, format, vars...)
}

func (logger *AppInsightsLogger) ErrorWith(format interface{}, vars ...interface{}) {
	logger.emitStructured(appinsights.Error, format, vars...)
}

func (logger *AppInsightsLogger) WarnWith(format interface{}, vars ...interface{}) {
	logger.emitStructured(appinsights.Warning, format, vars...)
}

func (logger *AppInsightsLogger) InfoWith(format interface{}, vars ...interface{}) {
	logger.emitStructured(appinsights.Information, format, vars...)
}

func (logger *AppInsightsLogger) DebugWith(format interface{}, vars ...interface{}) {
	logger.emitStructured(appinsights.Verbose, format, vars...)
}

// Flush flushes buffered logs
func (logger *AppInsightsLogger) Flush() {
	(*logger.client).Channel().Flush()
}

// GetChild returns a child logger
func (logger *AppInsightsLogger) GetChild(name string) AppInsightsLogger {
	return NewLogger(logger.client, fmt.Sprintf("%s.%s", logger.name, name))
}

func NewLogger(client *appinsights.TelemetryClient, name string) AppInsightsLogger {
	return AppInsightsLogger{client, name}
}

var logger AppInsightsLogger

func main() {

	client := appinsights.NewTelemetryClient("<app insights instrumentatoin key>")
	logger = NewLogger(&client, "RootLogger")

	logger.Error("Error message without properties")
	logger.ErrorWith("Error message with properties", "property1", "value", "property2", 100)

	childLogger := logger.GetChild("ChildLogger")
	childLogger.Info("Info Message without properties")
	childLogger.InfoWith("Info Message with properties", "logger", childLogger.name)

	time.Sleep(30 * time.Second)

	logger.Close()
}
