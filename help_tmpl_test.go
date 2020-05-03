package clop

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func checkUsage(b *bytes.Buffer) bool {
	buf := b.Bytes()
	lines := bytes.Split(buf, []byte("\n"))

	haveUsage := false
	for _, l := range lines {
		l = bytes.TrimSpace(l)
		if bytes.HasPrefix(l, []byte("Usage:")) {
			haveUsage = true
			continue
		}

		if haveUsage {
			if bytes.Index(l, []byte("<Subcommand>")) == -1 {
				return false
			}
			return true
		}
	}

	return false
}

func haveDefaultInfo(b []byte) bool {
	return bytes.Index(b, []byte("default")) != -1
}

func Test_Usage_tmpl_CloseDefault(t *testing.T) {
	ShowUsageDefault = false
	defer func() { ShowUsageDefault = true }()

	help := Help{
		ProcessName: "test",
		Version:     "clop v0.0.1",
		About:       "guonaihong development",
		//Usage:       "--output <output> [--] [FILE]...",
		Flags: []showOption{
			{Opt: "-d, --debug", Usage: "Activate debug mode", Env: "DEBUG=", Default: "true"},
			{"-h, --help", "Prints help information", "", ""},
			{"-V, --version", "Prints version information", "", ""},
			{"-v, --verbose", "Verbose mode (-v, -vv, -vvv, etc.)", "", ""},
		},
		Options: []showOption{
			{Opt: "-l, --level <level>...", Usage: "admin_level to consider", Env: "LEVEL=debug", Default: "info"},
			{"-c, --nb-cars <nb-cars>", "Number of cars", "", ""},
			{"-o, --output <output>", "Output file", "", ""},
			{"-s, --speed <speed>", "-s, --speed <speed>", "", ""},
		},
		Args: []showOption{
			{"<api-url>", "[env: API_URL=]", "", ""},
			{"<FILE>...", "Files to process", "", ""},
		},
		Subcommand: []showOption{
			{"add", "Add file contents to the index", "", ""},
			{"mv", "Move or rename a file, a directory, or a symlink", "", ""},
		},
		MaxNameLen:       30,
		ShowUsageDefault: ShowUsageDefault,
	}

	b := bytes.Buffer{}
	w := io.MultiWriter(os.Stdout, &b)

	tmpl := newTemplate()
	err := tmpl.Execute(w, help)
	assert.NoError(t, err)

	assert.False(t, haveDefaultInfo(b.Bytes()))
	assert.True(t, checkUsage(&b))
}

func Test_Usage_tmpl(t *testing.T) {
	help := Help{
		ProcessName: "test",
		Version:     "clop v0.0.1",
		About:       "guonaihong development",
		//Usage:       "--output <output> [--] [FILE]...",
		Flags: []showOption{
			{Opt: "-d, --debug", Usage: "Activate debug mode", Env: "DEBUG=", Default: "true"},
			{"-h, --help", "Prints help information", "", ""},
			{"-V, --version", "Prints version information", "", ""},
			{"-v, --verbose", "Verbose mode (-v, -vv, -vvv, etc.)", "", ""},
		},
		Options: []showOption{
			{Opt: "-l, --level <level>...", Usage: "admin_level to consider", Env: "LEVEL=debug", Default: "info"},
			{"-c, --nb-cars <nb-cars>", "Number of cars", "", ""},
			{"-o, --output <output>", "Output file", "", ""},
			{"-s, --speed <speed>", "-s, --speed <speed>", "", ""},
		},
		Args: []showOption{
			{"<api-url>", "[env: API_URL=]", "", ""},
			{"<FILE>...", "Files to process", "", ""},
		},
		Subcommand: []showOption{
			{"add", "Add file contents to the index", "", ""},
			{"mv", "Move or rename a file, a directory, or a symlink", "", ""},
		},
		MaxNameLen:       30,
		ShowUsageDefault: ShowUsageDefault,
	}

	b := bytes.Buffer{}
	w := io.MultiWriter(os.Stdout, &b)

	tmpl := newTemplate()
	err := tmpl.Execute(w, help)
	assert.NoError(t, err)

	assert.True(t, haveDefaultInfo(b.Bytes()))
	assert.True(t, checkUsage(&b))
}
