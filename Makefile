migrate:
	go run cmd/migrate.go

format:
	go fmt ./...


docker run:
	docker run -v /var/run/docker.sock:/var/run/docker.sock -p 8080:8080 --env-file api_server/.env comet-server

docker build:
	docker build -t coderhari/comet-server .
