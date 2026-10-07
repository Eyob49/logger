package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Logger struct {
	destination io.Writer
	minLevel    Level
	file        *os.File
	mu          sync.Mutex
}

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Mode int

const (
	ModeStdout Mode = iota
	ModeFile
	ModeBoth
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Close is idempotent so callers can safely defer it or call it explicitly.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		file := l.file
		l.file = nil
		return file.Close()
	}
	return nil
}

func fileOpen(filePath string) (*os.File, error) {
	writer, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	return writer, err
}

func New(mode Mode, level Level, filePath string) (*Logger, error) {
	var writer io.Writer
	var err error
	var file *os.File
	switch mode {
	case ModeStdout:
		writer = os.Stdout
	case ModeFile:
		file, err = fileOpen(filePath)
		if err != nil {
			return nil, err
		}
		writer = file
	case ModeBoth:
		file, err = fileOpen(filePath)
		if err != nil {
			return nil, err
		}
		writer = io.MultiWriter(os.Stdout, file)
	default:
		return nil, fmt.Errorf("invalid mode: %v", mode)
	}
	return &Logger{
		destination: writer,
		minLevel:    level,
		file:        file,
	}, nil
}

func (l *Logger) log(level Level, message string) {
	if l.minLevel <= level {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		logMessage := timestamp + " [" + level.String() + "] " + message + "\n"
		l.mu.Lock()
		defer l.mu.Unlock()
		_, err := l.destination.Write([]byte(logMessage))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to write log message:", err)
		}
	}
}

func (l *Logger) Debug(message string) {
	l.log(LevelDebug, message)
}

func (l *Logger) Info(message string) {
	l.log(LevelInfo, message)
}

func (l *Logger) Warn(message string) {
	l.log(LevelWarn, message)
}

func (l *Logger) Error(message string) {
	l.log(LevelError, message)
}
