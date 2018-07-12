package chesster

import (
	"fmt"
	"testing"
)

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
	ms = b.Pieces[1].GetPossibleMoves(&b)
	if len(ms) != 5 {
		t.Errorf("expected %d got %d", 5, len(ms))
	}

}

func TestValidTwoKingsPinnedBishop(t *testing.T) {
	b := Board{}
	b.Pieces = append(b.Pieces, Piece{4, 0, King, White, true})
	b.Pieces = append(b.Pieces, Piece{4, 6, King, Black, true})
	b.Pieces = append(b.Pieces, Piece{4, 3, Rook, White, true})
	b.Pieces = append(b.Pieces, Piece{4, 5, Bishop, Black, true})

	ms := b.Pieces[0].GetPossibleMoves(&b)
	if len(ms) != 5 {
		for _, m := range ms {
			fmt.Printf("%s\n", m.Notation(&b))
		}
		t.Errorf("expected %d got %d", 5, len(ms))
	}
	ms = b.Pieces[1].GetPossibleMoves(&b)
	if len(ms) != 7 {
		for _, m := range ms {
			fmt.Printf("%s\n", m.Notation(&b))
		}
		t.Errorf("expected %d got %d", 7, len(ms))
	}

	ms = b.Pieces[2].GetPossibleMoves(&b)
	if len(ms) != 11 {
		for _, m := range ms {
			fmt.Printf("%s\n", m.Notation(&b))
		}
		t.Errorf("expected %d got %d", 11, len(ms))
	}

	ms = b.Pieces[3].GetPossibleMoves(&b)
	if len(ms) != 0 {
		t.Errorf("expected %d got %d", 0, len(ms))
	}
}

func TestValidCastle(t *testing.T) {
	b := Board{}
	b.Pieces = append(b.Pieces, Piece{4, 0, King, White, false})
	b.Pieces = append(b.Pieces, Piece{4, 6, King, Black, true})
	b.Pieces = append(b.Pieces, Piece{7, 0, Rook, White, false})
	ms := b.Pieces[0].GetPossibleMoves(&b)
	if len(ms) != 6 {
		t.Errorf("expected %d got %d", 6, len(ms))
	}
	ms = b.Pieces[1].GetPossibleMoves(&b)
	if len(ms) != 8 {
		t.Errorf("expected %d got %d", 8, len(ms))
	}
	ms = b.Pieces[2].GetPossibleMoves(&b)
	if len(ms) != 9 {
		t.Errorf("expected %d got %d", 9, len(ms))
	}

	b.Pieces = append(b.Pieces, Piece{5, 7, Rook, Black, true})
	ms = b.Pieces[0].GetPossibleMoves(&b)
	if len(ms) != 3 {
		t.Errorf("expected %d got %d", 3, len(ms))
	}
	ms = b.Pieces[1].GetPossibleMoves(&b)
	if len(ms) != 7 {
		for _, m := range ms {
			fmt.Printf("%s\n", m.Notation(&b))
		}
		t.Errorf("expected %d got %d", 7, len(ms))
	}
	ms = b.Pieces[2].GetPossibleMoves(&b)
	if len(ms) != 9 {
		for _, m := range ms {
			fmt.Printf("%s\n", m.Notation(&b))
		}
		t.Errorf("expected %d got %d", 9, len(ms))
	}
	ms = b.Pieces[3].GetPossibleMoves(&b)
	if len(ms) != 14 {
		for _, m := range ms {
			fmt.Printf("%s\n", m.Notation(&b))
		}
		t.Errorf("expected %d got %d", 14, len(ms))
	}

}

// TODO check pawns
// lots of special cases there
