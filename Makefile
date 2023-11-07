MODULE_NAME=${shell cat go.mod | grep module | awk '{print $$2}' }

build:
	@export CGO_ENABLED=1
	@go build -o bin/${MODULE_NAME}

run: build
	bin/${MODULE_NAME}

clean:
	@rm ${MODULE_NAME} 2>/dev/null || true
	