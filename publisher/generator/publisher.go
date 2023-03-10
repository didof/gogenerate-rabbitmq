//go:build ignore

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/didof/gogenerate-rabbitmq/publisher/generator/internals"
)

const (
	suffix = "rabbitmq_publisher"
)

var (
	typeName = flag.String("type", "", "struct representing the amqp message; must be set")
	output   = flag.String("output", "", fmt.Sprintf("output file name; default srcdir/<type>_%s.go", suffix))
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", suffix))
	flag.Parse()

	if len(*typeName) == 0 {
		log.Fatal("TODO Usage")
		os.Exit(2)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	name := os.Getenv("GOFILE")
	if len(*output) == 0 {
		if name == "" {
			log.Fatalln("This file must be run via go:generate")
		}
		*output = AddSuffix(name, suffix)
	} else if filepath.Ext(*output) != ".go" {
		log.Fatalln("Output has wrong extention")
	}

	g := internals.NewGenerator(name, filepath.Join(dir, *output))

	s, err := g.GetStructByName(*typeName)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range s.Iter() {
		fmt.Println(f)
	}

	data := g.GetData()

	g.Sprintf("// Code generated by \"go run publisher.go %s\"; DO NOT EDIT.\n\n", strings.Join(os.Args[1:], " "))
	g.Sprintf("package %s\n\n", data.PackageName)
	g.Sprintf("import \"fmt\"\n\n")
	g.Sprintf("type RabbitMQPublisher%s struct {\n", *typeName)
	g.Sprintf("\tmsg *%s\n", *typeName)
	g.Sprintf("}\n\n")
	g.Sprintf("func (s *RabbitMQPublisher%s) Eureka() {\n\tfmt.Println(\"eureka!\")\n}", *typeName)

	g.Write()
}

func BaseWithoutExt(name string) string {
	return strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
}

func AddSuffix(base, suffix string) string {
	base = BaseWithoutExt(base)
	return fmt.Sprintf("%s_%s.go", base, suffix)
}

func CopyIntoDir(src, dst string) (string, error) {
	// Create destination file
	d, err := os.Create(filepath.Join(dst, filepath.Base(src)))
	if err != nil {
		return "", fmt.Errorf("error creating destination file: %v", err)
	}
	defer d.Close()

	// Open the source file
	s, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("error opening file %s: %v", src, err)
	}
	defer s.Close()

	// Copy the source file contents into the destination file
	_, err = io.Copy(d, s)
	if err != nil {
		return "", fmt.Errorf("error copying file contents from %s to %s: %v", s.Name(), d.Name(), err)
	}

	return d.Name(), nil
}
