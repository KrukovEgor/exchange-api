ENTRY_POINT=./cmd/api/main.go

BUILD_DIR=./build
BIN_FILE=${BUILD_DIR}/api

.PHONY: build run lint format clean rebuild

build:
	@mkdir -p ${BUILD_DIR}
	go build -o ${BIN_FILE} ${ENTRY_POINT}

run: build
	${BIN_FILE}

lint:
	golangci-lint run

format:
	go fmt ./...

clean:
	rm -f ${BIN_FILE}

rebuild: clean build
