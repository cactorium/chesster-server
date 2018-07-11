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
	Side     Side
	HasMoved bool
}

func isInBounds(x, y int) bool {
	isOutOfBounds = (x < 0) || (x > 8) || (y < 0) || (y > 8)
	return !isOutOfBounds
}

// gets possible moves
func (p Piece) GetPossibleMoves(b *Board) []Move {
	isClear := func(x, y int) bool {
		return b.getPiece(x, y) == nil
	}
	hasEnemy := func(x, y int) bool {
		maybePiece := b.getPiece(x, y)
		if maybePiece == nil {
			return false
		}
		return maybePiece.Side == p.Side.Opposite()
	}
	raycast := func(ary []Move, dx, dy int) []Move {
		ret := ary
		x := p.X + dx
		y := p.Y + dy
		for isInBounds(x, y) && isClear(x, y) {
			newPiece := p
			newPiece.X = x
			newPiece.Y = y
			newPiece.HasMoved = true
			ret = append(ret, Move{p, newPiece})
		}
		if isInBounds(x, y) && hasEnemy(x, y) {
			newPiece := p
			newPiece.X = x
			newPiece.Y = y
			newPiece.HasMoved = true
			ret = append(ret, Move{p, newPiece})
		}
		return ret
	}

	moves := make([]Move, 0)
	switch p.Type {
	case InvalidPiece:
		return nil
	case Pawn:
		var forward int
		if p.Side == White {
			forward = 1
		} else {
			forward = -1
		}
		// add forward one square if it's not blocked
		if isInBounds(p.X, p.Y+forward) && isClear(p.X, p.Y+forward) {
			newPiece := p
			newPiece.Y = p.Y + forward
			newPiece.HasMoved = true
			if newPiece.Y == 0 || newPiece.Y == 7 {
				for pieceType := range []PieceType{Rook, Knight, Bishop, Queen} {
					newNewPiece := newPiece
					newNewPiece.Type = pieceType
					moves = append(moves, Move{p, newNewPiece, true})
				}
			} else {
				moves = append(moves, Move{p, newPiece})
			}
		}
		// add forward two squares if it's not blocked
		if !p.HasMoved {
			hasPieceTwoForward := b.getPiece(p.X, p.Y+2*forward) == nil
			if isInBounds(p.X, p.Y+2*forward) && isClear(p.X, p.Y+forward) && isClear(p.X, p.Y+2*forward) {
				newPiece := p
				newPiece.Y = p.Y + 2*forward
				newPiece.HasMoved = true
				moves = append(moves, Move{p, newPiece})
			}
		}
		// check possible captures
		if isInBounds(p.X-1, p.Y+forward) && hasEnemy(p.X-1, p.Y+forward) {
			newPiece := p
			newPiece.X = p.X - 1
			newPiece.Y = p.Y + forward
			newPiece.HasMoved = true
			moves = append(moves, Move{p, newPiece})
		}
		if isInBounds(p.X+1, p.Y+forward) && hasEnemy(p.X+1, p.Y+forward) {
			newPiece := p
			newPiece.X = p.X + 1
			newPiece.Y = p.Y + forward
			newPiece.HasMoved = true
			moves = append(moves, Move{p, newPiece})
		}

		// check for en passant
		var enPassantActive bool
		var enPassantFile int
		if p.Side == White {
			enPassantFile = b.BlackEnPassent
		} else {
			enPassantFile = b.WhiteEnPassant
		}
		if enPassantFile != -1 && (p.X-1 == enPassantFile || p.X+1 == enPassantFile) {
			if p.Side == White && p.Y == 4 {
				newPiece := p
				newPiece.X = enPassantFile
				newPiece.Y = 5
				newPiece.HasMoved = true
				moves = append(moves, Move{p, newPiece})
			}
			if p.Side == Black && p.Y == 3 {
				newPiece := p
				newPiece.X = enPassantFile
				newPiece.Y = 2
				newPiece.HasMoved = true
				moves = append(moves, Move{p, newPiece})
			}
		}
	case Rook:
		// check all four directions
		moves = raycast(moves, 1, 0)
		moves = raycast(moves, -1, 0)
		moves = raycast(moves, 0, 1)
		moves = raycast(moves, 0, -1)
	case Knight:
		// check all eight L shapes
		dx := []int{1, 1, -1, -1, 2, 2, -2, -2}
		dy := []int{2, -2, 2, -2, 1, -1, 1, -1}
		for i := 0; i < 8; i++ {
			if isInBounds(p.X+dx, p.Y+dy) && (isClear(p.X+dx, p.Y+dy) || hasEnemy(p.X+dx, p.Y+dy)) {
				newPiece := p
				newPiece.X = p.X + dx
				newPiece.Y = p.Y + dy
				newPiece.HasMoved = true
				moves = append(moves, Move{p, newPiece})
			}
		}
	case Bishop:
		// check all four directions
		moves = raycast(moves, 1, -1)
		moves = raycast(moves, -1, 1)
		moves = raycast(moves, 1, 1)
		moves = raycast(moves, -1, -1)
	case King:
		// check his eight moves
		dx := []int{0, 0, -1, -1, -1, 1, 1, 1}
		dy := []int{1, -1, -1, 0, 1, -1, 0, 1}
		for i := 0; i < 8; i++ {
			if isInBounds(p.X+dx, p.Y+dy) && (isClear(p.X+dx, p.Y+dy) || hasEnemy(p.X+dx, p.Y+dy)) {
				newPiece := p
				newPiece.X = p.X + dx
				newPiece.Y = p.Y + dy
				newPiece.HasMoved = true
				moves = append(moves, Move{p, newPiece})
			}
		}
	case Queen:
		// check all eight directions
		moves = raycast(moves, 1, 0)
		moves = raycast(moves, -1, 0)
		moves = raycast(moves, 0, 1)
		moves = raycast(moves, 0, -1)

		moves = raycast(moves, 1, -1)
		moves = raycast(moves, -1, 1)
		moves = raycast(moves, 1, 1)
		moves = raycast(moves, -1, -1)
	}

	// TODO: check moves for moves that would lead to own check or are out of bounds
	// TODO: special case castle; check all squares between start and end position
	return moves
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
	Moves    []Move
	Pieces   []Piece
	Captured []Piece
	State    BoardState
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
	return Board{
		Pieces: []Piece{
			Piece{0, 0, Rook, White, false},
			Piece{1, 0, Knight, White, false},
			Piece{2, 0, Bishop, White, false},
			Piece{3, 0, Queen, White, false},
			Piece{4, 0, King, White, false},
			Piece{5, 0, Bishop, White, false},
			Piece{6, 0, Knight, White, false},
			Piece{7, 0, Rook, White, false},
			Piece{0, 1, Pawn, White, false},
			Piece{1, 1, Pawn, White, false},
			Piece{2, 1, Pawn, White, false},
			Piece{3, 1, Pawn, White, false},
			Piece{4, 1, Pawn, White, false},
			Piece{5, 1, Pawn, White, false},
			Piece{6, 1, Pawn, White, false},
			Piece{7, 1, Pawn, White, false},
			Piece{0, 7, Rook, Black, false},
			Piece{1, 7, Knight, Black, false},
			Piece{2, 7, Bishop, Black, false},
			Piece{3, 7, Queen, Black, false},
			Piece{4, 7, King, Black, false},
			Piece{5, 7, Bishop, Black, false},
			Piece{6, 7, Knight, Black, false},
			Piece{7, 7, Rook, Black, false},
			Piece{0, 6, Pawn, Black, false},
			Piece{1, 6, Pawn, Black, false},
			Piece{2, 6, Pawn, Black, false},
			Piece{3, 6, Pawn, Black, false},
			Piece{4, 6, Pawn, Black, false},
			Piece{5, 6, Pawn, Black, false},
			Piece{6, 6, Pawn, Black, false},
			Piece{7, 6, Pawn, Black, false},
		},
		Captured:          []PieceType{},
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
		Pieces:            make([]Piece, len(b.Pieces)),
		Captured:          make([]Piece, len(b.Captured)),
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
	copy(newBoard.Pieces, b.Pieces)
	copy(newBoard.Captured, b.Captured)
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
	// TODO: check for check; complicated due to pinning
	// TODO: reset draw ask
	// TODO: update moves since capture
	// TODO: update en passant states
	// TODO: check for 50 move draw
	// TODO: check for 3 fold repetition
	panic("unimplemented!")
}

func (b *Board) IsValid() bool {
	// make sure there are the appropriate number of pieces on the board and captured
	if (len(b.Pieces) + len(b.Captured)) != 32 {
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

func (s Side) Opposite() Side {
	if s == White {
		return Black
	} else {
		return White
	}
}

type Move struct {
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
	DisloyaltyForbidden
	WrongSide
	OutOfBounds
	PieceNotFound
	TypeChangeNotAllowed
	InvalidMove
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

	// make sure the piece doesn't change sides
	if m.Start.Side != m.End.Side {
		return nil, DisloyaltyForbidden
	}

	// make sure the move is from the right player
	if (b.State == WhiteMove && m.Start.Side != White) || (b.State == BlackMove && m.Start.Side != Black) {
		return nil, WrongSide
	}

	// make sure the piece exists
	pieceFound := false
	for _, piece := range toCheck {
		if (piece.X == m.Start.X) && (piece.Y == m.End.Y) && (piece.Side == m.Start.Side) {
			pieceFound = true
			break
		}
	}
	if !pieceFound {
		return nil, PieceNotFound
	}

	// make sure it's a possible move
	possibleMoves := m.Start.GetPossibleMoves(b)
	moveFound := false
	for _, move := range moves {
		if move == m {
			moveFound = true
		}
	}
	if !moveFound {
		return nil, InvalidMove
	}

	// else it's good
	return &afterMove, MoveOkay
}

func (m Move) String() string {
	// TODO: write move as standard chess notation
	panic("unimplemented!")
}
