NAME = "hw7"

h help:
	@echo "h help - Help"
	@echo "build - Build Project"
	@echo "run - Build and Run"
	@echo "test - Run tests"
	@echo "swag - Reinitialize swagger"
	@echo "dev - Reinitialize swagger and build and run the project"
.PHONY: h help

build: swag
	go build -o ${NAME}
.PHONY: build

run: build
	./${NAME}
.PHONY: run

test:
	go test ./... -v
.PHONY: test

swag:
	swag init
.PHONY: swag

dev: swag run
.PHONY: dev
