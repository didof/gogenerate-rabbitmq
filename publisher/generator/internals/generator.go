package internals

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Generator struct {
	in, out string
	buffer  strings.Builder
	parser  *Parser
}

func (g *Generator) GetStructByName(structName string) (*Struct, error) {
	return g.parser.getStructByName(g.in, structName)
}

func NewGenerator(in, out string) *Generator {
	parser := NewParser(in)

	parser.data.packageName = parser.getPackageName()

	return &Generator{
		in:  in,
		out: out,

		parser: parser,
	}
}

func (g *Generator) Sprintf(format string, a ...any) error {
	_, err := g.buffer.Write([]byte(fmt.Sprintf(format, a...)))
	return err
}

func (g *Generator) Write() {
	if err := os.WriteFile(g.out, []byte(g.buffer.String()), os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

type GeneratorData struct {
	PackageName string
}

func (g *Generator) GetData() GeneratorData {
	return GeneratorData{
		PackageName: g.parser.data.packageName,
	}
}
