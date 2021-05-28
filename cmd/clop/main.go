package main

import (
	"github.com/guonaihong/clop"
)

type Cmd struct {
	FileName string `clop:"short;long" usage:"go file" valid:"required"`
}

func main() {
	c := Cmd{}
	clop.Bind(&c)

	p := clop.NewParseFlag().FromFile(c.FileName)
	all, err := p.Parse()
	if err != nil {
		fmt.Printf("parser:%v\n", err)
		return
	}
	os.Stdout.Write(all)
}
