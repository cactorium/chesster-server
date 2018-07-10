package chesster

type PieceType int

const (
	Pawn PieceType = iota
	Rook
	Knight
	Bishop
	King
	Queen
)

type Piece struct {
	X    int
	Y    int
	Type PieceType
}

type Board struct {
	White []Piece
	Black []Piece
}
