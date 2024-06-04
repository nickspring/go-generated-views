package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"text/template"

	"golang.org/x/tools/imports"
)

const (
	commonView = "*"
)

// Regular expression to catch special tags
var specialTagRegexp = regexp.MustCompile(`^([^)(]+)\((.+)\)`)

// Generator is responsible for generating validation files for the given in a go source file.
type Generator struct {
	Version   string
	Revision  string
	BuildDate string
	BuiltBy   string
	fileSet   *token.FileSet
	t         *template.Template
	buildTags []string
}

// NewGenerator creates new Generator
func NewGenerator() *Generator {
	g := &Generator{
		Version:   "-",
		Revision:  "-",
		BuildDate: "-",
		BuiltBy:   "-",
		fileSet:   token.NewFileSet(),
		t:         template.New("generator"),
	}
	g.addEmbeddedTemplates()
	return g
}

// WithBuildTags will add build tags to the generated file.
func (g *Generator) WithBuildTags(tags ...string) *Generator {
	g.buildTags = append(g.buildTags, tags...)
	return g
}

// GenerateFromFile is responsible for orchestrating the Code generation.  It results in a byte array
// that can be written to any file desired.
func (g *Generator) GenerateFromFile(inputFile string) ([]byte, error) {
	node, err := parser.ParseFile(g.fileSet, inputFile, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("generate: error parsing input file '%s': %s", inputFile, err)
	}
	return g.Generate(node)
}

// Generate is responsible for orchestrating the code generation process.
func (g *Generator) Generate(node *ast.File) ([]byte, error) {
	structs := NewStructList(node)

	// work with header template
	pkg := node.Name.Name
	vBuff := bytes.NewBuffer([]byte{})
	err := g.t.ExecuteTemplate(vBuff, "header", map[string]interface{}{
		"package":   pkg,
		"version":   g.Version,
		"revision":  g.Revision,
		"buildDate": g.BuildDate,
		"builtBy":   g.BuiltBy,
		"buildTags": g.buildTags,
	})
	if err != nil {
		return nil, fmt.Errorf("failed writing header: %w", err)
	}

	// work with views template
	err = g.t.ExecuteTemplate(vBuff, "views", map[string]interface{}{
		"commonView": commonView,
		"structs":    structs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed writing views: %w", err)
	}

	// work with views methods
	err = g.t.ExecuteTemplate(vBuff, "methods", map[string]interface{}{
		"commonView": commonView,
		"structs":    structs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed writing methods: %w", err)
	}

	// formating code
	formatted, err := imports.Process(pkg, vBuff.Bytes(), nil)
	if err != nil {
		err = fmt.Errorf("generate: error formatting code %s\n\n%s", err, vBuff.String())
	}
	return formatted, err
}
