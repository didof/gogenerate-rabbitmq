BINARY_NAME=rabbitmq_publisher

generate: build
	mv -f ${BINARY_NAME} ${GOPATH}/bin
	go generate ./...

build:
	go build -o ${BINARY_NAME} publisher.go

test:
	go test ./... -v

example:
	go generate example.go

clean:
	rm -f ${BINARY_NAME}
	rm -f ${GOPATH}/bin/${BINARY_NAME}
	rm -f example_rabbitmq_publisher.go