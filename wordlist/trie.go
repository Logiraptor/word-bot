package wordlist

import (
	"bytes"
	"fmt"

	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/definitions"
)

//go:generate go-bindata -pkg wordlist -o ./words.go ./words.txt

func MakeDefaultWordList() *Trie {
	builder := NewTrieBuilder(151434)
	buf := MustAsset("words.txt")
	err := definitions.LoadDefinitionsReader(bytes.NewReader(buf), builder)
	if err != nil {
		panic(fmt.Sprintf("Cannot load embedded word list file"))
	}
	return builder.Build()
}

func MakeDefaultWordListGaddag() *Gaddag {
	builder := NewGaddag()
	buf := MustAsset("words.txt")
	err := definitions.LoadDefinitionsReader(bytes.NewReader(buf), builder)
	if err != nil {
		panic(fmt.Sprintf("Cannot load embedded word list file"))
	}
	return builder
}

type Trie struct {
	nodes    [26]*Trie
	terminal bool
}

func NewTrie() *Trie {
	return &Trie{}
}

func (t *Trie) Contains(word core.Word) bool {
	current := t
	for _, letter := range word {
		if current.nodes[letter] == nil {
			return false
		}
		current = current.nodes[letter]
	}
	return current.terminal
}

func (t *Trie) IsTerminal() bool {
	return t.terminal
}

func (t *Trie) CanBranch(tile core.Tile) (*Trie, bool) {
	next := t.nodes[tile.ToLetter()]
	return next, next != nil
}

func (t *Trie) AddWord(word string) {
	current := t
	for _, letter := range word {
		i := core.Rune2Letter(letter)
		if current.nodes[i] == nil {
			current.nodes[i] = NewTrie()
		}

		current = current.nodes[i]
	}

	current.terminal = true
}

type TrieBuilder struct {
	nodes []Trie
	ptr   int
}

func NewTrieBuilder(size int) *TrieBuilder {
	nodes := make([]Trie, size)
	nodes = append(nodes, Trie{})
	return &TrieBuilder{
		nodes: nodes,
		ptr:   1,
	}
}

func (t *TrieBuilder) allocateTrie() *Trie {
	t.ptr++
	if t.ptr >= len(t.nodes) {
		t.nodes = append(t.nodes, Trie{})
		t.nodes = t.nodes[:cap(t.nodes)]
	}
	return &t.nodes[t.ptr]
}

func (t *TrieBuilder) AddWord(word string) {
	current := &t.nodes[0]
	for _, letter := range word {
		i := core.Rune2Letter(letter)
		if current.nodes[i] == nil {
			current.nodes[i] = t.allocateTrie()
		}
		current = current.nodes[i]
	}

	current.terminal = true
}

func (t *TrieBuilder) Build() *Trie {
	return &t.nodes[0]
}
