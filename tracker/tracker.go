package tracker

import "strconv"

const maxBufferSize = 10000

type Tracker interface {
	GetPeerList(tracker string) ([]string, error)
}

func NewTCPTracker(piece map[string]interface{}) *TCPTracker {
	tracker := TCPTracker{
		infoHash:   piece["info_hash"].([20]byte),
		peerID:     piece["peer_id"].([20]byte),
		port:       piece["port"].(int32),
		downloaded: piece["downloaded"].(int64),
		left:       piece["left"].(int64),
		uploaded:   piece["uploaded"].(int64),
		event:      piece["event"].(string),
		compact:    piece["compact"].(bool),
	}
	return &tracker
}

func NewUDPTracker(piece map[string]interface{}) *UDPTracker {
	tracker := UDPTracker{
		infoHash:   piece["info_hash"].([20]byte),
		peerID:     piece["peer_id"].([20]byte),
		port:       piece["port"].(int32),
		downloaded: piece["downloaded"].(int64),
		left:       piece["left"].(int64),
		uploaded:   piece["uploaded"].(int64),
		event:      piece["event"].(string),
		compact:    piece["compact"].(bool),
	}
	return &tracker
}

func getPort(add string) int32 {
	ret := int32(0)

	for i, pow := len(add)-1, int32(1); i >= 0; i-- {
		if add[i] == ':' {
			break
		}

		ret += int32(add[i]-'0') * pow
		pow *= 10
	}

	return ret
}

// converting Big indian bytes array to bo
func toBit(val []byte) []byte {
	ret := make([]byte, 8*len(val))

	for i, k := len(val)-1, 0; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			if int32(((val[i]) & (1 << j))) != 0 {
				ret[k] = 1
			}
			k++
		}
	}
	return ret
}

func toInt(val []byte) interface{} {
	// MSB is always zero, check the sign
	switch len(val) {
	case 4:
		var ret int32
		val = toBit(val)

		for i := 0; i < 32; i++ {
			ret |= (int32(val[i])) << i
		}

		return ret
	case 8:
		var ret int64
		val = toBit(val)

		for i := 0; i < 64; i++ {
			ret |= int64(val[i]) << i
		}

		return ret
	default: // 2
		var ret int32
		val = toBit(val)

		for i := 0; i < 16; i++ {
			ret |= int32(val[i]) << i
		}
		return ret
	}

}

func toBuf(valInterface interface{}) []byte {
	var ret []byte
	var size int
	var val int64

	switch valInterface.(type) {
	case int64:
		ret = make([]byte, 8)
		size = 64
		val = valInterface.(int64)
	default:
		ret = make([]byte, 4)
		size = 32
		val = int64(valInterface.(int32))
	}

	// checking neg val as MSB is always zero :(
	neg := false
	idx := -1

	if val < 0 {
		neg = true
		val *= -1
	}

	for counter, tempVal, pow := 0, 0, 1; counter < size; counter++ {
		if val%2 == 1 && idx == -1 {
			idx = counter
		}

		tempVal = tempVal + pow*int(val%2)
		pow = pow << 1
		val = val >> 1

		if counter%8 == 7 {
			ret[len(ret)-counter/8-1] = byte(tempVal)
			tempVal = 0
			pow = 1
		}
	}

	if neg {
		for i := 0; i < len(ret); i++ {
			ret[i] ^= 0xff
		}

		for i := 0; i < idx/8; i++ {
			ret[i] ^= 0xff
		}
		ret[len(ret)-idx/8-1] ^= (1 << (idx%8 + 1)) - 1
	}

	return ret
}

func toIP(val []byte) string {
	return strconv.Itoa(int(val[0])) + "." + strconv.Itoa(int(val[1])) + "." + strconv.Itoa(int(val[2])) + "." + strconv.Itoa(int(val[3]))
}

func toPort(val []byte) string {
	return strconv.Itoa(int(val[0]) + int(val[1]))
}
