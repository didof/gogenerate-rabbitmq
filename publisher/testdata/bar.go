package testdata

//go:generate go run publisher.go -type=Bar

type Bar struct {
	Foo bool
	Baz int
}
