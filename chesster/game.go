package chesster

type GameState int

const (
	InPlay GameState = iota
	// black checkmated; white wins!
	WhiteCheckmate
	// white checkmated; black wins!
	BlackCheckmate
	// white stalemated; black move caused stalemate
	WhiteStalemate
	// black stalemated; white move caused stalemate
	BlackStalemate
	// white resigned
	WhiteResigned
	// black resigned
	BlackResigned
	// both sides agreed to draw
	DrawAgreed
	// FIDE rule; 50 moves since capture or pawn move
	Draw50Moves
	// FIDE rule; threefold repetition
	Draw3Fold
)

type Game struct {
	Moves []Move
	Board Board
	State GameState
	// white's king is currently threatened
	WhiteCheck bool
	// black's king is currently threatened
	BlackCheck bool
	// white asked for draw
	WhiteDrawAsk bool
	// black asked for draw
	BlackDrawAsk      bool
	MovesSinceCapture int
}

func NewGame() Game {
	return Game{
		Moves:             []Move{},
		Board:             NewBoard(),
		State:             InPlay,
		WhiteCheck:        false,
		BlackCheck:        false,
		WhiteDrawAsk:      false,
		BlackDrawAsk:      false,
		MovesSinceCapture: 0,
	}
}

func (g *Game) Clone() Game {
	newGame := Game{
		Moves:             make([]Move, len(g.Moves)),
		Board:             g.Board.Clone(),
		State:             g.State,
		WhiteCheck:        g.WhiteCheck,
		BlackCheck:        g.BlackCheck,
		WhiteDrawAsk:      g.WhiteDrawAsk,
		BlackDrawAsk:      g.BlackDrawAsk,
		MovesSinceCapture: g.MovesSinceCapture,
	}
	copy(newGame.Moves, g.Moves)
	return newGame
}

func (g *Game) GameEnded() bool {
	return g.State != InPlay
}

func (g *Game) WhiteWon() bool {
	return g.State == BlackCheckmate
}

func (g *Game) BlackWon() bool {
	return g.State == WhiteCheckmate
}

func (g *Game) Draw() bool {
	return (g.State == DrawAgreed) || (g.State == Draw50Moves) || (g.State == Draw3Fold)
}

func (g *Game) OfferDraw(s Side) {
	if s == White {
		g.WhiteDrawAsk = true
	} else {
		g.BlackDrawAsk = true
	}

	if g.WhiteDrawAsk && g.BlackDrawAsk {
		g.State = DrawAgreed
	}
}

func (g *Game) RescindDraw(s Side) {
	if g.State != DrawAgreed {
		if s == White {
			g.WhiteDrawAsk = false
		} else {
			g.BlackDrawAsk = false
		}
	}
}

func (g *Game) DoMove(m Move) (b bool, r InvalidMoveReason) {
	ocl := len(g.Board.Captured)
	// do move and update board state as needed
	if b, r = g.Board.TryMove(m); !b {
		return
	}

	// append to movelist
	g.Moves = append(g.Moves, m)

	// check for check
	g.BlackCheck = g.Board.InCheck(Black)
	g.WhiteCheck = g.Board.InCheck(White)

	// check for checkmate and stalemate
	if g.Board.IsMove(White) {
		if g.Board.Stalemate(White) {
			g.State = WhiteStalemate
			return
		}
		if g.Board.Checkmate(White) {
			g.State = BlackCheckmate
			return
		}
	} else {
		if g.Board.Stalemate(Black) {
			g.State = BlackStalemate
			return
		}
		if g.Board.Checkmate(Black) {
			g.State = WhiteCheckmate
			return
		}
	}

	// update moves since capture
	if len(g.Board.Captured) > ocl {
		g.MovesSinceCapture = 0
	} else {
		g.MovesSinceCapture += 1
	}

	// check for 50 move draw
	if g.MovesSinceCapture >= 50 {
		g.State = Draw50Moves
	}

	// TODO: check for 3 fold repetition
	panic("unimplemented!")
	return
}
