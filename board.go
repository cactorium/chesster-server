package chesster

type PieceType int

const (
	InvalidPiece PieceType = iota
	Pawn
	Rook
	Knight
	Bishop
	King
	Queen
)

type Piece struct {
	X        int
	Y        int
	Type     PieceType
	HasMoved bool
}

type BoardState int

const (
	// white's turn to move
	WhiteMove BoardState = iota
	// black's turn to move
	BlackMove
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

type Board struct {
	Moves []Move
	White []Piece
	Black []Piece
	// pieces captured by white
	CapturedWhite []PieceType
	// pieces captured by black
	CapturedBlack []PieceType
	State         BoardState
	// white's king is currently threatened
	WhiteCheck bool
	// black's king is currently threatened
	BlackCheck bool
	// white asked for draw
	WhiteDrawAsk bool
	// black asked for draw
	BlackDrawAsk      bool
	MovesSinceCapture int
	// -1 if not marked, set if white moved a pawn two spaces last turn, set to the pawn's location
	WhiteEnPassant int
	// -1 if not marked, set if black moved a pawn two spaces last turn, set to the pawn's location
	BlackEnPassant int
}

func NewBoard() Board {
	// TODO: generate a fresh board
	panic("unimplemented!")
	return Board{
		White: []Piece{
			Piece{0, 0, Rook, false},
		},
		Black:             []Piece{},
		CapturedWhite:     []PieceType{},
		CapturedBlack:     []PieceType{},
		State:             WhiteMove,
		WhiteCheck:        false,
		BlackCheck:        false,
		WhiteDrawAsk:      false,
		BlackDrawAsk:      false,
		MovesSinceCapture: 0,
		WhiteEnPassant:    -1,
		BlackEnPassant:    -1,
	}
}

func (b *Board) Clone() Board {
	newBoard := Board{
		Moves:             make([]Move, len(b.Moves)),
		White:             make([]Piece, len(b.White)),
		Black:             make([]Piece, len(b.Black)),
		CapturedWhite:     make([]PieceType, len(b.CapturedWhite)),
		CapturedBlack:     make([]PieceType, len(b.CapturedBlack)),
		State:             b.State,
		WhiteCheck:        b.WhiteCheck,
		BlackCheck:        b.BlackCheck,
		WhiteDrawAsk:      b.WhiteDrawAsk,
		BlackDrawAsk:      b.BlackDrawAsk,
		MovesSinceCapture: b.MovesSinceCapture,
		WhiteEnPassant:    b.WhiteEnPassant,
		BlackEnPassant:    b.BlackEnPassant,
	}
	copy(newBoard.Moves, b.Moves)
	copy(newBoard.White, b.White)
	copy(newBoard.Black, b.Black)
	copy(newBoard.CapturedWhite, b.CapturedWhite)
	copy(newBoard.CapturedBlack, b.CapturedBlack)
	return newBoard
}

func (b *Board) IsMove(s Side) bool {
	return (s == White && b.State == WhiteMove) || (s == Black && b.State == BlackMove)
}

func (b *Board) GameEnded() bool {
	return b.State != WhiteMove && b.State != BlackMove
}

func (b *Board) WhiteWon() bool {
	return b.State == BlackCheckmate
}

func (b *Board) BlackWon() bool {
	return b.State == WhiteCheckmate
}

func (b *Board) Draw() bool {
	return (b.State == DrawAgreed) || (b.State == Draw50Moves) || (b.State == Draw3Fold)
}

func (b *Board) OfferDraw(s Side) {
	if s == White {
		b.WhiteDrawAsk = true
	} else {
		b.BlackDrawAsk = true
	}

	if b.WhiteDrawAsk && b.BlackDrawAsk {
		b.State = DrawAgreed
	}
}

func (b *Board) DoMove(m Move) {
	// TODO: do move and update state as needed
	// TODO: append to movelist
	// TODO: check to see if it's en passant
	// TODO: update captured pieces
	// TODO: update check states
	// TODO: reset draw ask
	// TODO: update moves since capture
	// TODO: update en passant states
	// TODO: check for 50 move draw
	// TODO: check for 3 fold repetition
	panic("unimplemented!")
}

func (b *Board) IsValid() bool {
	// make sure there are the appropriate number of pieces on the board
	if (len(b.White)+len(b.CapturedWhite)) > 8 || (len(b.Black)+len(b.CapturedBlack)) > 8 {
		return false
	}
	// make sure all the pieces are in bounds
	for _, piece := range b.White {
		if piece.X < 0 || piece.X > 8 || piece.Y < 0 || piece.Y > 8 {
			return false
		}
	}
	for _, piece := range b.Black {
		if piece.X < 0 || piece.X > 8 || piece.Y < 0 || piece.Y > 8 {
			return false
		}
	}

	// make sure both kings aren't in check
	if b.WhiteCheck && b.BlackCheck {
		return false
	}
	// make sure en passant isn't set for both players
	if b.WhiteEnPassant != -1 && b.BlackEnPassant != -1 {
		return false
	}
	return true
}

type Side int

const (
	White Side = iota
	Black
)

type Move struct {
	Side        Side
	Start       Piece
	End         Piece
	IsPromotion bool
	// if the move is a castle; overrides all the above attributes
	IsCastle bool
	// true is kingside, false if queenside; ignored if not castle
	IsKingsideCastle bool
}

type InvalidMoveReason int

const (
	MoveOkay InvalidMoveReason = iota
	GameEnded
	WrongSide
	OutOfBounds
	PieceNotFound
	SpaceOccupied
	PieceCantMoveThatWay
	TypeChange
	StillInCheck
	OnlyOneKing
	CantCastle
)

func (b *Board) TryMove(m Move) (*Board, InvalidMoveReason) {
	// sanity checks
	// make sure the game is still in session
	if b.State != WhiteMove && b.State != BlackMove {
		return nil, GameEnded
	}
	// make sure the move is from the right player
	if (b.State == WhiteMove && m.Side != White) || (b.State == BlackMove && m.Side != Black) {
		return nil, WrongSide
	}

	// make sure the move ends within bounds
	if m.End.X > 8 || m.End.X < 0 || m.End.Y > 8 || m.End.X < 0 {
		return nil, OutOfBounds
	}

	var toCheck []Piece
	if m.Side == White {
		toCheck = b.White
	} else {
		toCheck = b.Black
	}

	// make sure the piece exists
	pieceFound := false
	for _, piece := range toCheck {
		if (piece.X == m.Start.X) && (piece.Y == m.End.Y) {
			pieceFound = true
			break
		}
	}
	if !pieceFound {
		return nil, PieceNotFound
	}

	// make sure the move lands on top of a piece of opposite color or an empty square (basically doesn't land on top of one of their own pieces)
	for _, piece := range toCheck {
		if (piece.X == m.End.X) && (piece.Y == m.End.Y) {
			return nil, SpaceOccupied
		}
	}

	// TODO: make sure the move is valid for the piece type
	// TODO: make sure the move ends with the same piece type as it started (assuming it is not a promotion)
	// TODO: if the move is a promotion, the end piece type is not a king
	// TODO: if the move is a kind of castle, make sure all the conditions are met:
	//   king and appropriate rook haven't moved
	//   the space between them is clear
	//   the king is not in check
	//   after the move, the king is not in check
	var isInCheck bool
	if b.State == WhiteMove {
		isInCheck = b.WhiteCheck
	} else {
		isInCheck = b.BlackCheck
	}
	afterMove := b.Clone()
	afterMove.DoMove(m)
	// see if the move leaves the player who moved in check
	if isInCheck {
		if (m.Side == White && afterMove.WhiteCheck) || (m.Side == Black && afterMove.BlackCheck) {
			return nil, StillInCheck
		}
	}

	// else it's good
	return &afterMove, MoveOkay
}

func (m Move) String() string {
	// TODO: write move as standard chess notation
	panic("unimplemented!")
}
