# Description: Makefile for catsocial
# Run to debug
.PHONY: run
run:
	go run cmd/main.go

# Build the project
.PHONY: build
build:
	env GOARCH=amd64 GOOS=linux go build -v -o main_syarif_04 cmd/main.go

# Build the docker image
.PHONY: docker-build
docker-build:
	docker build -t syarif/halosuster:latest .

# Down, Drop, and Up New Database
.PHONY: reset-db
reset-db:
	migrate -database "postgres://postgres:password@localhost:5432/belimang?sslmode=disable" -path db/migrations down
	migrate -database "postgres://postgres:password@localhost:5432/belimang?sslmode=disable" -path db/migrations drop
	migrate -database "postgres://postgres:password@localhost:5432/belimang?sslmode=disable" -path db/migrations up