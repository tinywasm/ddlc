package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tinywasm/fmt"
)

const (
	msgExporting     = "Exporting DDL schema..."
	msgExportSuccess = "Export successful. Written to: "
	msgExportFailed  = "Export failed: "
)

func collectLogs(t *testing.T, h *Handler) *[]string {
	t.Helper()
	var logs []string
	h.SetLog(func(messages ...any) {
		var strMsgs []string
		for _, m := range messages {
			strMsgs = append(strMsgs, fmt.Sprint(m))
		}
		logs = append(logs, strings.Join(strMsgs, " "))
	})
	return &logs
}

func TestHandler_WithoutExportFunc(t *testing.T) {
	h := New()
	logs := collectLogs(t, h)

	err := h.ExecuteErr()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := "ddlc handler: export function not configured"
	if err.Error() != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}

	found := false
	for _, log := range *logs {
		if strings.Contains(log, expectedErr) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error message in logs, logs were: %v", *logs)
	}
}

func TestHandler_WithExportFunc_Success_WritesFile(t *testing.T) {
	h := New()
	logs := collectLogs(t, h)

	outputPath := filepath.Join(t.TempDir(), "nested", "db.sql")
	h.SetOutputPath(outputPath)

	expectedSQL := "CREATE TABLE users (id INTEGER PRIMARY KEY);"
	h.SetExport(func() (string, error) {
		return expectedSQL, nil
	})

	if err := h.ExecuteErr(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	written, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("expected output file to exist: %v", err)
	}
	if string(written) != expectedSQL {
		t.Errorf("expected written SQL %q, got %q", expectedSQL, string(written))
	}

	foundOpen, foundClose := false, false
	for _, log := range *logs {
		if strings.Contains(log, msgExporting) {
			foundOpen = true
		}
		if strings.Contains(log, msgExportSuccess) && strings.Contains(log, outputPath) {
			foundClose = true
		}
	}
	if !foundOpen {
		t.Error("expected LogOpen message in logs")
	}
	if !foundClose {
		t.Errorf("expected LogClose success message with output path in logs, logs were: %v", *logs)
	}
}

func TestHandler_DefaultOutputPath(t *testing.T) {
	h := New()
	if h.outputPath != DefaultOutputPath {
		t.Errorf("expected default output path %q, got %q", DefaultOutputPath, h.outputPath)
	}
}

func TestHandler_WithExportFunc_Error(t *testing.T) {
	h := New()
	logs := collectLogs(t, h)

	expectedErr := fmt.Err("database connection lost")
	h.SetExport(func() (string, error) {
		return "", expectedErr
	})

	err := h.ExecuteErr()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected propagated error to be verbatim %v, got %v", expectedErr, err)
	}

	foundErr := false
	for _, log := range *logs {
		if strings.Contains(log, msgExportFailed) && strings.Contains(log, expectedErr.Error()) {
			foundErr = true
		}
	}
	if !foundErr {
		t.Errorf("expected error message in logs, logs were: %v", *logs)
	}
}

func TestHandler_Execute_Structural(t *testing.T) {
	h := New()
	h.SetOutputPath(filepath.Join(t.TempDir(), "db.sql"))
	called := false
	h.SetExport(func() (string, error) {
		called = true
		return "", nil
	})
	h.Execute()
	if !called {
		t.Error("expected export function to be called by Execute()")
	}
}
