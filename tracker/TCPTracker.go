package tracker

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jackpal/bencode-go"
)

type TCPTracker struct {
	infoHash   [20]byte
	peerID     [20]byte
	port       int32
	uploaded   int64
	downloaded int64
	left       int64
	event      string
	compact    bool
	interval   int64
	seeders    int64
	leachers   int64
}

// tracker: tracker address
// return list of peers

func (t *TCPTracker) GetPeerList(tracker string) ([]string, error) {
	req, err := http.NewRequest("GET", tracker, nil)
	if err != nil {
		return nil, err
	}

	urlParameters := t.trackerParameters()
	req.URL.RawQuery = urlParameters.Encode()
	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	buf := make([]byte, maxBufferSize)
	n, err := res.Body.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}

	resReader := strings.NewReader(string(buf[:n]))
	data, err := bencode.Decode(resReader)

	if err != nil {
		return nil, err
	}

	dataMap := data.(map[string]interface{})
	peers := []string{}

	switch dataMap["peers"].(type) {
	case []string:
		peers = append(peers, dataMap["peers"].([]string)...)

	default: // string
		tempPeers := dataMap["peers"].(string)

		for i := 0; i < len(tempPeers); i += 6 {
			peers = append(peers, toIP([]byte(tempPeers[i:i+4]))+":"+toPort([]byte(tempPeers[i+4:i+6])))
		}
	}

	t.interval = dataMap["min interval"].(int64)
	t.seeders = dataMap["complete"].(int64)
	t.leachers = dataMap["incomplete"].(int64)
	return peers, nil
}

func (t *TCPTracker) trackerParameters() url.Values {
	tempCompact := "1"
	if !t.compact {
		tempCompact = "0"
	}

	urlParameters := url.Values{
		"info_hash":  []string{string(t.infoHash[:])},
		"peer_id":    []string{string(t.peerID[:])},
		"port":       []string{strconv.Itoa(int(t.port))},
		"uploaded":   []string{strconv.FormatInt(t.uploaded, 10)},
		"downloaded": []string{strconv.FormatInt(t.downloaded, 10)},
		"left":       []string{strconv.FormatInt(t.left, 10)},
		"event":      []string{t.event},
		"compact":    []string{tempCompact},
	}

	return urlParameters
}
