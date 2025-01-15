package cli_log

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	logFile string
	maxSize int64
	mutex   sync.Mutex
	appName string
}

func NewLogger(logFile string, maxSize int64, appName string) *Logger {
	if logFile == "" {
		logFile = "log.txt"
	}
	if maxSize <= 0 {
		maxSize = 1024 * 1024 // Default to 1MB
	}
	return &Logger{logFile: logFile, maxSize: maxSize, appName: appName}
}

func (l *Logger) AddToLog(errorType, msg string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	timestamp := time.Now().Format(time.RFC3339)

	if err := l.checkLogFile(); err != nil {
		return err
	}

	file, err := os.OpenFile(l.logFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	logEntry := fmt.Sprintf("[%s] [%s] %s\n", strings.ToUpper(errorType), timestamp, msg)
	if _, err = file.WriteString(logEntry); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}
	return nil
}

func (l *Logger) checkLogFile() error {
	info, err := os.Stat(l.logFile)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(l.logFile)
			if err != nil {
				return fmt.Errorf("failed to create log file: %w", err)
			}
			file.Close()
			return nil
		}
		return fmt.Errorf("failed to check log file: %w", err)
	}

	if info.Size() > l.maxSize {
		if err = os.Remove(l.logFile); err != nil {
			return fmt.Errorf("failed to remove old log file: %w", err)
		}
		file, err := os.Create(l.logFile)
		if err != nil {
			return fmt.Errorf("failed to create new log file: %w", err)
		}
		file.Close()
	}
	return nil
}

func (l *Logger) Notify(msg string) error {
	cmd := exec.Command("notify-send", "-i", "preferences-system", l.appName, msg)
	return cmd.Run()
}
