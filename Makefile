
.PHONY: linux
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cashier_linux ./cmd/cashier

.PHONY: darwin
darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o cashier_darwin ./cmd/cashier
