package main

import "fmt"

//go:generate go run ../generator/publisher.go -type=Msg
// You can customize output via flag -output=custom.go

type Msg struct {
	Id string
}

func main() {
	// FIXME Even if found in the IDE, it is not found when building.
	r := &RabbitMQPublisherMsg{}

	fmt.Println(r)
}
