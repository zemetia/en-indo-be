# Every Nation Indonesia - Backend Commands

# Dependencies
dep: 
	go mod tidy

# Local development
run: 
	go run main.go

build: 
	go build -o main main.go

run-build: build
	./main

# Testing
test:
	go test -v ./tests

# Database operations (run locally)
migrate:
	go run main.go --migrate

seed: 
	go run main.go --seed

migrate-seed: 
	go run main.go --migrate --seed