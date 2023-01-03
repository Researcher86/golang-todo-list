


up:
	docker-compose up -d

down:
	docker-compose down -v


init:
	rm -rf vendor/ &&\
	go mod download &&\
	go mod vendor

run:
	go run cmd/main.go

build-app:
	go build cmd/main.go