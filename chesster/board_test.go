package chesster

import "testing"

func TestValidMovesOneKing(t *testing.T) {
	b := Board{}
	b.Pieces = append(b.Pieces, Piece{4, 4, King, White, true})
	ms := b.Pieces[0].GetPossibleMoves(&b)
	if len(ms) != 8 {
		t.Errorf("expected %d got %d", 8, len(ms))
	}
}

func TestValidMovesTwoKing(t *testing.T) {
	b := Board{}
	b.Pieces = append(b.Pieces, Piece{4, 4, King, White, true})
	b.Pieces = append(b.Pieces, Piece{4, 6, King, Black, true})
	ms := b.Pieces[0].GetPossibleMoves(&b)
	if len(ms) != 5 {
		t.Errorf("expected %d got %d", 5, len(ms))
	}
}
