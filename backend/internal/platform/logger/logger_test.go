package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input   string
		want    slog.Level
		wantErr bool
	}{
		{input: "info", want: slog.LevelInfo},
		{input: "INFO", want: slog.LevelInfo},
		{input: "warn", want: slog.LevelWarn},
		{input: "error", want: slog.LevelError},
		{input: " error ", want: slog.LevelError},
		{input: "debug", wantErr: true},
		{input: "", wantErr: true},
	}

	for _, tt := range tests {
		got, err := ParseLevel(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseLevel(%q) = %v, want error", tt.input, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseLevel(%q) returned unexpected error: %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestSetupFiltersBelowConfiguredLevel(t *testing.T) {
	log, err := Setup("warn")
	if err != nil {
		t.Fatalf("Setup() error = %v", err)
	}

	if log.Enabled(nil, slog.LevelInfo) {
		t.Fatal("expected info level to be filtered out at warn level")
	}
	if !log.Enabled(nil, slog.LevelWarn) {
		t.Fatal("expected warn level to be enabled at warn level")
	}
	if !log.Enabled(nil, slog.LevelError) {
		t.Fatal("expected error level to be enabled at warn level")
	}
}

func TestSetupRejectsInvalidLevel(t *testing.T) {
	if _, err := Setup("verbose"); err == nil {
		t.Fatal("Setup() succeeded, want error for invalid level")
	}
}

func TestSetupWritesJSON(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	log := slog.New(handler)

	log.Info("test message", "event", "test.event")

	var decoded map[string]any
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if decoded["msg"] != "test message" {
		t.Fatalf("msg = %v, want %q", decoded["msg"], "test message")
	}
	if decoded["event"] != "test.event" {
		t.Fatalf("event = %v, want %q", decoded["event"], "test.event")
	}
}
