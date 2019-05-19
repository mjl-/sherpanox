default:
	-mkdir assets 2>/dev/null
	go build
	go vet
	go run vendor/golang.org/x/lint/golint/*.go . client
	go run vendor/github.com/mjl-/sherpadoc/cmd/sherpadoc/*.go Example >assets/example.json
	./sherpanox
