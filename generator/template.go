package generator

import (
	"embed"
	"text/template"
)

//go:embed view.tpl
var content embed.FS

// addEmbeddedTemplates parses and adds embedded template files to the template instance
func (g *Generator) addEmbeddedTemplates() {
	g.t = template.Must(g.t.ParseFS(content, "*.tpl"))
}
