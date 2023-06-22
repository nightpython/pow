start-server:
	go run cmd/server/main.go

start-client:
	go run cmd/client/main.go

start:
	docker-compose up --force-recreate --build server --build client