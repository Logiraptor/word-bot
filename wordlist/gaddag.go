package wordlist

import (
	"fmt"
	"io"

	"github.com/Logiraptor/word-bot/core"
)

const reverseToken = ('z' - 'a') + 1

type Gaddag struct {
	nodes    [27]*Gaddag
	terminal bool
}

func NewGaddag() *Gaddag {
	return &Gaddag{}
}

func (g *Gaddag) AddWord(word string) {
	contents := []rune(word)
	for i := range contents {
		g.insertLinearString(contents[:i], contents[i:])
	}
}

func (g *Gaddag) insertLinearString(start, end []rune) {
	current := g
	for _, r := range end {
		if current.nodes[r-'a'] == nil {
			current.nodes[r-'a'] = NewGaddag()
		}
		current = current.nodes[r-'a']
	}

	if current.nodes[reverseToken] == nil {
		current.nodes[reverseToken] = NewGaddag()
	}
	current = current.nodes[reverseToken]

	for i := len(start) - 1; i >= 0; i-- {
		r := start[i]
		if current.nodes[r-'a'] == nil {
			current.nodes[r-'a'] = NewGaddag()
		}
		current = current.nodes[r-'a']
	}
	current.terminal = true
}

func (g *Gaddag) CanBranch(l core.Tile) bool {
	return g.nodes[l.ToLetter()] != nil
}

func (g *Gaddag) Branch(l core.Tile) *Gaddag {
	return g.nodes[l.ToLetter()]
}

func (g *Gaddag) CanReverse() bool {
	return g.nodes[reverseToken] != nil
}

func (g *Gaddag) Reverse() *Gaddag {
	return g.nodes[reverseToken]
}

func (g *Gaddag) IsTerminal() bool {
	return g.terminal
}

func (g *Gaddag) DumpOptions() []string {
	output := []string{}
	if g.IsTerminal() {
		output = append(output, ".")
	}
	for i, n := range g.nodes {
		if n == nil {
			continue
		}
		r := string(core.Letter(i).ToRune())
		if i == reverseToken {
			r = "#"
		}
		subStrings := n.DumpOptions()
		for _, s := range subStrings {
			output = append(output, r+s)
		}
	}
	return output
}

func (g *Gaddag) DumpToDot(wr io.Writer) {
	fmt.Fprintln(wr, "digraph {")
	g.dumpToDot(wr)
	fmt.Fprintln(wr, "}")
}

func (g *Gaddag) dumpToDot(wr io.Writer) {
	for i, n := range g.nodes {
		if n == nil {
			continue
		}
		r := string(core.Letter(i).ToRune())
		if i == reverseToken {
			r = "#"
		}

		color := "white"
		if n.IsTerminal() {
			color = "red"
		} else if n.CanReverse() {
			color = "blue"
		}

		fmt.Fprintf(wr, "\"%p\" -> \"%p\" [label=\"%s\"];\n", g, n, r)
		fmt.Fprintf(wr, "\"%p\" [label=\"%s\" color=\"%s\"];\n", n, "__", color)
		n.dumpToDot(wr)
	}
}
