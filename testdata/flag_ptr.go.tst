package main

import "flag"

var (
	cmdString    string
	cmdInt       int
	flumeCfgFile string
	logLevel     string
)

func main() {
	flag.StringVar(&cmdString, "cs", "aug", "example cmd string")
	flag.IntVar(&cmdInt, "ci", 29, "example cmd int")
	flag.StringVar(&flumeCfgFile, "c", "flume.json", "flume config file name")
	flag.StringVar(&logLevel, "l", "i", "d debug, i info, w warning, e error")
	flag.Parse()
}
