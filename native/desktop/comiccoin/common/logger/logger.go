package logger

import (
	"log/slog"
	"os"
)

// NewProvider creates a new logger instance with a configurable logging level.
// The logger is set to log to the standard output and includes source file information.
func NewProvider() *slog.Logger {
	// Create a logging level variable to control the verbosity of the logger.
	// The level is set to Info by default.
	var loggingLevel = new(slog.LevelVar)

	// Create a new logger instance with the logging level variable.
	// The logger is set to log to the standard output and includes source file information.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// AddSource is set to true to include the source file information in the log output.
		AddSource: true,
		// The logging level is set to the loggingLevel variable, allowing it to be changed dynamically.
		Level: loggingLevel,
	}))

	// Set the logging level to Debug to include all log messages.
	// This can be changed later to a different level (e.g. Info, Warn, Error) to filter out less important messages.
	loggingLevel.Set(slog.LevelDebug)

	// // Set the logger as the default logger for the application.
	// // This is commented out to allow for a custom logger to be used instead.
	// slog.SetDefault(logger)

	return logger
}
