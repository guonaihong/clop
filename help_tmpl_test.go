package clop

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_Usage_tmpl(t *testing.T) {
	help := Help{
		Version: "clop v0.0.1",
		About:   "guonaihong development",
		Usage:   "test [FLAGS] [OPTIONS] --output <output> [--] [FILE]...",
		Flags: []showOption{
			{"-d, --debug", "Activate debug mode", "DEBUG="},
			{"-h, --help", "Prints help information", ""},
			{"-V, --version", "Prints version information", ""},
			{"-v, --verbose", "Verbose mode (-v, -vv, -vvv, etc.)", ""},
		},
		Options: []showOption{
			{"-l, --level <level>...", "admin_level to consider", "LEVEL=debug"},
			{"-c, --nb-cars <nb-cars>", "Number of cars", ""},
			{"-o, --output <output>", "Output file", ""},
			{"-s, --speed <speed>", "-s, --speed <speed>", ""},
		},
		Args:       showOption{"<FILE>...", "Files to process", ""},
		MaxNameLen: 30,
	}

	tmpl := newTemplate()
	err := tmpl.Execute(os.Stdout, help)
	assert.NoError(t, err)
}
