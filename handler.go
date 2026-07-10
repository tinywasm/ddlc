package ddlc

import (
	"github.com/tinywasm/fmt"
)

// ExportFunc produces the DDL SQL for the current project.
type ExportFunc func() (sql string, err error)

// Handler implements the TUI HandlerExecution and Loggable interfaces structurally.
type Handler struct {
	logFn  func(messages ...any)
	export ExportFunc
}

// New creates a new DDLC TUI handler.
func New() *Handler {
	return &Handler{
		logFn: func(messages ...any) {}, // default no-op
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

// Log message prefix and message constants to avoid repeating string literals.
const (
	LogOpen          = "[..."
	LogClose         = "...]"
	MsgPrefix        = "ddlc handler: "
	MsgExporting     = "Exporting DDL schema..."
	MsgExportFailed  = "Export failed: "
	MsgExportSuccess = "Export successful. Length: "
	MsgBytes         = " bytes"
	ErrNotConfigured = "ddlc handler: export function not configured"
)

// Execute triggers the DDL export. Implements devtui.HandlerExecution structurally.
func (h *Handler) Execute() {
	_ = h.ExecuteErr()
}

// ExecuteErr executes the export and returns the error for test assertion.
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

	h.logFn(LogClose, MsgPrefix+MsgExportSuccess+fmt.Sprint(len(sql))+MsgBytes)
	return nil
}
