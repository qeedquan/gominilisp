package main

import "log"

func main() {
	log.SetFlags(0)
	ip := NewInterp()
	ip.Repl()
}
