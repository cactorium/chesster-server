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
	isOutOfBounds := (x < 0) || (x > 8) || (y < 0) || (y > 8)
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
			ret = append(ret, Move{Start: p, End: newPiece})
		}
		if isInBounds(x, y) && hasEnemy(x, y) {
			newPiece := p
			newPiece.X = x
			newPiece.Y = y
			newPiece.HasMoved = true
			ret = append(ret, Move{Start: p, End: newPiece})
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
				for _, pieceType := range []PieceType{Rook, Knight, Bishop, Queen} {
					newNewPiece := newPiece
					newNewPiece.Type = pieceType
					moves = append(moves, Move{Start: p, End: newNewPiece, IsPromotion: true})
				}
			} else {
				moves = append(moves, Move{Start: p, End: newPiece})
			}
		}
		// add forward two squares if it's not blocked
		if !p.HasMoved {
			if isInBounds(p.X, p.Y+2*forward) && isClear(p.X, p.Y+forward) && isClear(p.X, p.Y+2*forward) {
				newPiece := p
				newPiece.Y = p.Y + 2*forward
				newPiece.HasMoved = true
				moves = append(moves, Move{Start: p, End: newPiece})
			}
		}
		// check possible captures
		if isInBounds(p.X-1, p.Y+forward) && hasEnemy(p.X-1, p.Y+forward) {
			newPiece := p
			newPiece.X = p.X - 1
			newPiece.Y = p.Y + forward
			newPiece.HasMoved = true
			moves = append(moves, Move{Start: p, End: newPiece})
		}
		if isInBounds(p.X+1, p.Y+forward) && hasEnemy(p.X+1, p.Y+forward) {
			newPiece := p
			newPiece.X = p.X + 1
			newPiece.Y = p.Y + forward
			newPiece.HasMoved = true
			moves = append(moves, Move{Start: p, End: newPiece})
		}

		// check for en passant
		var enPassantFile int
		if p.Side == White {
			enPassantFile = b.BlackEnPassant
		} else {
			enPassantFile = b.WhiteEnPassant
		}
		if enPassantFile != -1 && (p.X-1 == enPassantFile || p.X+1 == enPassantFile) {
			if p.Side == White && p.Y == 4 {
				newPiece := p
				newPiece.X = enPassantFile
				newPiece.Y = 5
				newPiece.HasMoved = true
				moves = append(moves, Move{Start: p, End: newPiece})
			}
			if p.Side == Black && p.Y == 3 {
				newPiece := p
				newPiece.X = enPassantFile
				newPiece.Y = 2
				newPiece.HasMoved = true
				moves = append(moves, Move{Start: p, End: newPiece})
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
			if isInBounds(p.X+dx[i], p.Y+dy[i]) && (isClear(p.X+dx[i], p.Y+dy[i]) || hasEnemy(p.X+dx[i], p.Y+dy[i])) {
				newPiece := p
				newPiece.X = p.X + dx[i]
				newPiece.Y = p.Y + dy[i]
				newPiece.HasMoved = true
				moves = append(moves, Move{Start: p, End: newPiece})
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
			if isInBounds(p.X+dx[i], p.Y+dy[i]) && (isClear(p.X+dx[i], p.Y+dy[i]) || hasEnemy(p.X+dx[i], p.Y+dy[i])) {
				newPiece := p
				newPiece.X = p.X + dx[i]
				newPiece.Y = p.Y + dy[i]
				newPiece.HasMoved = true
				moves = append(moves, Move{Start: p, End: newPiece})
			}
		}
		// add castling
		if !p.HasMoved {
			for _, piece := range b.Pieces {
				if piece.Type == Rook && piece.Side == p.Side && !piece.HasMoved {
					if piece.X == 0 && isClear(1, p.Y) && isClear(2, p.Y) && isClear(3, p.Y) {
						moves = append(moves, Move{IsCastle: true})
					}
					if piece.X == 7 && isClear(6, p.Y) && isClear(5, p.Y) {
						moves = append(moves, Move{IsCastle: true, IsKingsideCastle: true})
					}
				}
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

	isValid := make([]bool, len(moves))
	// check for moves that would lead to own check
	// TODO: special case castle; check all squares between start and end position

	for i, m := range moves {
		if !isInBounds(m.End.X, m.End.Y) {
			isValid[i] = false
			continue
		}
		if !m.IsCastle {
			// copy the pieces into a new board, do the move and see if it causes
			// the same colored king to be in check

			nb := b.Clone()
			nb.CommitMove(m)
			isValid[i] = !nb.InCheck(p.Side)
		} else {
			// TODO handle castling
		}
	}

	validMoves := []Move{}
	for i, m := range moves {
		if isValid[i] {
			validMoves = append(validMoves, m)
		}
	}
	return validMoves
}

type BoardState int

const (
	// white's turn to move
	WhiteMove BoardState = iota
	// black's turn to move
	BlackMove
)

// TODO: probably refactor board into game state with board embedded
type Board struct {
	Pieces   []Piece
	Captured []Piece
	State    BoardState
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
		Captured:       []Piece{},
		State:          WhiteMove,
		WhiteEnPassant: -1,
		BlackEnPassant: -1,
	}
}

func EmptyBoard() Board {
	return Board{
		Pieces:         []Piece{},
		Captured:       []Piece{},
		State:          WhiteMove,
		WhiteEnPassant: -1,
		BlackEnPassant: -1,
	}
}

func (b *Board) Clone() Board {
	newBoard := Board{
		Pieces:         make([]Piece, len(b.Pieces)),
		Captured:       make([]Piece, len(b.Captured)),
		State:          b.State,
		WhiteEnPassant: b.WhiteEnPassant,
		BlackEnPassant: b.BlackEnPassant,
	}
	copy(newBoard.Pieces, b.Pieces)
	copy(newBoard.Captured, b.Captured)
	return newBoard
}

func (b *Board) IsMove(s Side) bool {
	return (s == White && b.State == WhiteMove) || (s == Black && b.State == BlackMove)
}

func (b *Board) getPiece(x, y int) *Piece {
	for i, p := range b.Pieces {
		if p.X == x && p.Y == y {
			return &b.Pieces[i]
		}
	}
	return nil
}

func (b *Board) getKing(s Side) *Piece {
	for i, p := range b.Pieces {
		if p.Type == King || p.Side == s {
			return &b.Pieces[i]
		}
	}
	return nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
func minInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
func absInt(a int) int {
	return maxInt(a, -a)
}

func (b *Board) CommitMove(m Move) bool {
	if m.IsCastle {
		// TODO special case castle
		panic("unimplemented!")
		return false
	}
	moving := b.getPiece(m.Start.X, m.Start.Y)
	if moving == nil || moving.Side != m.Start.Side || moving.Type != m.Start.Type {
		return false
	}
	captured := b.getPiece(m.End.X, m.End.Y)
	// special case en passant
	if m.Start.Type == Pawn && (absInt(m.End.Y-m.Start.Y) == 2 && absInt(m.End.X-m.Start.X) == 1) {
		var d int
		if m.End.Y > m.Start.Y {
			d = 1
		} else {
			d = -1
		}
		captured = b.getPiece(m.End.X, m.Start.Y+d)
		if captured == nil || captured.Type != Pawn {
			return false
		}
	}
	if captured != nil && captured.Side != m.Start.Side.Opposite() {
		return false
	}

	*moving = m.End
	moving.HasMoved = true
	if captured != nil {
		// replace it with the last piece and shift it all down one to remove it
		*captured = b.Pieces[len(b.Pieces)-1]
		b.Pieces = b.Pieces[:len(b.Pieces)-1]
	}
	return true
}

func (b *Board) InCheck(s Side) bool {
	king := b.getKing(s)
	if king == nil {
		// NOTE: probably should print out some error info
		return false
	}
	kx := king.X
	ky := king.Y

	isBlocking := func(x, y int) bool {
		return b.getPiece(x, y) != nil
	}
	for _, p := range b.Pieces {
		// no team kills; probably the least realistic part of chess
		if p.Side == s {
			continue
		}
		switch p.Type {
		case Pawn:
			var forward int
			if p.Side == White {
				forward = 1
			} else {
				forward = -1
			}
			if absInt(p.X-kx) == 1 && p.Y == ky-forward {
				return true
			}
		case Rook:
			if p.X == kx {
				start := minInt(p.Y, ky)
				end := maxInt(p.Y, ky)

				canThreaten := true
				x := kx
				for y := start; y < end; y++ {
					if isBlocking(x, y) {
						canThreaten = false
						break
					}
				}
				if canThreaten {
					return true
				}
			}
			if p.Y == ky {
				start := minInt(p.X, kx)
				end := maxInt(p.X, kx)

				canThreaten := true
				y := ky
				for x := start; x < end; x++ {
					if isBlocking(x, y) {
						canThreaten = false
						break
					}
				}
				if canThreaten {
					return true
				}
			}
		case Knight:
			if absInt(p.X-kx) == 1 && absInt(p.Y-ky) == 2 {
				return true
			}
			if absInt(p.X-kx) == 2 && absInt(p.Y-ky) == 1 {
				return true
			}
		case Bishop:
			// on the +x +y diagonal, the difference between x and y stays the same
			// on the +x -y diagonal, their sum stays the same
			if kx-ky == p.X-p.Y {
				panic("unimplemented")
				// TODO
			}
			if kx+ky == p.X+p.Y {
				panic("unimplemented")
				// TODO
			}
		case Queen:
			panic("unimplemented")
			// TODO
		case King:
			if absInt(p.X-kx) <= 1 || absInt(p.Y-ky) <= 1 {
				return true
			}
		}
	}
	return false
}

func (b *Board) IsValid() bool {
	// make sure there are the appropriate number of pieces on the board and captured
	if (len(b.Pieces) + len(b.Captured)) != 32 {
		return false
	}
	// make sure all the pieces are in bounds
	for _, piece := range b.Pieces {
		if piece.X < 0 || piece.X > 8 || piece.Y < 0 || piece.Y > 8 {
			return false
		}
	}

	// make sure that each side has a king
	if b.getKing(Black) == nil || b.getKing(White) == nil {
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
	AfraidOfCommitment
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
	if maybePiece := b.getPiece(m.Start.X, m.Start.Y); maybePiece != nil {
		if maybePiece.Side != m.Start.Side || maybePiece.Type != m.Start.Type {
			return nil, PieceNotFound
		}
	} else {
		return nil, PieceNotFound
	}

	// make sure it's a possible move
	moves := m.Start.GetPossibleMoves(b)
	moveFound := false
	for _, move := range moves {
		if move == m {
			moveFound = true
		}
	}
	if !moveFound {
		return nil, InvalidMove
	}

	afterMove := b.Clone()
	if !afterMove.CommitMove(m) {
		return nil, AfraidOfCommitment
	}

	// else it's good
	return &afterMove, MoveOkay
}

func (m Move) String() string {
	// TODO: write move as standard chess notation
	panic("unimplemented!")
}
