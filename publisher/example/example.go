package main

//go:generate go run ../generator/publisher.go -type=Msg
// You can customize output via flag -output=custom.go

type Msg struct {
	Id string
}

func main() {
	r := &RabbitMQPublisherMsg{}

	r.Eureka()
}
