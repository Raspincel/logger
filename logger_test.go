package logger

import (
	"fmt"
	"testing"
)

func TestLogger(t *testing.T) {
	l := NewLogger(LoggerConfig{})
	const writer Writer = "default"

	testCases := []LogData{
		{Message: "Hey, Planet!", Metadata: map[string]any{"key": "value"}, Writer: writer, Level: Info, Category: Default},
		{Message: "Error", Metadata: map[string]any{"key": "value"}, Writer: "window", Level: Info, Category: Default},
		{Message: "Hello, World!", Metadata: map[string]any{"key": "value"}, Writer: writer, Level: Debug, Category: "main"},
		{Message: "Hello, World!", Metadata: map[string]any{"key": "value"}, Writer: writer, Level: Info, Category: "misc"},
	}

	var actual LogEntry

	l.AddWriter(writer, func(e LogEntry) {
		actual = e
	})

	for _, tc := range testCases {
		err := l.Log(LogData{
			Message:  tc.Message,
			Metadata: tc.Metadata,
			Writer:   tc.Writer,
			Level:    tc.Level,
			Category: tc.Category,
		})

		if err != nil && err.Error() == fmt.Sprintf("Writer %s not found", tc.Writer) {
			if actual.Message == tc.Message {
				t.Errorf("Expected message %v to be different from %v", tc.Message, actual.Message)
			}

			continue
		}

		if actual.Message != tc.Message {
			fmt.Println(actual.Message, tc.Message)
			t.Errorf("Expected message %v, got %v", tc.Message, actual.Message)
		}
		if actual.Metadata["key"] != tc.Metadata["key"] {
			t.Errorf("Expected metadata %v, got %v", tc.Metadata["key"], actual.Metadata["key"])
		}
		if actual.Writer != tc.Writer {
			t.Errorf("Expected writer %v, got %v", tc.Writer, actual.Writer)
		}
		if actual.Level != tc.Level {
			t.Errorf("Expected level %v, got %v", tc.Level, actual.Level)
		}
		if actual.Category != tc.Category {
			t.Errorf("Expected category %v, got %v", tc.Category, actual.Category)
		}
	}
}

func TestRules(t *testing.T) {
	l := NewLogger(LoggerConfig{ForceCategoryEnabling: true, ForceLevelEnabling: true})

	const writer Writer = "default"

	l.EnableCategory("cat1")
	l.EnableLevel("level1")

	l.EnableCategory("cat2")
	l.EnableLevel("level2")

	l.DisableCategory("cat2")
	l.DisableLevel("level2")

	testCases := []LogData{
		{Message: "Error", Metadata: map[string]any{"result": "error"}, Writer: "window", Level: "1", Category: Default},
		{Message: "Error", Metadata: map[string]any{"result": "error"}, Writer: "window", Level: Error, Category: "main"},
		{Message: "Success", Metadata: map[string]any{"result": "success"}, Writer: writer, Level: "level1", Category: Default},
		{Message: "Success", Metadata: map[string]any{"result": "success"}, Writer: writer, Level: Info, Category: "cat1"},
		{Message: "Error", Metadata: map[string]any{"result": "error"}, Writer: writer, Level: "level2", Category: Default},
		{Message: "Error", Metadata: map[string]any{"result": "error"}, Writer: writer, Level: Info, Category: "cat2"},
	}

	l.AddWriter(writer, func(e LogEntry) {})

	for _, tc := range testCases {
		err := l.Log(LogData{
			Message:  tc.Message,
			Metadata: tc.Metadata,
			Writer:   tc.Writer,
			Level:    tc.Level,
			Category: tc.Category,
		})

		if tc.Metadata["result"].(string) == "error" && err == nil {
			t.Errorf("Expected error, got nil. Level: %v, Category: %v", tc.Level, tc.Category)
		}

		if tc.Metadata["result"].(string) != "error" && err != nil {
			t.Errorf("Expected nil, got error. Level: %v, Category: %v", tc.Level, tc.Category)
		}
	}
}
