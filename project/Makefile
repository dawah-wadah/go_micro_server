FRONT_END_BINARY=frontAPP
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker Images..."
	docker-compose up -docker
	@echo "Docker images started!"
	@echo "================================Nothing Good Starts In a Getaway Car==================================="

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth
	@echo "Stopping docker images (if running....)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started"
	
## down: stops docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	@echo "================================It was the best of times, the worst of crimes==================================="
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "done!"

## build_front: builds the broker binary as a linux executable
build_front:
	@echo "Building broker binary..."
	@echo "================================It was the best of times, the worst of crimes==================================="
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ../front-end && ./${FRONT_END_BINARY} &

## stop: stop the front-end
stop:
	@echo "Stopping front end..."
	@-pkilll -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end"


## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	@echo "================================A struck of match that blew your mind==================================="
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "done!"