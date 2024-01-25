mock:
	go generate -v ./...

# Docker compose commands

develop:
	echo "Starting docker environment"
	docker compose  -f docker-compose.yaml up --build 

docker_delve:
	echo "Starting docker debug environment"
	docker compose -f docker-compose.delve.yaml up --build

local:
	echo "Starting local environment"
	docker compose -f docker-compose.local.yaml up --build

stop:
	docker compose down


proto:
	rm -rf pk/pb/*.go
	protoc --proto_path=pkg/proto --go_out=pkg/pb --go_opt=paths=source_relative \
	--go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional \
	--grpc-gateway_out=pkg/pb --grpc-gateway_opt=paths=source_relative \
	pkg/proto/*.proto

local-web:
	CompileDaemon -include=".html" -include=".js" \
	 -directory="./cmd/web" -recursive=true -command="go run cmd/web/main.go"

evans:
	evans --port 8000 -r repl

.PHONY: proto