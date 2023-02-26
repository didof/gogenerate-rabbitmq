BINARY_NAME=rabbitmq_publisher

generate: build
	mv -f ${BINARY_NAME} ${GOPATH}/bin
	go generate ./...

build:
	go build -o ${BINARY_NAME} publisher.go

test:
	go test ./... -v

example: FORCE
	go generate ./example/example.go
	go build -o ./example/e ./example/example.go ./example/example_rabbitmq_publisher.go
	@echo "Running the executable!\n"
	@./example/e
	@rm -f ./example/e

FORCE: ;

example-run: example
	go build -o ./example/example.go
	./e
	rm ./e

clean:
	rm -f ${BINARY_NAME}
	rm -f ${GOPATH}/bin/${BINARY_NAME}
	rm -f example/e example/example_rabbitmq_publisher.go