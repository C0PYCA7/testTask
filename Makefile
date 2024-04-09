.SILENT:

run:
	go run cmd/test/main.go

docker-compose:
	docker-compose up

test-cover:
	go test --cover ./...

