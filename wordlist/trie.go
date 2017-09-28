package wordlist

import (
	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
)

type Trie struct {
	nodes    [26]*Trie
	terminal bool
}

var _ ai.WordTree = NewTrie()

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

func (t *Trie) CanBranch(tile core.Tile) (ai.WordTree, bool) {
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

type TrieInConstruction struct {
	nodes    [26]int
	terminal bool
}

type TrieBuilder struct {
	nodes []TrieInConstruction
}

func NewTrieBuilder() *TrieBuilder {
	return &TrieBuilder{
		nodes: []TrieInConstruction{
			{},
		},
	}
}

func (t *TrieBuilder) AddWord(word string) {
	current := &t.nodes[0]
	for _, letter := range word {
		i := core.Rune2Letter(letter)
		if current.nodes[i] == 0 {
			current.nodes[i] = len(t.nodes)
			t.nodes = append(t.nodes, TrieInConstruction{})
		}
		current = &t.nodes[current.nodes[i]]
	}

	current.terminal = true
}

func (t *TrieBuilder) Build() *Trie {
	finalTries := make([]Trie, len(t.nodes))
	for i, node := range t.nodes {
		for j, link := range node.nodes {
			if link != 0 {
				finalTries[i].nodes[j] = &finalTries[link]
			}
		}
		finalTries[i].terminal = node.terminal
	}
	return &finalTries[0]
}
