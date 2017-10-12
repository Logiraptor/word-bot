package wordlist

import (
	"testing"

	"github.com/Logiraptor/word-bot/core"
	"github.com/stretchr/testify/assert"
)

func TestGaddagCanBranch(t *testing.T) {
	gaddag := NewGaddag()
	gaddag.AddWord("abc")
	a := core.Rune2Letter('a')
	b := core.Rune2Letter('b')
	c := core.Rune2Letter('c')

	assert.True(t, gaddag.CanBranch(a))
	assert.True(t, gaddag.CanBranch(b))
	assert.True(t, gaddag.CanBranch(c))
}

func TestGaddagBranch(t *testing.T) {
	gaddag := NewGaddag()
	gaddag.AddWord("abc")
	a := core.Rune2Letter('a')
	b := core.Rune2Letter('b')
	c := core.Rune2Letter('c')

	aBranch := gaddag.Branch(a)
	assert.False(t, aBranch.CanBranch(a))
	assert.True(t, aBranch.CanBranch(b))
	assert.False(t, aBranch.CanBranch(c))
}

func TestGaddagCanReverse(t *testing.T) {
	gaddag := NewGaddag()
	gaddag.AddWord("abc")
	a := core.Rune2Letter('a')
	b := core.Rune2Letter('b')
	c := core.Rune2Letter('c')

	aBranch := gaddag.Branch(a)
	bBranch := gaddag.Branch(b)
	cBranch := gaddag.Branch(c)
	assert.False(t, aBranch.CanReverse())
	assert.False(t, bBranch.CanReverse())
	assert.True(t, cBranch.CanReverse())
}

func TestGaddagReverse(t *testing.T) {
	gaddag := NewGaddag()
	gaddag.AddWord("abc")
	a := core.Rune2Letter('a')
	b := core.Rune2Letter('b')
	c := core.Rune2Letter('c')

	cBranch := gaddag.Branch(c).Reverse()
	assert.False(t, cBranch.CanBranch(a))
	assert.True(t, cBranch.CanBranch(b))
	assert.False(t, cBranch.CanBranch(c))
}

func TestGaddagTerminal(t *testing.T) {
	gaddag := NewGaddag()
	gaddag.AddWord("abc")
	a := core.Rune2Letter('a')
	b := core.Rune2Letter('b')
	c := core.Rune2Letter('c')

	t1 := gaddag.Branch(c)
	t2 := t1.Reverse()
	t3 := t2.Branch(b)
	t4 := t3.Branch(a)

	assert.False(t, t1.IsTerminal())
	assert.False(t, t2.IsTerminal())
	assert.False(t, t3.IsTerminal())
	assert.True(t, t4.IsTerminal())
}
