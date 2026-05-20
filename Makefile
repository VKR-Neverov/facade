# Constants

PROJECT_NAME=facade
USER=fidesy-pay


PHONY: gg
gg:
	go mod tidy
	go run github.com/99designs/gqlgen generate

PHONY: clean
clean:
	 if docker inspect ${PROJECT_NAME} > /dev/null 2>&1; then docker rm -f ${PROJECT_NAME} && docker rmi -f ${PROJECT_NAME}; else echo "Container not found."; fi

PHONY: go-build
go-build:
	GOOS=linux GOARCH=amd64 go build -o ./main ./cmd/${PROJECT_NAME}
	mkdir -p bin
	mv main bin

PHONY: build
build:
	make go-build
	docker build --tag ${PROJECT_NAME} .

PHONY: run
run:
	make clean
	make build
	docker run --name ${PROJECT_NAME} --network=zoo -dp 7090:7090 -p 7091:7091 -e PORT=7090 -e METRICS_PORT=7091 -e APP_NAME=${PROJECT_NAME} -e ENV=local ${PROJECT_NAME}