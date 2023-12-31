MODULE_NAME=${shell cat go.mod | grep module | awk '{print $$2}' }

build:
	mkdir -p bin
	export CGO_ENABLED=1
	go build -o bin/${MODULE_NAME}

run: build
	bin/${MODULE_NAME}

test:
	go test -v ./...

quality:
	go vet ./...
