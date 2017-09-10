package main

type Trie struct {
	nodes    [26]*Trie
	terminal bool
}

func NewTrie() *Trie {
	return &Trie{}
}

func (t *Trie) Contains(word string) bool {
	current := t
	for _, letter := range word {
		i := rune2Letter(letter)
		if current.nodes[i] == nil {
			return false
		}
		current = current.nodes[i]
	}
	return current.terminal
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
