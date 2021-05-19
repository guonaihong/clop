package main

import "flag"

func main() {
	str := flag.String("string", "1234", "string opt usage")
	size := flag.Int("int", 0, "int opt usage")
	flag.Parse()
}
