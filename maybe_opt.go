package clop

import (
	"fmt"

	"github.com/antlabs/strsim"
)

func (c *Clop) maybeOpt(optionName string) string {
	opts := make([]string, len(c.shortAndLong))
	index := 0
	for k := range c.shortAndLong {
		opts[index] = k
		index++
	}

	m := strsim.FindBestMatchOne(optionName, opts)
	if m.Score > 0.0 {
		return m.S
	}

	return ""
}

func (c *Clop) genMaybeHelpMsg(optionName string) string {
	if s := c.maybeOpt(optionName); len(s) > 0 {
		return fmt.Sprintf("\n	Did you mean --%s?\n", s)
	}

	return ""
}
