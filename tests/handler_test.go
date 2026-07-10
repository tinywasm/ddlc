package tests

import (
	"strings"
	"testing"

	"github.com/tinywasm/ddlc"
	"github.com/tinywasm/fmt"
)

const (
	msgExporting     = "Exporting DDL schema..."
	msgExportSuccess = "Export successful. Length: "
	msgExportFailed  = "Export failed: "
)

func TestHandler_WithoutExportFunc(t *testing.T) {
	h := ddlc.New()
	var logs []string
	h.SetLog(func(messages ...any) {
		var strMsgs []string
		for _, m := range messages {
			strMsgs = append(strMsgs, fmt.Sprint(m))
		}
		logs = append(logs, strings.Join(strMsgs, " "))
	})

	err := h.ExecuteErr()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := "ddlc handler: export function not configured"
	if err.Error() != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}

	// Verify the log contains the error message prefix and contents
	found := false
	for _, log := range logs {
		if strings.Contains(log, expectedErr) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error message in logs, logs were: %v", logs)
	}
}

func TestHandler_WithExportFunc_Success(t *testing.T) {
	h := ddlc.New()
	var logs []string
	h.SetLog(func(messages ...any) {
		var strMsgs []string
		for _, m := range messages {
			strMsgs = append(strMsgs, fmt.Sprint(m))
		}
		logs = append(logs, strings.Join(strMsgs, " "))
	})

	expectedSQL := "CREATE TABLE users (id INTEGER PRIMARY KEY);"
	h.SetExport(func() (string, error) {
		return expectedSQL, nil
	})

	err := h.ExecuteErr()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify execution triggered the export func and logged success
	foundOpen := false
	foundClose := false
	for _, log := range logs {
		if strings.Contains(log, msgExporting) {
			foundOpen = true
		}
		if strings.Contains(log, msgExportSuccess) && strings.Contains(log, fmt.Sprint(len(expectedSQL))) {
			foundClose = true
		}
	}

	if !foundOpen {
		t.Error("expected LogOpen message in logs")
	}
	if !foundClose {
		t.Errorf("expected LogClose success message with length in logs, logs were: %v", logs)
	}
}

func TestHandler_WithExportFunc_Error(t *testing.T) {
	h := ddlc.New()
	var logs []string
	h.SetLog(func(messages ...any) {
		var strMsgs []string
		for _, m := range messages {
			strMsgs = append(strMsgs, fmt.Sprint(m))
		}
		logs = append(logs, strings.Join(strMsgs, " "))
	})

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

	// Verify error log message
	foundErr := false
	for _, log := range logs {
		if strings.Contains(log, msgExportFailed) && strings.Contains(log, expectedErr.Error()) {
			foundErr = true
		}
	}
	if !foundErr {
		t.Errorf("expected error message in logs, logs were: %v", logs)
	}
}

func TestHandler_Execute_Structural(t *testing.T) {
	h := ddlc.New()
	called := false
	h.SetExport(func() (string, error) {
		called = true
		return "", nil
	})
	// Should run without panic or returning error
	h.Execute()
	if !called {
		t.Error("expected export function to be called by Execute()")
	}
}
