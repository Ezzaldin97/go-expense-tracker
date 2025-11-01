package main

import (
	"go-expense-tracker/src"
	"io"
	"log"
	"os"
)

func main() {
	// Ensure the logs directory exists
	src.LoggerDir()
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Failed to open log file: %v", err)
	}
	defer file.Close()

	mw := io.MultiWriter(os.Stdout, file)
	logger := log.New(mw, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.SetPrefix("INFO: ")
	logger.Println("Hello, World!")
}
