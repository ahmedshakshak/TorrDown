package main

import (
	"TorrDown/tracker"
	"fmt"
	"net"
	"reflect"
	"time"
)

const timeOut = 2

type P2P struct {
	address  string
	pstrlen  byte
	pstr     string
	reserved []byte
	infoHash []byte
	peerID   []byte
}

type Peer struct {
	conn   net.Conn
	ID     []byte
	pieces []int32
}

// peerID Length should be 20 byte.
// infohash length should be 20 byte.
func NewP2P(address *string, infoHash []byte, peerID []byte) (*P2P, error) {

	if len(infoHash) != 20 {
		return nil, fmt.Errorf("infoHash length should be 20 but got %v", len(infoHash))
	}

	if len(peerID) != 20 {
		return nil, fmt.Errorf("peerID length should be 20 but got %v", len(peerID))
	}

	p2p := P2P{
		peerID:   make([]byte, len(peerID)),
		infoHash: make([]byte, len(infoHash)),
		reserved: make([]byte, 8),
		pstr:     "BitTorrent protocol",
		pstrlen:  byte(len("BitTorrent protocol")),
		address:  *address,
	}

	copy(p2p.infoHash, infoHash)
	copy(p2p.peerID, peerID)

	return &p2p, nil
}

func (p *P2P) sendHandShake(errChan chan *Peer) (*Peer, error) {
	conn, err := net.Dial("tcp", p.address)

	if err != nil {
		errChan <- nil
		return nil, err
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(timeOut * time.Second))

	message := string(p.pstrlen) + p.pstr + string(p.reserved) +
		string(p.infoHash) + string(p.peerID)

	//	fmt.Println("peer message: ", message)
	res, err := tracker.SendReq(&conn, &message)

	if err != nil {
		errChan <- nil
		return nil, err
	}

	//fmt.Println("peer res1: ", res[:len(message)-len(p.peerID)])
	//fmt.Println("peer res2: ", []byte(message[:len(message)-len(p.peerID)]))

	if !reflect.DeepEqual(res[:len(p.pstr)+1], []byte(message[:len(p.pstr)+1])) ||
		!reflect.DeepEqual(res[len(p.pstr)+9:len(message)-len(p.peerID)],
			[]byte(message[len(p.pstr)+9:len(message)-len(p.peerID)])) {
		errChan <- nil
		return nil, fmt.Errorf("Peers are Not talking about the same torrent file")
	}

	peer := &Peer{
		conn:   conn,
		ID:     res[len(message)-len(p.peerID) : len(message)],
		pieces: []int32{},
	}

	fmt.Printf("peer res:\nsize: %v\nmessage:%v \n", len(res)-len(message), res[len(message):])
	for i, j := len(message), int32(0); i < len(res); i++ {
		for k := 0; k < 8; k++ {
			if res[i]&(1<<k) != 0 {
				peer.pieces = append(peer.pieces, j)
			}

			j++
		}
	}

	errChan <- peer
	return peer, nil
}

func (p *P2P) sendKeepAliveMessage(conn *net.Conn) ([]byte, error) {
	keepAliveMessage := string(tracker.ToBuf(int32(0)))
	return tracker.SendReq(conn, &keepAliveMessage)
}

func (p *P2P) sendChokeMessage(conn *net.Conn) ([]byte, error) {
	message := string(tracker.ToBuf(int32(1))) + string(0)
	return tracker.SendReq(conn, &message)
}

func (p *P2P) sendUnchokeMessage(conn *net.Conn) ([]byte, error) {
	message := string(tracker.ToBuf(int32(1))) + string(1)
	return tracker.SendReq(conn, &message)
}

func (p *P2P) sendInterestedMessage(conn *net.Conn) ([]byte, error) {
	message := string(tracker.ToBuf(int32(1))) + string(2)
	return tracker.SendReq(conn, &message)
}

func (p *P2P) sendNotInterestedMessage(conn *net.Conn) ([]byte, error) {
	message := string(tracker.ToBuf(int32(1))) + string(3)
	return tracker.SendReq(conn, &message)
}

// pieceIdx: zero based
func (p *P2P) sendHaveMessage(conn *net.Conn, pieceIdx int32) ([]byte, error) {
	message := string(tracker.ToBuf(int32(5))) + string(4) + string(tracker.ToBuf(int32(pieceIdx)))
	return tracker.SendReq(conn, &message)
}

// The payload is a bitfield representing the pieces that have been successfully downloaded.
// The high bit in the first byte corresponds to piece index 0. Bits that are cleared indicated a missing piece,
// and set bits indicate a valid and available piece.
// Spare bits at the end are set to zero.
func (p *P2P) sendBitfieldMessage(conn *net.Conn, pieceIdx int32, begin int32, bitField []byte) ([]byte, error) {
	message := string(tracker.ToBuf(int32(9+len(bitField)))) + string(5) + string(tracker.ToBuf(int32(pieceIdx))) +
		string(tracker.ToBuf(int32(begin))) + string(bitField)
	return tracker.SendReq(conn, &message)
}

// used to send a block request
func (p *P2P) sendRequestMessage(conn *net.Conn, pieceIdx int32, begin int32, length int32) ([]byte, error) {
	message := string(tracker.ToBuf(int32(13))) + string(6) + string(tracker.ToBuf(int32(pieceIdx))) +
		string(tracker.ToBuf(int32(begin))) + string(tracker.ToBuf(int32(length)))
	return tracker.SendReq(conn, &message)
}

// used to send a requested block
func (p *P2P) sendPiecetMessage(conn *net.Conn, pieceIdx int32, begin int32, block []byte) ([]byte, error) {
	message := string(tracker.ToBuf(int32(9+len(block)))) + string(7) + string(tracker.ToBuf(int32(pieceIdx))) +
		string(tracker.ToBuf(int32(begin))) + string(block)
	return tracker.SendReq(conn, &message)
}

// used to cancel a block request
func (p *P2P) sendCanceltMessage(conn *net.Conn, pieceIdx int32, begin int32, length int32) ([]byte, error) {
	message := string(tracker.ToBuf(int32(13))) + string(8) + string(tracker.ToBuf(int32(pieceIdx))) +
		string(tracker.ToBuf(int32(begin))) + string(tracker.ToBuf(int32(length)))
	return tracker.SendReq(conn, &message)
}
