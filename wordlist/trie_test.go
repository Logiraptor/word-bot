package wordlist

import (
	"testing"

	"github.com/Logiraptor/word-bot/core"
)

func loadWordsDirect() *Trie {
	return MakeDefaultWordList()
}

var blankA = core.Rune2Letter('a').ToTile(false)
var blankZ = core.Rune2Letter('z').ToTile(false)

func dfs(t *Trie) {
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
	trie := loadWordsDirect()

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
	trie := loadWordsDirect()
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
		loadWordsDirect()
	}
}
