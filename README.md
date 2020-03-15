# clop
[![Go](https://github.com/guonaihong/clop/workflows/Go/badge.svg)](https://github.com/guonaihong/clop/actions)
[![codecov](https://codecov.io/gh/guonaihong/clop/branch/master/graph/badge.svg)](https://codecov.io/gh/guonaihong/clop)
[![Go Report Card](https://goreportcard.com/badge/github.com/guonaihong/clop)](https://goreportcard.com/report/github.com/guonaihong/clop)

clop 是一款基于struct的命令行解析器，麻雀虽小，五脏俱全。(从零实现)

## feature
* 支持环境变量绑定 ```env DEBUG=xx ./proc```
* 支持参数搜集 ```cat a.txt b.txt```，可以把```a.txt, b.txt```散装成员归归类，收集到你指定的结构体成员里
* 支持短选项```proc -d``` 或者长选项```proc --debug```不在话下
* posix风格命令行支持，支持命令组合```ls -ltr```是```ls -l -t -r```简写形式，方便实现普通posix 标准命令
* 子命令支持，方便实现git风格子命令```git add ```，简洁的子命令注册方式，只要会写结构提就行，3,4,5到无穷尽子命令也支持，只要你喜欢，用上clop就可以实现
* 默认值支持```default:"1"```，支持多种数据类型，让你省去类型转换的烦恼
* 贴心的重复命令报错
* 严格的短选项，长选项报错。避免二义性选项诞生
* 效验模式支持，不需要写一堆的```if x!= "" ``` or ```if y!=0```浪费青春的代码

## 内容
- [Installation](#Installation)
- [Quick start](#quick-start)
- [example](#example)
	- [1. How to use required tags](#required-flag)
	- [2.Support environment variables](#support-environment-variables)
	- [3.Set default value](#set-default-value)
	- [4. How to implement git style commands](#subcommand)
- [Implementing linux command options](#Implementing-linux-command-options)
	- [cat](#cat)
## Installation
```
go get github.com/guonaihong/clop
```

## Quick start
```go
package main

import (
	"fmt"
	"github.com/guonaihong/clop"
)

type Hello struct {
	File string `clop:"-f; --file" usage:"file"`
}

func main() {

	h := Hello{}
	clop.Bind(&h)
	fmt.Printf("%#v\n", h)
}
// ./one -f test
// main.Hello{File:"test"}
// ./one --file test
// main.Hello{File:"test"}

```
## example
### required flag
```go
 package main

import (
    "fmt"
    "github.com/guonaihong/clop"
)

type curl struct {
    Url string `clop:"-u; --url" usage:"url" valid:"required"`
}

func main() {

    c := curl{}
    clop.Bind(&c)
}

```
#### set default value
可以使用default tag设置默认值，普通类型直接写，复合类型用json表示
```go
package main

import (
    "fmt"
    "github.com/guonaihong/clop"
)

type defaultExample struct {
    Int          int       `default:"1"`
    Float64      float64   `default:"3.64"`
    Float32      float32   `default:"3.32"`
    SliceString  []string  `default:"[\"one\", \"two\"]"`
    SliceInt     []int     `default:"[1,2,3,4,5]"`
    SliceFloat64 []float64 `default:"[1.1,2.2,3.3,4.4,5.5]"`
}

func main() {
    de := defaultExample{}
    clop.Bind(&de)
    fmt.Printf("%v\n", de) 
}
// run
//         ./use_def
// output:
//         {1 3.64 3.32 [one two] [1 2 3 4 5] [1.1 2.2 3.3 4.4 5.5]}
```
### Support environment variables
```go
// file name use_env.go
package main

import (
	"fmt"
	"github.com/guonaihong/clop"
)

type env struct {
	OmpNumThread string `clop:"env=omp_num_thread" usage:"omp num thread"`
	Path         string `clop:"env=XPATH" usage:"xpath"`
	Max          int    `clop:"env=MAX" usage:"max thread"`
}

func main() {
	e := env{}
	clop.Bind(&e)
	fmt.Printf("%#v\n", e)
}
// run
// env XPATH=`pwd` omp_num_thread=3 MAX=4 ./use_env 
// output
// main.env{OmpNumThread:"3", Path:"/home/guo", Max:4}
```
### subcommand
```go
package main

import (
	"fmt"
	"github.com/guonaihong/clop"
)

type add struct {
	All      bool     `clop:"-A; --all" usage:"add changes from all tracked and untracked files"`
	Force    bool     `clop:"-f; --force" usage:"allow adding otherwise ignored files"`
	Pathspec []string `clop:"args=pathspec"`
}

type mv struct {
	Force bool `clop:"-f; --force" usage:"allow adding otherwise ignored files"`
}

type git struct {
	Add add `clop:"subcommand=add" usage:"Add file contents to the index"`
	Mv  mv  `clop:"subcommand=mv" usage:"Move or rename a file, a directory, or a symlink"`
}

func main() {
	g := git{}
	clop.Bind(&g)
	fmt.Printf("git:%#v\n", g)
	fmt.Printf("git:set mv(%t) or set add(%t)\n", clop.IsSetSubcommand("mv"), clop.IsSetSubcommand("add"))
}
// run:
// ./git add -f

// output:
// git:main.git{Add:main.add{All:false, Force:true, Pathspec:[]string(nil)}, Mv:main.mv{Force:false}}
// git:set mv(false) or set add(true)

```
## Implementing linux command options
### cat
```go
package main

import (
	"fmt"
	"github.com/guonaihong/clop"
)

type cat struct {
	NumberNonblank bool `clop:"-c;--number-nonblank" 
	                     usage:"number nonempty output lines, overrides"`

	ShowEnds bool `clop:"-E;--show-ends" 
	               usage:"display $ at end of each line"`

	Number bool `clop:"-n;--number" 
	             usage:"number all output lines"`

	SqueezeBlank bool `clop:"-s;--squeeze-blank" 
	                   usage:"suppress repeated empty output lines"`

	ShowTab bool `clop:"-T;--show-tabs" 
	              usage:"display TAB characters as ^I"`

	ShowNonprinting bool `clop:"-v;--show-nonprinting" 
	                      usage:"use ^ and M- notation, except for LFD and TAB" `

	Files []string `clop:"args=files"`
}

func main() {

	c := cat{}
	err := clop.Bind(&c)

	fmt.Printf("%#v, %s\n", c, err)
}

/*
Usage:
    ./cat [Flags] <files> 

Flags:
    -E,--show-ends           display $ at end of each line 
    -T,--show-tabs           display TAB characters as ^I 
    -c,--number-nonblank     number nonempty output lines, overrides 
    -n,--number              number all output lines 
    -s,--squeeze-blank       suppress repeated empty output lines 
    -v,--show-nonprinting    use ^ and M- notation, except for LFD and TAB 

Args:
    <files>
*/
```
