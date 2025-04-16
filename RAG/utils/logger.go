package utils

import (
	"log"
	"os"
)

var Logger *log.Logger

func InitLogger() {
	logFile := Config.GetString("log.file")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	Logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}
