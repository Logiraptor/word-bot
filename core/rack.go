package core

// ConsumableRack manages a rack of tiles and allows efficient consumption of tiles
type ConsumableRack struct {
	Rack     []Tile
	consumed int
}

func NewConsumableRack(tiles []Tile) ConsumableRack {
	return ConsumableRack{
		Rack:     tiles,
		consumed: 0,
	}
}

// Consume uses up a tile in the rack
func (c ConsumableRack) Consume(i int) ConsumableRack {
	return ConsumableRack{
		Rack:     c.Rack,
		consumed: c.consumed | (1 << uint(i)),
	}
}

// CanConsume returns true if the ith tile is available to use
func (c ConsumableRack) CanConsume(i int) bool {
	return c.consumed&(1<<uint(i)) == 0
}
