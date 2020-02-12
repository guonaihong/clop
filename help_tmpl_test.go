package clop

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_Usage_tmpl(t *testing.T) {
	help := Help{
		ProcessName: "test",
		Version:     "clop v0.0.1",
		About:       "guonaihong development",
		//Usage:       "--output <output> [--] [FILE]...",
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
		Args: []showOption{
			{"<api-url>", "[env: API_URL=]", ""},
			{"<FILE>...", "Files to process", ""},
		},
		Subcommand: []showOption{
			{"add", "Add file contents to the index", ""},
			{"mv", "Move or rename a file, a directory, or a symlink", ""},
		},
		MaxNameLen: 30,
	}

	tmpl := newTemplate()
	err := tmpl.Execute(os.Stdout, help)
	assert.NoError(t, err)
}
