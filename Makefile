version ?= latest
img      = egurnov/maze-api:$(version)
imgdev   = egurnov/maze-api:$(version)

.PHONY: test
test:
	go test ./maze-api/...

.PHONY: e2e-test
e2e-test:
	go test ./test/...

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run

.PHONY: build
build:
	go build -o ./build/maze-api ./cmd/maze-api

.PHONY: run
run:
	JWT_SIGNING_KEY=changeme \
    DB_URL='root:root@(localhost:3306)/maze_api' \
	go run ./cmd/maze-api

.PHONY: clean
clean:
	rm -rf ./build

.PHONY: all
all: lint test e2e-test clean build

.PHONY: docker-build
docker-build:
	docker build . -t $(img)

.PHONY: docker-build-dev
docker-build-dev:
	docker build . --target builder -t $(imgdev)

.PHONY: docker-run
docker-run:
	docker run --rm -p 8080:8080 \
      -e JWT_SINGING_KEY=changeme \
      -e DB_URL='root@(host.docker.internal:3306)/maze-api' \
      $(img)

.PHONY: docker-lint
docker-lint:
	docker run \
	  --rm \
	  --volume $(shell pwd):/app \
	  --workdir /app \
	  golangci/golangci-lint:v1.30.0 \
	  golangci-lint run -v

.PHONY: docker-test
docker-test:
	docker run --rm $(imgdev) go test ./maze-api/... ./test/...

.PHONY: docker-all
docker-all: docker-lint docker-build-dev docker-test docker-build

.PHONY: up
up:
	docker-compose up

.PHONY: down
down:
	docker-compose down

.PHONY: swag
swag:
	swag init -dir ./maze-api --generalInfo ./app/app.go

.PHONY: generate
generate:
	(cd ./maze-api/app && go generate)
