build:
	go build pkg/diagnostics/tracing.go
	go build pkg/diagnostics/httpTracing.go
	go build pkg/client/client.go
	go build pkg/server/server.go
	go build pkg/runtime/runtime.go

run: build
	go run cmd/server/main.go

test:
	go run cmd/client/main.go
tests:
	go run cmd/tests/main.go