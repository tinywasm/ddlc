module github.com/tinywasm/ddlc

go 1.25.2

// Leaf guarantee: github.com/tinywasm/ddlc must remain a leaf package in the SQL ecosystem.
// It must depend ONLY on tinywasm/model and tinywasm/fmt to keep it portable for WASM/frontend.
require (
	github.com/tinywasm/fmt v0.25.5
	github.com/tinywasm/model v0.1.0
)

require github.com/tinywasm/tui v0.1.1
