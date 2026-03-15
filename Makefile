include .env
export

service-dev-db:
	docker compose up -d --build db

service-dev-redis:
	docker compose up -d redis

service-run:
	@go run cmd/api/main.go || true

service-deploy:
	docker compose up -d --build db redis application 
	
service-undeploy:
	docker compose down