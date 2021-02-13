package p2p

import (
	"fmt"
	"net"
	"reflect"
	"tracker"
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

func (p *P2P) SendHandShake(errChan chan *Peer) (*Peer, error) {
	conn, err := net.Dial("tcp", p.address)

	if err != nil {
		errChan <- nil
		return nil, err
	}

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

	if len(res) < len(message) || !reflect.DeepEqual(res[:len(p.pstr)+1], []byte(message[:len(p.pstr)+1])) ||
		!reflect.DeepEqual(res[len(p.pstr)+9:len(message)-len(p.peerID)],
			[]byte(message[len(p.pstr)+9:len(message)-len(p.peerID)])) {

		errChan <- nil
		return nil, fmt.Errorf("Peers are Not talking about the same torrent file")
	}

	peer := &Peer{
		Conn:   conn,
		ID:     res[len(message)-len(p.peerID) : len(message)],
		Pieces: []int32{},
	}

	//	fmt.Printf("peer res:\nsize: %v\nmessage:%v \n", len(res)-len(message), res[len(message):])
	for i, j := len(message), int32(0); i < len(res); i++ {
		for k := 0; k < 8; k++ {
			if res[i]&(1<<k) != 0 {
				peer.Pieces = append(peer.Pieces, j)
			}

			j++
		}
	}

	errChan <- peer
	return peer, nil
}
