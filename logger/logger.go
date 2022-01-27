package logger

import (
	"log"
	"os"
)

type Logger struct {
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
}

func (logger *Logger) initializeFileLogger() {
	file, err := os.OpenFile("output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logger.InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (logger *Logger) initializeConsoleLogger() {
	logger.InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.WarningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func New(appendFile bool) (logger *Logger) {
	logger = new(Logger)

	if appendFile {
		logger.initializeFileLogger()
	} else {
		logger.initializeConsoleLogger()
	}

	return logger
}

func (logger *Logger) Info(msg string) {
	logger.InfoLogger.Println(msg)
}

func (logger *Logger) Warning(msg string) {
	logger.WarningLogger.Println(msg)
}

func (logger *Logger) Error(msg string) {
	logger.ErrorLogger.Println(msg)
}

func (logger *Logger) PrintError(err error) {
	logger.ErrorLogger.Println(err)
}

func (logger *Logger) Fatal(err error) {
	log.Fatal(err)
}
