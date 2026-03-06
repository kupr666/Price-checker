include .env
export

service-dev-db:
	docker compose up -d db

service-run:
	@go run cmd/api/main.go || true

service-deploy:
	docker compose up -d --build db application 
	
service-undeploy:
	docker compose down