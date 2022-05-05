package clop

import (
	"io"
	"sort"
	"strings"
	"text/template"
)

func init() {
	funcMap = template.FuncMap{
		"addSpace": addSpace,
		"sub":      sub,
	}
}

var funcMap map[string]interface{}

func addSpace(max, cur int) string {
	return strings.Repeat(" ", max-cur)
}

func sub(index int) int {
	index--
	return index
}

type showOption struct {
	Opt     string
	Usage   string
	Env     string
	Default string
}

type Help struct {
	ProcessName      string
	Version          string
	About            string
	Flags            []showOption
	Options          []showOption
	Args             []showOption
	Envs             []showOption
	Subcommand       []showOption
	MaxNameLen       int
	ShowUsageDefault bool
}

func (h *Help) output(w io.Writer) error {
	sort.Slice(h.Flags, func(i, j int) bool {
		return h.Flags[i].Opt < h.Flags[j].Opt
	})

	sort.Slice(h.Options, func(i, j int) bool {
		return h.Options[i].Opt < h.Options[j].Opt
	})

	sort.Slice(h.Subcommand, func(i, j int) bool {
		return h.Subcommand[i].Opt < h.Subcommand[j].Opt
	})
	tmpl := newTemplate()
	return tmpl.Execute(w, *h)
}

var usageDefaultTmpl = `{{- $ShowUsageDefault := .ShowUsageDefault}}{{- if gt (len .About) 0}}
{{- .About}}

{{end}}
{{- if or (gt (len .Flags) 0) (gt (len .Options) 0) (gt (len .Args) 0) (gt (len .Subcommand) 0)}}Usage:
    {{if gt (len .ProcessName) 0}}{{.ProcessName}} {{end}}
{{- if gt (len .Flags) 0}}[Flags] {{end}}
{{- if gt (len .Options) 0}}[Options] {{end}}
{{- range $_, $flag := .Args}}{{$flag.Opt}} {{end}}
{{- if gt (len .Subcommand) 0}}<Subcommand> {{end}}
{{- end}}
{{- $maxNameLen :=.MaxNameLen}}

{{- if gt (len .Flags) 0 }}

Flags:
{{- $length := len .Flags}}
{{- $length = sub $length}}
{{range $index, $flag:= .Flags}}    {{addSpace $maxNameLen (len $flag.Opt)|printf "%s%s" $flag.Opt}}    {{$flag.Usage}}
{{- if gt (len $flag.Env) 0 }} [env: {{$flag.Env}}] {{- end}}
{{- if and (gt (len $flag.Default) 0) $ShowUsageDefault}} [default: {{$flag.Default}}] {{- end}}
{{- if ne $index $length}}
{{end}}
{{- end}}

{{- end}}


{{- if gt (len .Options) 0 }}

Options:
{{- $length := len .Options}}
{{- $length = sub $length}}
{{range $index, $flag:= .Options}}    {{addSpace $maxNameLen (len $flag.Opt)|printf "%s%s" $flag.Opt}}    {{$flag.Usage}} 
{{- if gt (len $flag.Env) 0 }} [env: {{$flag.Env}}]{{- end}}
{{- if and (gt (len $flag.Default) 0 ) $ShowUsageDefault}} [default: {{$flag.Default}}]{{- end}}
{{- if ne $index $length}}
{{end}}

{{- end}}
{{- end}}


{{- if gt (len .Args) 0}}
Args:
{{- $length := len .Args}}
{{- $length = sub $length}}
{{range $index, $flag:= .Args}}    {{addSpace $maxNameLen (len $flag.Opt)|printf "%s%s" $flag.Opt}}    {{$flag.Usage}}
{{- if gt (len $flag.Env) 0 }} [env: {{$flag.Env}}]{{- end}}
{{- if ne $index $length}}
{{end}}

{{- end}}
{{- end}}


{{- if gt (len .Envs) 0}}

Environment Variable:
{{- $length := len .Envs}}
{{- $length = sub $length}}
{{range $index, $flag:= .Envs}}    {{addSpace $maxNameLen (len $flag.Opt)|printf "%s%s" $flag.Opt}}    {{$flag.Usage}}
{{- if ne $index $length}}
{{end}}

{{- end}}
{{- end}}

{{- if gt (len .Subcommand) 0 }}

Subcommand:
{{- $length := len .Subcommand}}
{{- $length = sub $length}}
{{range $index, $flag:= .Subcommand}}    {{addSpace $maxNameLen (len $flag.Opt)|printf "%s%s" $flag.Opt}}    {{$flag.Usage}} 
{{- if gt (len $flag.Env) 0 }} [env: {{$flag.Env}}]{{- end}}
{{- if ne $index $length}}
{{end}}

{{- end}}
{{- end}}
`

func newTemplate() *template.Template {
	tmpl := usageDefaultTmpl
	return template.Must(template.New("clop-default-usage").Funcs(funcMap).Parse(tmpl))
}
