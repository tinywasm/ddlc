// Package tui provides the DevTUI execution handler for DDL export.
//
// This package is intentionally separate from the ddlc root package: the
// root stays a WASM-safe leaf contract (Exporter, FieldExt, TopologicalSort)
// consumed by sqlt/postgres/ormc/sqlmcp, while this package pulls in os/file
// I/O and the DevTUI execution contract and is only meant to be imported by
// tinywasm/app's dev console.
package tui

import (
	"os"
	"path/filepath"

	"github.com/tinywasm/fmt"
)

// ExportFunc produces the DDL SQL for the current project.
type ExportFunc func() (sql string, err error)

// DefaultOutputPath is used when SetOutputPath is never called.
const DefaultOutputPath = "config/db.sql"

// Handler implements the TUI HandlerExecution and Loggable interfaces structurally.
type Handler struct {
	logFn      func(messages ...any)
	export     ExportFunc
	outputPath string
}

// New creates a new DDLC TUI handler.
func New() *Handler {
	return &Handler{
		logFn:      func(messages ...any) {}, // default no-op
		outputPath: DefaultOutputPath,
	}
}

// Name returns the TUI handler identifier.
func (h *Handler) Name() string {
	return "DDLC"
}

// Label returns the button label for DevTUI.
func (h *Handler) Label() string {
	return "Export DDL"
}

// SetLog injects the logging function.
func (h *Handler) SetLog(fn func(messages ...any)) {
	if fn != nil {
		h.logFn = fn
	} else {
		h.logFn = func(messages ...any) {}
	}
}

// SetExport injects the DDL export logic.
func (h *Handler) SetExport(fn ExportFunc) {
	h.export = fn
}

// SetOutputPath overrides the file the generated DDL is written to.
// Empty path resets to DefaultOutputPath.
func (h *Handler) SetOutputPath(path string) {
	if path == "" {
		path = DefaultOutputPath
	}
	h.outputPath = path
}

// Log message prefix and message constants to avoid repeating string literals.
const (
	LogOpen          = "[..."
	LogClose         = "...]"
	MsgPrefix        = "ddlc handler: "
	MsgExporting     = "Exporting DDL schema..."
	MsgExportFailed  = "Export failed: "
	MsgExportSuccess = "Export successful. Written to: "
	MsgWriteFailed   = "Write failed: "
	ErrNotConfigured = "ddlc handler: export function not configured"
)

// Execute triggers the DDL export. Implements devtui.HandlerExecution structurally.
func (h *Handler) Execute() {
	_ = h.ExecuteErr()
}

// ExecuteErr executes the export, writes the result to outputPath, and
// returns the error for test assertion.
func (h *Handler) ExecuteErr() error {
	if h.export == nil {
		err := fmt.Err(ErrNotConfigured)
		h.logFn(MsgPrefix + err.Error())
		return err
	}

	h.logFn(LogOpen, MsgPrefix+MsgExporting)
	sql, err := h.export()
	if err != nil {
		h.logFn(LogClose, MsgPrefix+MsgExportFailed+err.Error())
		return err
	}

	if err := h.writeOutput(sql); err != nil {
		h.logFn(LogClose, MsgPrefix+MsgWriteFailed+err.Error())
		return err
	}

	h.logFn(LogClose, MsgPrefix+MsgExportSuccess+h.outputPath)
	return nil
}

func (h *Handler) writeOutput(sql string) error {
	dir := filepath.Dir(h.outputPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Err(err.Error())
		}
	}
	if err := os.WriteFile(h.outputPath, []byte(sql), 0o644); err != nil {
		return fmt.Err(err.Error())
	}
	return nil
}
