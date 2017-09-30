package core

import (
	"reflect"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
)

func TestConsumableBagCanConsume(t *testing.T) {
	err := quick.Check(func(i byte) bool {
		idx := int(i) % len(allTiles)
		originalBag := NewConsumableBag()
		changedBag := originalBag.Consume(idx)
		return assert.True(t, originalBag.CanConsume(idx), "New Bag should be able to consume tile %d", idx) &&
			assert.False(t, changedBag.CanConsume(idx), "Changed bag must not be able to consume tile %d", idx) &&
			assert.True(t, originalBag.CanConsume(idx), "Original bag should be unchanged")
	}, nil)
	assert.NoError(t, err)
}

func TestConsumableBagNumTiles(t *testing.T) {
	err := quick.Check(func(i byte) bool {
		idx := int(i) % len(allTiles)
		originalBag := NewConsumableBag()
		for i := 0; i < idx; i++ {
			originalBag = originalBag.Consume(i)
		}
		return assert.Equal(t, len(allTiles)-idx, originalBag.Count(), "Consuming %d tiles should leave %d in the bag", idx, len(allTiles)-idx)
	}, nil)
	assert.NoError(t, err)
}

func TestConsumableBagFillRack(t *testing.T) {
	err := quick.Check(func(i byte) bool {
		idx := int(i) % len(allTiles)
		originalBag := NewConsumableBag()
		rack := NewConsumableRack(nil)
		originalBag, rack.Rack = originalBag.FillRack(rack.Rack, idx)
		return assert.Equal(t, len(allTiles)-idx, originalBag.Count(), "Filling a rack with %d tiles should leave %d in the bag", idx, len(allTiles)-idx)
	}, nil)
	assert.NoError(t, err)
}

func TestConsumableBagShuffling(t *testing.T) {
	err := quick.Check(func(i byte) bool {
		idx := int(i) % len(allTiles)
		if idx == 0 {
			idx = 1
		}
		originalBag := NewConsumableBag()
		shuffledBag := originalBag.Shuffle()

		originalRack := NewConsumableRack(nil)
		_, originalRack.Rack = originalBag.FillRack(originalRack.Rack, idx)

		shuffledRack := NewConsumableRack(nil)
		_, shuffledRack.Rack = shuffledBag.FillRack(shuffledRack.Rack, idx)

		return assert.Equal(t, originalBag.Count(), shuffledBag.Count(), "Shuffling should not change count") &&
			assert.False(t, reflect.DeepEqual(originalRack.Rack, shuffledRack.Rack),
				"Rack filled from shuffled bag should be different, got %q / %q",
				tiles2String(originalRack.Rack),
				tiles2String(shuffledRack.Rack))
	}, nil)
	assert.NoError(t, err)
}
