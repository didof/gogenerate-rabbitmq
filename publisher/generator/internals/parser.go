package internals

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type Parser struct {
	fset    *token.FileSet
	astFile *ast.File

	data ParserData
}

func NewParser(filename string) *Parser {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		panic("cannot parse file")
	}

	return &Parser{
		fset:    fset,
		astFile: astFile,

		data: ParserData{
			structType: make(map[string]*Struct),
		},
	}
}

type ParserData struct {
	packageName string
	structType  map[string]*Struct
}

func (p *Parser) getPackageName() string {
	return p.astFile.Name.Name
}

func (p *Parser) getStructByName(filename, structName string) (*Struct, error) {
	s, found := p.data.structType[structName]
	if found {
		return s, nil
	}

	for _, decl := range p.astFile.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok && typeSpec.Name.Name == structName {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						// Found the struct by name.
						s := &Struct{*structType}
						p.data.structType[structName] = s
						return s, nil
					}
				}
			}
		}
	}

	// Struct not found by name.
	return nil, fmt.Errorf("struct not found: %v", structName)
}
