package internals

import (
	"go/ast"
)

type Struct struct {
	ast.StructType
}

type Field struct {
	name      *ast.Ident
	fieldType ast.Expr
}

func (s *Struct) Iter() []*Field {
	list := []*Field{}

	for _, field := range s.Fields.List {
		list = append(list, &Field{
			name:      field.Names[0],
			fieldType: field.Type,
		})
	}

	return list
}
