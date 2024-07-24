package default_logger

import (
	"log"
	"os"

	"github.com/MohammadAsDev/geo_tracker/src/interfaces"
)

type DefaultLogger struct {
	_InfoLogger    *log.Logger
	_SuccessLogger *log.Logger
	_WarningLogger *log.Logger
	_ErrorLogger   *log.Logger
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		_InfoLogger:    log.New(os.Stdout, interfaces.InfoColor.Sprint("INFO:"), log.Ldate|log.Ltime),
		_SuccessLogger: log.New(os.Stdout, interfaces.CompleteColor.Sprint("SUCCESS:"), log.Ldate|log.Ltime),
		_WarningLogger: log.New(os.Stdout, interfaces.WarningColor.Sprint("WARNING:"), log.Ldate|log.Ltime),
		_ErrorLogger:   log.New(os.Stdout, interfaces.ErrorColor.Sprint("ERROR:"), log.Ldate|log.Ltime),
	}
}

func (logger *DefaultLogger) PrintSuccessLogMessage(msg string) error {
	logger._SuccessLogger.Println(msg)
	return nil
}

func (logger *DefaultLogger) PrintWarningLogMessage(msg string) error {
	logger._WarningLogger.Println(msg)
	return nil
}

func (logger *DefaultLogger) PrintErrorLogMessage(msg string) error {
	logger._ErrorLogger.Println(msg)
	return nil
}

func (logger *DefaultLogger) PrintInfoLogMessage(msg string) error {
	logger._InfoLogger.Println(msg)
	return nil
}
