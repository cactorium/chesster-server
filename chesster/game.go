package chesster

type GameState int

const (
	InPlay GameState = iota
	// white checkmated
	WhiteCheckmate
	// black checkmated
	BlackCheckmate
	// white stallmated; black move caused stalemate
	WhiteStallmate
	// black stallmated; white move caused stalemate
	BlackStallmate
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

func (g *Game) DoMove(m Move) {
	// do move and update board state as needed
	g.Board.CommitMove(m)
	// append to movelist
	g.Moves = append(g.Moves, m)
	// TODO: check to see if it's en passant
	// TODO: update en passant states
	// TODO: update captured pieces
	// TODO: check for check
	// TODO: reset draw ask
	// TODO: update moves since capture
	// TODO: check for 50 move draw
	// TODO: check for 3 fold repetition
	panic("unimplemented!")
}

func (g *Game) IsValid() bool {
	if !g.Board.IsValid() {
		return false
	}

	// make sure both kings aren't in check
	if g.WhiteCheck && g.BlackCheck {
		return false
	}
	return true
}
