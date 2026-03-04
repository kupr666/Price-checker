include .env
export

service-run:
	@go run cmd/api/main.go || true

service-deploy:
	docker compose up -d --build db application 
	
service-undeploy:
	docker compose down