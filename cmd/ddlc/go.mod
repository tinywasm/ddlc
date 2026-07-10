module github.com/tinywasm/ddlc/cmd/ddlc

go 1.25.2

require (
	github.com/tinywasm/ddlc v0.0.2
	github.com/tinywasm/ormc v0.0.1
	github.com/tinywasm/postgres v0.3.4
	github.com/tinywasm/sqlt v0.0.6
)

require (
	github.com/lib/pq v1.11.2 // indirect
	github.com/tinywasm/fmt v0.25.1 // indirect
	github.com/tinywasm/model v0.0.8 // indirect
	github.com/tinywasm/modfind v0.0.4 // indirect
	github.com/tinywasm/orm v0.9.27 // indirect
)

// TODO(publish)
replace github.com/tinywasm/ddlc => ../..

// TODO(publish)
replace github.com/tinywasm/ormc => ../../../ormc
