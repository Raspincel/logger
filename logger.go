package logger

import (
	"fmt"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Category  Category
	Writer    Writer
	Level     LogLevel
	Message   string
	Metadata  map[string]any
}

type Writer string
type Category string
type LogLevel string

const (
	Default Category = "DEFAULT_CATEGORY"
)

const (
	Debug LogLevel = "DEBUG_LEVEL"
	Error LogLevel = "ERROR_LEVEL"
	Info  LogLevel = "INFO_LEVEL"
	Warn  LogLevel = "WARN_LEVEL"
	Misc  LogLevel = "MISC_LEVEL"
)

type LoggerConfig struct {
	AllowLoggingDisabled  bool // Allow logging with disabled levels and categories without error
	EnableDefaultLevels   bool // Register the default levels of logging (error, info, warning, etc)
	EnableDefaultCategory bool // Register the default category of logging
	UseLock               bool // Lock the logger for concurrent writes
}

type write func(e LogEntry)

type out struct {
	write write
	name  Category
}

type Logger struct {
	config  LoggerConfig
	writers map[Writer]out
	m       *sync.Mutex

	enabledCategories map[Category]bool
	enabledLevels     map[LogLevel]bool
}

type LogData struct {
	Message  string
	Metadata map[string]any
	Writer   Writer
	Level    LogLevel
	Category Category
}

func NewLogger(config LoggerConfig) *Logger {
	l := &Logger{config: config}

	l.enabledCategories = make(map[Category]bool)
	l.enabledLevels = make(map[LogLevel]bool)

	l.writers = make(map[Writer]out)

	l.config = config

	if l.config.EnableDefaultLevels {
		l.EnableLevel(Debug)
		l.EnableLevel(Error)
		l.EnableLevel(Info)
		l.EnableLevel(Warn)
		l.EnableLevel(Misc)
	}

	if l.config.EnableDefaultCategory {
		l.EnableCategory(Default)
	}

	if l.config.UseLock {
		l.m = &sync.Mutex{}
	}

	return l
}

func (l *Logger) AddWriter(n Writer, w write) {
	l.writers[n] = out{write: w, name: Category(n)}
}

func (l Logger) Log(d LogData) error {
	if !l.enabledCategories[d.Category] && !l.config.AllowLoggingDisabled {
		return fmt.Errorf("Category %s not enabled", d.Category)
	}

	if !l.enabledLevels[d.Level] && !l.config.AllowLoggingDisabled {
		return fmt.Errorf("Level %s not enabled", d.Level)
	}

	w, ok := l.writers[d.Writer]

	if !ok {
		return fmt.Errorf("Writer %s not found", d.Writer)
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Category:  d.Category,
		Level:     d.Level,
		Message:   d.Message,
		Metadata:  d.Metadata,
		Writer:    d.Writer,
	}

	if l.config.UseLock {
		l.m.Lock()
		defer l.m.Unlock()
	}

	w.write(entry)

	return nil
}

func (l *Logger) DisableAllCategories() {
	for k := range l.enabledCategories {
		l.enabledCategories[k] = false
	}
}

func (l *Logger) DisableAllLevels() {
	for k := range l.enabledLevels {
		l.enabledLevels[k] = false
	}
}

func (l *Logger) EnableCategory(c Category) {
	l.enabledCategories[c] = true
}

func (l *Logger) DisableCategory(c Category) {
	l.enabledCategories[c] = false
}

func (l *Logger) EnableLevel(lv LogLevel) {
	l.enabledLevels[lv] = true
}

func (l *Logger) DisableLevel(lv LogLevel) {
	l.enabledLevels[lv] = false
}
