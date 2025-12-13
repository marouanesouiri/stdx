/*
Package log provides a simple way to write logs in Go applications.

It has two types of loggers:
  - JSONLogger: Writes logs in JSON format. Good for servers.
  - TextLogger: Writes logs as text with colors. Good for local work.

# Usage

Create a logger:

	// For servers (JSON)
	logger := log.NewJSONLogger(os.Stdout, log.LogLevelInfoLevel)

	// For local work (Text)
	logger := log.NewTextLogger(os.Stdout, log.LogLevelDebugLevel)

Write logs:

	logger.Info("Hello world")
	logger.Infof("User %s joined", "Alice")
	logger.WithField("id", 123).Info("User login")

Control details:

	logger.SetLevel(log.LogLevelDebugLevel)
*/
package log
