package wordlist

import (
	"testing"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/definitions"

	"github.com/Logiraptor/word-bot/core"
)

func trashMemory() []int {
	var x []int
	for i := 0; i < 1000; i++ {
		x = make([]int, 10)
	}
	return x
}

func loadWordsDirect() *Trie {
	wordDB := NewTrie()
	definitions.LoadWords("../words.txt", wordDB)
	return wordDB
}

func loadWordsWithBuilder() *Trie {
	builder := NewTrieBuilder(151434)
	definitions.LoadWords("../words.txt", builder)
	return builder.Build()
}

var blankA = core.Rune2Letter('a').ToTile(false)
var blankZ = core.Rune2Letter('z').ToTile(false)

func dfs(t ai.WordTree) {
	for x := blankA; x <= blankZ; x++ {
		if next, ok := t.CanBranch(x); ok {
			dfs(next)
		}
		t.IsTerminal()
	}
}

func TestContainsDirect(t *testing.T) {
	trie := loadWordsDirect()

	if !trie.Contains(core.MakeWord("foot")) {
		t.Errorf("Trie did not contain foot!")
	}
	if trie.Contains(core.MakeWord("dugz")) {
		t.Errorf("Trie contained dugz!")
	}
}

func TestContainsBuilder(t *testing.T) {
	trie := loadWordsWithBuilder()

	if !trie.Contains(core.MakeWord("foot")) {
		t.Errorf("Trie did not contain foot!")
	}
	if trie.Contains(core.MakeWord("dugz")) {
		t.Errorf("Trie contained dugz!")
	}
}

func BenchmarkDFSDirect(b *testing.B) {
	trie := loadWordsDirect()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dfs(trie)
	}
}

func BenchmarkDFSBuilder(b *testing.B) {
	trie := loadWordsWithBuilder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dfs(trie)
	}
}

func BenchmarkLoadWordsDirect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loadWordsDirect()
	}
}

func BenchmarkLoadWordsWithBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loadWordsWithBuilder()
	}
}
