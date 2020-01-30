package clop

import (
	"io"
	"text/template"
)

type showOption struct {
	Opt   string
	Usage string
}

type Help struct {
	Version string
	About   string
	Usage   string
	Flags   []showOption
	Options []showOption
	Args    showOption
}

func (h *Help) output(w io.Writer) error {
	tmpl := newTemplate()
	return tmpl.Execute(w, *h)
}

var usageDefaultTmpl = `{{if gt (len .Version) 0}}{{.Version}}{{end}}
{{if gt (len .About) 0}}{{.About}}{{end}}
{{if gt (len .Usage) 0 }}Usage:
    {{.Usage}}
{{end}}
{{if gt (len .Flags) 0 }}Flags:
{{range $_, $flag:= .Flags}}    {{$flag.Opt}}    {{$flag.Usage}}
{{end}}{{end}}
{{if gt (len .Options) 0 }}Options:
{{range $_, $flag:= .Options}}    {{$flag.Opt}}    {{$flag.Usage}}
{{end}}{{end}}
{{if gt (len .Args.Opt) 0}}Args:
    {{.Args.Opt}}    {{.Args.Usage}}
{{end}}`

func newTemplate() *template.Template {
	tmpl := usageDefaultTmpl
	return template.Must(template.New("clop-default-usage").Parse(tmpl))
}
