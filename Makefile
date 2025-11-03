.PHONY: proto build docker-up docker-down test

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user.proto

build:
	go build -o bin/user-service cmd/user-service/main.go

test:
	go test ./...

clean:
	rm -rf bin/