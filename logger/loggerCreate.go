package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func CreateLogger(processingName string) (*log.Logger, *os.File, error) {
	now := time.Now()
	filePath := fmt.Sprintf("logger/logs/processing_%s_%s.log", processingName, now.UTC().Format("20060102_150405"))
	logFile, err := os.Create(filePath)

	if err != nil {
		return nil, nil, fmt.Errorf("ReadFileId: ошибка чтения файла: \n%w", err)
	}

	mw := io.MultiWriter(logFile, os.Stdout)

	logger := log.New(
		mw,
		"",
		log.Ldate|log.Ltime|log.Lmicroseconds,
	)
	return logger, logFile, nil
}
