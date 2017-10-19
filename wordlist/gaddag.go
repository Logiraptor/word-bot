package wordlist

import "github.com/Logiraptor/word-bot/core"

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
	for i, n := range g.nodes {
		if n == nil {
			continue
		}
		r := string(core.Letter(i).ToRune())
		if i == reverseToken {
			r = "#"
		}
		if n.IsTerminal() {
			output = append(output, r)
		}
		subStrings := n.DumpOptions()
		for _, s := range subStrings {
			output = append(subStrings, r+s)
		}
	}
	return output
}
