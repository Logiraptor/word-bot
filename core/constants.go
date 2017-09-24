package core

// Direction constants
const (
	Vertical   Direction = true
	Horizontal           = !Vertical
)

// Shorthands for bonuses
const (
	xx Bonus = iota
	DW
	TW
	DL
	TL
)

// Longer names for bonus spaces
const (
	None         = xx
	DoubleWord   = DW
	TripleWord   = TW
	DoubleLetter = DL
	TripleLetter = TL
)

var normalBonus = [...][15]Bonus{
	{TW, xx, xx, DL, xx, xx, xx, TW, xx, xx, xx, DL, xx, xx, TW},
	{xx, DW, xx, xx, xx, TL, xx, xx, xx, TL, xx, xx, xx, DW, xx},
	{xx, xx, DW, xx, xx, xx, DL, xx, DL, xx, xx, xx, DW, xx, xx},
	{DL, xx, xx, DW, xx, xx, xx, DL, xx, xx, xx, DW, xx, xx, DL},
	{xx, xx, xx, xx, DW, xx, xx, xx, xx, xx, DW, xx, xx, xx, xx},
	{xx, TL, xx, xx, xx, TL, xx, xx, xx, TL, xx, xx, xx, TL, xx},
	{xx, xx, DL, xx, xx, xx, DL, xx, DL, xx, xx, xx, DL, xx, xx},
	{TW, xx, xx, DL, xx, xx, xx, DW, xx, xx, xx, DL, xx, xx, TW},
	{xx, xx, DL, xx, xx, xx, DL, xx, DL, xx, xx, xx, DL, xx, xx},
	{xx, TL, xx, xx, xx, TL, xx, xx, xx, TL, xx, xx, xx, TL, xx},
	{xx, xx, xx, xx, DW, xx, xx, xx, xx, xx, DW, xx, xx, xx, xx},
	{DL, xx, xx, DW, xx, xx, xx, DL, xx, xx, xx, DW, xx, xx, DL},
	{xx, xx, DW, xx, xx, xx, DL, xx, DL, xx, xx, xx, DW, xx, xx},
	{xx, DW, xx, xx, xx, TL, xx, xx, xx, TL, xx, xx, xx, DW, xx},
	{TW, xx, xx, DL, xx, xx, xx, TW, xx, xx, xx, DL, xx, xx, TW},
}
