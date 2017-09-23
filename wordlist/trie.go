package wordlist

import (
	"word-bot/ai"
	"word-bot/core"
)

type Trie struct {
	nodes    [26]*Trie
	terminal bool
}

var _ ai.WordList = NewTrie()

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

func (t *Trie) CanBranch(tile core.Tile) (ai.WordList, bool) {
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
