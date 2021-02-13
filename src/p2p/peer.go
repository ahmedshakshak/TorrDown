package p2p

import (
	"net"
	"tracker"
)

type Peer struct {
	Conn   net.Conn
	ID     []byte
	Pieces []int32
}

type PeerResponse int32

const (
	keepAlive PeerResponse = iota
	Choke
	Unchoke
	Interested
	NotInterested
	Have
	Bitfield
	Request
	Piece
	Cancel
)

func (p *Peer) SendKeepAliveMessage(conn *net.Conn) ([]byte, error) {
	keepAliveMessage := string(tracker.ToBuf(int32(0)))
	return tracker.SendReq(conn, &keepAliveMessage)
}

func (p *Peer) SendChokeMessage(conn *net.Conn) ([]byte, error) {
	message := string(tracker.ToBuf(int32(1))) + string(0)
	return tracker.SendReq(conn, &message)
}

func (p *Peer) SendUnchokeMessage(conn *net.Conn) ([]byte, error) {
	message := string(tracker.ToBuf(int32(1))) + string(1)
	return tracker.SendReq(conn, &message)
}

func (p *Peer) SendInterestedMessage(conn *net.Conn) ([]byte, error) {
	message := string(tracker.ToBuf(int32(1))) + string(2)
	return tracker.SendReq(conn, &message)
}

func (p *Peer) SendNotInterestedMessage(conn *net.Conn) ([]byte, error) {
	message := string(tracker.ToBuf(int32(1))) + string(3)
	return tracker.SendReq(conn, &message)
}

// pieceIdx: zero based
func (p *Peer) SendHaveMessage(conn *net.Conn, pieceIdx int32) ([]byte, error) {
	message := string(tracker.ToBuf(int32(5))) + string(4) + string(tracker.ToBuf(int32(pieceIdx)))
	return tracker.SendReq(conn, &message)
}

// The payload is a bitfield representing the pieces that have been successfully downloaded.
// The high bit in the first byte corresponds to piece index 0. Bits that are cleared indicated a missing piece,
// and set bits indicate a valid and available piece.
// Spare bits at the end are set to zero.
func (p *Peer) SendBitfieldMessage(conn *net.Conn, pieceIdx int32, begin int32, bitField []byte) ([]byte, error) {
	message := string(tracker.ToBuf(int32(9+len(bitField)))) + string(5) + string(tracker.ToBuf(int32(pieceIdx))) +
		string(tracker.ToBuf(int32(begin))) + string(bitField)
	return tracker.SendReq(conn, &message)
}

// used to send a block request
func (p *Peer) SendRequestMessage(conn *net.Conn, pieceIdx int32, begin int32, length int32) ([]byte, error) {
	message := string(tracker.ToBuf(int32(13))) + string(6) + string(tracker.ToBuf(int32(pieceIdx))) +
		string(tracker.ToBuf(int32(begin))) + string(tracker.ToBuf(int32(length)))
	return tracker.SendReq(conn, &message)
}

// used to send a requested block
func (p *Peer) SendPiecetMessage(conn *net.Conn, pieceIdx int32, begin int32, block []byte) ([]byte, error) {
	message := string(tracker.ToBuf(int32(9+len(block)))) + string(7) + string(tracker.ToBuf(int32(pieceIdx))) +
		string(tracker.ToBuf(int32(begin))) + string(block)
	return tracker.SendReq(conn, &message)
}

// used to cancel a block request
func (p *Peer) SendCanceltMessage(conn *net.Conn, pieceIdx int32, begin int32, length int32) ([]byte, error) {
	message := string(tracker.ToBuf(int32(13))) + string(8) + string(tracker.ToBuf(int32(pieceIdx))) +
		string(tracker.ToBuf(int32(begin))) + string(tracker.ToBuf(int32(length)))
	return tracker.SendReq(conn, &message)
}
