package pieces_manager

type Block struct {
	data  []byte
	start int64
}

type Piece struct {
	blocks []Block
	start  int64
	size   int32
}

func NewPiece(start int64, size int32) Piece {
	p := Piece{
		start: start,
		size:  size,
	}

	return p
}

func NewBlocke(start int64, data []byte) Block {
	b := Block{
		start: start,
		data:  make([]byte, len(data)),
	}

	copy(b.data, data)
	return b
}

func (p *Piece) AddBLock(block *Block) {
	p.blocks = append(p.blocks, *block)
}
