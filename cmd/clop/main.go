package main

import (
	"fmt"
	"github.com/guonaihong/clop"
	"os"
)

type Cmd struct {
	FileName   string `clop:"short;long" usage:"go file" valid:"required"`
	OnlyStruct bool   `clop:"short;long" usage:"only struct" valid:"required"`
}

func main() {
	c := Cmd{}
	clop.Bind(&c)

	p := clop.NewParseFlag().FromFile(c.FileName)

	if c.OnlyStruct {
		p.OnlyStruct()
	} else {
		p.All()
	}

	all, err := p.Parse()
	if err != nil {
		fmt.Printf("parser:%v\n", err)
		return
	}
	os.Stdout.Write(all)
}
