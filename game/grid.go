package game

type TurnResult uint8

type Grid struct {
	Data [13]uint8
}

const (
	TurnResultMiss = TurnResult(iota)
	TurnResultHit
)

func (g *Grid) Turn(x, y uint) TurnResult {
	pos := y*10 + x
	if pos > 99 {
		return TurnResultMiss
	}

	bytePos := pos / 8
	byte := g.Data[bytePos]

	bitPos := pos % 8
	if byte&(1<<bitPos) != 0 {
		return TurnResultHit
	}

	return TurnResultMiss
}
