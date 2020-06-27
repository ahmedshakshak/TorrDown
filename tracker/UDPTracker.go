package tracker

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type UDPTracker struct {
	infoHash      [20]byte
	peerID        [20]byte
	port          int32
	uploaded      int64
	downloaded    int64
	left          int64
	event         string
	compact       bool
	connectingID  int64
	transactionID int32
	interval      int64
	seeders       int64
	leachers      int64
}

// tracker: tracker's address
// return list of peers
func (t *UDPTracker) GetPeerList(address string) ([]string, error) {
	return t.getPeerList(address)
}

func (t *UDPTracker) getPeerList(address string) ([]string, error) {
	// connecting to the tracker
	tempAddress := address[6 : len(address)-9] // removing `udp//` and `/announce`
	conn, err := t.connectTracker(&tempAddress)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	//setting reeading timeout
	deadline := time.Now().Add(5 * time.Second)
	err = conn.SetReadDeadline(deadline)
	if err != nil {
		return nil, err
	}

	// sending connecting packet
	t.transactionID = 123
	buffer, err := t.sendConnectingPacket(&conn)
	if err != nil {
		return nil, err
	}

	// sending announcinig packet
	buffer, err = t.sendAnnouncingPacket(&conn, buffer)
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for i := 20; i < len(buffer); i += 6 {
		peerIP := buffer[i : i+4]
		peerPort := buffer[i+4 : i+6]
		ret = append(ret, toIP(peerIP)+":"+toPort(peerPort))
	}

	t.interval = int64(toInt(buffer[8:12]).(int32))
	t.leachers = int64(toInt(buffer[12:16]).(int32))
	t.seeders = int64(toInt(buffer[16:20]).(int32))

	// sending scrape packet
	buffer, err = t.sendScapingPacket(&conn, buffer)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *UDPTracker) sendScapingPacket(conn *net.Conn, annPacketRes []byte) ([]byte, error) {
	action := int32(2)
	message := string(toBuf(t.connectingID)) + string(toBuf(action)) + string(toBuf(t.transactionID))
	buffer, err := t.sendUDPReq(conn, &message)

	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func (t *UDPTracker) sendAnnouncingPacket(conn *net.Conn, connPacketRes []byte) ([]byte, error) {
	action := int32(1)
	t.connectingID = toInt(connPacketRes[8:16]).(int64)
	message := string(toBuf(t.connectingID)) + string(toBuf(action)) + string(toBuf(t.transactionID)) + string(t.infoHash[:]) +
		string(t.peerID[:]) + string(toBuf(t.downloaded)) + string(toBuf(t.left)) + string(toBuf(t.uploaded)) +
		string(toBuf(int32(2))) + string(toBuf(int32(0))) + string(toBuf(int32(1234))) + string(toBuf(int32(-1))) + string(getPort((*conn).LocalAddr().String()))

	buffer, err := t.sendUDPReq(conn, &message)
	if err != nil {
		return nil, err
	}

	if len(buffer) < 20 || action != toInt(buffer[0:4]).(int32) || t.transactionID != toInt(buffer[4:8]).(int32) {
		return nil, errors.New("announcing packet error: " + string(buffer[8:]))
	}

	return buffer, nil
}

func (t *UDPTracker) sendConnectingPacket(conn *net.Conn) ([]byte, error) {
	connectingID := int64(0x41727101980)
	action := int32(0)
	message := string(toBuf(connectingID)) + string(toBuf(action)) + string(toBuf(t.transactionID))

	buffer, err := t.sendUDPReq(conn, &message)
	if err != nil {
		return nil, err
	}

	if len(buffer) < 16 || action != toInt(buffer[0:4]).(int32) || t.transactionID != toInt(buffer[4:8]).(int32) {
		return nil, errors.New("Couldn't send connecting packet to the tracker")
	}

	return buffer, nil
}

func (t *UDPTracker) connectTracker(tracker *string) (net.Conn, error) {
	UDPaddr, err := net.ResolveUDPAddr("udp", *tracker) // resolving host names

	if err != nil {
		return nil, err
	}

	// connecting to the tracker
	conn, err := net.Dial("udp", UDPaddr.String())
	if err != nil {
		return nil, err
	}

	return conn, err
}

func (t *UDPTracker) sendUDPReq(conn *net.Conn, message *string) ([]byte, error) {
	//sending message to the tracker
	fmt.Fprintf(*conn, *message)
	// reading tracker response
	buffer := make([]byte, maxBufferSize)
	n, err := (*conn).Read(buffer)
	return buffer[:n], err
}
