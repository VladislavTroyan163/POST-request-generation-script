package logger

import (
	"log"
	"os"
)

func CreateLogger() (*log.Logger, *os.File) {
	logFile, err := os.OpenFile(
		"logger\\processing.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(
		logFile,
		"",
		log.Ldate|log.Ltime|log.Lmicroseconds,
	)
	return logger, logFile
}
