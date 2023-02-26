package testdata

//go:generate go run publisher.go -type=Foo

type Foo struct {
	Bar bool
	Baz int
}
