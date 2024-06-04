.PHONY: api
api:
	protoc --proto_path=. \
           --go_out=paths=source_relative:. \
           --go-grpc_out=paths=source_relative:. \
           api/*.proto

.PHONY: build
build:
	go build -o bin/dynamic-grpc-server ./cmd/main.go