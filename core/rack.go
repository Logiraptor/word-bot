package core

func NewConsumableRack(tiles []Tile) TileSet {
	s := NewEmptyTileSet()
	for _, t := range tiles {
		s.Add(t)
	}
	return s
}
