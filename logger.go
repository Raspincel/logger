package logger

import (
	"fmt"
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
	Async,
	ForceCategoryEnabling,
	ForceLevelEnabling bool
}

type write func(e LogEntry)

type out struct {
	write write
	name  Category
}

type Logger struct {
	config            LoggerConfig
	writers           map[Writer]out
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

	l.enabledLevels[Debug] = true
	l.enabledLevels[Error] = true
	l.enabledLevels[Info] = true
	l.enabledLevels[Warn] = true
	l.enabledLevels[Misc] = true

	l.enabledCategories[Default] = true

	return l
}

func (l *Logger) AddWriter(n Writer, w write) {
	l.writers[n] = out{write: w, name: Category(n)}
}

func (l Logger) Log(d LogData) error {
	if l.config.ForceCategoryEnabling && !l.enabledCategories[d.Category] {
		return fmt.Errorf("Category %s not enabled", d.Category)
	}

	if l.config.ForceLevelEnabling && !l.enabledLevels[d.Level] {
		return fmt.Errorf("Level %s not enabled", d.Level)
	}

	fmt.Println("level", d.Level, "enabled", l.enabledLevels[d.Level])

	entry := LogEntry{
		Timestamp: time.Now(),
		Category:  d.Category,
		Level:     d.Level,
		Message:   d.Message,
		Metadata:  d.Metadata,
		Writer:    d.Writer,
	}

	w, ok := l.writers[d.Writer]

	if !ok {
		return fmt.Errorf("Writer %s not found", d.Writer)
	}

	w.write(entry)

	return nil
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
