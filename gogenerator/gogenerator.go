package gogenerator

import (
	"fmt"
	"io"
	"text/template"

	"golang.org/x/mod/modfile"
)

// Imports ...
type Imports struct {
	PackageName string
	ImportsArr  []*modfile.Require
}

const tmplImports = `package {{ .PackageName }}

import (
	{{- range $i, $imp := .ImportsArr}}
	_ "{{$imp.Mod.Path}}"
	{{- end }}
)
`

// GenerateImportsFile ...
func GenerateImportsFile(imports Imports, w io.Writer) error {
	required := make([]*modfile.Require, 0)
	for _, v := range imports.ImportsArr {
		if !v.Indirect {
			required = append(required, v)
		}
	}
	if len(required) == 0 {
		fmt.Printf("\nWarning! no required imports present in the go.mod file")
		return nil
	}
	imports.ImportsArr = required
	tmpl := template.Must(template.New("test").Parse(tmplImports))
	return tmpl.Execute(w, imports)
}
