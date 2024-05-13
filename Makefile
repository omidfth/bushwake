run:
	cd ./cmd/api && go run main.go

tidy:
	export GOPROXY=https://goproxy.io,direct && go mod tidy

build:
	cd ./cmd/api && go build -o bushwake

up:
	docker compose up -d