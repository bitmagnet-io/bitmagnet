package directive

import (
	"github.com/vektah/gqlparser/v2/ast"
)

type (
	Directive        map[string]string
	FieldDirectives  map[string]Directive
	TypeDirectives   map[string]FieldDirectives
	SchemaDirectives map[string]TypeDirectives
)

func ExtractSchemaDirectives(schema *ast.Schema) SchemaDirectives {
	directives := make(SchemaDirectives)

	for _, def := range schema.Types {
		thisType := ExtractTypeDirectives(def)

		if len(thisType) > 0 {
			directives[def.Name] = thisType
		}
	}

	return directives
}

func ExtractTypeDirectives(def *ast.Definition) TypeDirectives {
	directives := make(TypeDirectives)

	for _, def := range def.Fields {
		thisField := ExtractFieldDirectives(def)

		if len(thisField) > 0 {
			directives[def.Name] = thisField
		}
	}

	return directives
}

func ExtractFieldDirectives(def *ast.FieldDefinition) FieldDirectives {
	directives := make(FieldDirectives)

	for _, def := range def.Directives {
		thisDir := make(Directive)

		for _, arg := range def.Arguments {
			thisDir[arg.Name] = arg.Value.Raw
		}

		directives[def.Name] = thisDir
	}

	return directives
}
