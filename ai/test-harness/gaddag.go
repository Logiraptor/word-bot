package main

import (
	"os"

	"github.com/Logiraptor/word-bot/wordlist"
)

func main() {

	f, err := os.Create("out.dot")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	g := wordlist.NewGaddag()
	g.AddWord("shaded")
	g.AddWord("doggo")
	g.AddWord("oar")
	g.DumpToDot(f)

}
