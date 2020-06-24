package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type URLInfo struct {
	infoHash   [20]byte
	peerID     [20]byte
	port       int32
	uploaded   int64
	downloaded int64
	left       int64
	event      string
	compact    bool
}

func NewURLInfo(piece map[string]interface{}) *URLInfo {
	url := URLInfo{
		infoHash:   piece["info_hash"].([20]byte),
		peerID:     piece["peer_id"].([20]byte),
		port:       piece["port"].(int32),
		uploaded:   piece["uploaded"].(int64),
		downloaded: piece["downloaded"].(int64),
		left:       piece["left"].(int64),
		event:      piece["event"].(string),
		compact:    piece["compact"].(bool),
	}
	return &url
}

func (urlInfo *URLInfo) SendTrackerRequest(tracker string) error {
	fmt.Println("sending req: ", tracker)
	tempCompact := "1"
	if !urlInfo.compact {
		tempCompact = "0"
	}

	req, err := http.NewRequest("GET", tracker, nil)
	if err != nil {
		return err
	}

	urlParameters := url.Values{
		"info_hash":  []string{string(urlInfo.infoHash[:])},
		"peer_id":    []string{string(urlInfo.peerID[:])},
		"port":       []string{strconv.Itoa(int(urlInfo.port))},
		"uploaded":   []string{strconv.FormatInt(urlInfo.uploaded, 10)},
		"downloaded": []string{strconv.FormatInt(urlInfo.downloaded, 10)},
		"left":       []string{strconv.FormatInt(urlInfo.left, 10)},
		"event":      []string{urlInfo.event},
		"compact":    []string{tempCompact},
	}

	req.URL.RawQuery = urlParameters.Encode()
	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("err: ", err)
		return err
	}

	defer res.Body.Close()
	fmt.Println("reqUrl: ", req.URL.String())
	fmt.Println("res: ", res)
	return nil
}
