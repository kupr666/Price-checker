include .env
export

service-run:

	@go run cmd/api/main.go || true
	