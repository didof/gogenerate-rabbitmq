//go:build ignore

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	g := NewGenerator(filepath.Join(dir, *output))

	g.Add("package main")

	_, err = g.parser.findStructByName(name, *typeName)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range g.parser.Fields() {
		fmt.Printf("Field: %v, Type: %v\n", f.name, f.fieldType)
	}

	g.Write()
}

type Parser struct {
	structType *ast.StructType
}

type Field struct {
	name      *ast.Ident
	fieldType ast.Expr
}

func (p *Parser) Fields() []*Field {
	list := []*Field{}

	for _, field := range p.structType.Fields.List {
		list = append(list, &Field{
			name:      field.Names[0],
			fieldType: field.Type,
		})
	}

	return list
}

func (p *Parser) findStructByName(filename, structName string) (*ast.StructType, error) {
	// Create a new file set and parse the Go file.
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing file: %v", err)
	}

	// Traverse the AST and find the struct by name.
	for _, decl := range astFile.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok && typeSpec.Name.Name == structName {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						// Found the struct by name.
						p.structType = structType
						return structType, nil
					}
				}
			}
		}
	}

	// Struct not found by name.
	return nil, fmt.Errorf("struct not found: %v", structName)
}

type Generator struct {
	path   string
	buffer strings.Builder
	parser Parser
}

func NewGenerator(path string) *Generator {
	return &Generator{
		path: path,
	}
}

func (g *Generator) Add(input string) error {
	_, err := g.buffer.Write([]byte(input))
	return err
}

func (g *Generator) Write() {
	if err := os.WriteFile(g.path, []byte(g.buffer.String()), os.ModePerm); err != nil {
		log.Fatal(err)
	}
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
