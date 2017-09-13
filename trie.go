package main

type Trie struct {
	nodes    [26]*Trie
	terminal bool
}

func NewTrie() *Trie {
	return &Trie{}
}

func (t *Trie) Contains(word Word) bool {
	current := t
	for _, letter := range word {
		if current.nodes[letter] == nil {
			return false
		}
		current = current.nodes[letter]
	}
	return current.terminal
}

func (t *Trie) CanBranch(tile Tile) (*Trie, bool) {
	next := t.nodes[tile&LetterMask]
	return next, next != nil
}

func (t *Trie) AddWord(word string) {
	current := t
	for _, letter := range word {
		i := rune2Letter(letter)
		if current.nodes[i] == nil {
			current.nodes[i] = NewTrie()
		}

		current = current.nodes[i]
	}

	current.terminal = true
}
