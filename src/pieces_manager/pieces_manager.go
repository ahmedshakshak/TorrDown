package pieces_manager

import "os"

const perm = 0744

type PiecesManager struct {
	path      string
	fileName  string
	pieces    []Piece
	fileSize  int64
	fileExist bool
}

func NewPiecesManager(absPath string, fileSize int64) PiecesManager {
	fileName := ""
	pathLastIndex := len(absPath)

	for i := len(absPath) - 1; i >= 0; i-- {
		if absPath[i] == '/' {
			pathLastIndex = i
			break
		}

		fileName = string(absPath[i]) + fileName
	}

	path := absPath[:pathLastIndex]
	p := PiecesManager{
		path:      path,
		fileName:  fileName,
		fileSize:  fileSize,
		fileExist: false,
	}

	return p
}

func (p *PiecesManager) createPath() error {
	return os.MkdirAll(p.path, perm)
}

func (p *PiecesManager) createFile() error {
	/*if _, err := os.State(p.path + "/" + p.fileName); os.IsNotExist(err) {
		f, err := ioutil.Create(p.path + "/" + p.fileName)
		if err != nil {
			return err
		}

		defer f.Close()
		_, err = f.Write(make([]byte, p.fileSize))
		return err
	}*/
	return nil
}

func (p *PiecesManager) WritePiece(piece *Piece) error {
	if !p.fileExist {
		err := p.createPath()
		if err != nil {
			return err
		}

		err = p.createFile()
		if err != nil {
			return err
		}

		p.fileExist = true
	}

	f, err := os.OpenFile(p.path+"/"+p.fileName, os.O_RDWR, perm)
	if err != nil {
		return err
	}

	for _, block := range piece.blocks {
		f.WriteAt(block.data, piece.start+block.start)
	}

	f.Close()
	return nil
}
