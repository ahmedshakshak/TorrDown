package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/jackpal/bencode-go"
)

// torrent file information
type TorrentFile struct {
	announceList []string
	comment      string
	createdBy    string
	length       int64
	creationDate int64
	encoding     string
	info         map[string]interface{}
}

func (t *TorrentFile) Read(filePath string) error {
	torretnFile, err := ioutil.ReadFile(filePath)

	if err != nil {
		return err
	}

	torrentFileReader := strings.NewReader(string(torretnFile))
	data, err := bencode.Decode(torrentFileReader)

	if err != nil {
		return err
	}

	dataMap := data.(map[string]interface{})
	//	fmt.Println(dataMap)
	switch dataMap["announce"].(type) {
	case string:
		t.announceList = []string{dataMap["announce"].(string)}
	default:
		t.announceList = []string{}
	}

	switch dataMap["announce-list"].(type) {
	case []interface{}:

		for i := 0; i < len(dataMap["announce-list"].([]interface{})); i++ {
			t.announceList = append(t.announceList, dataMap["announce-list"].([]interface{})[i].([]interface{})[0].(string))
		}
	default:
		fmt.Println(reflect.TypeOf(dataMap["announce-list"]))
		// do nothing
	}

	if len(t.announceList) < 0 {
		return fmt.Errorf("There is no tracker in torrent file")
	}

	switch dataMap["comment"].(type) {
	case string:
		t.comment = dataMap["comment"].(string)
	default:
		t.comment = ""
	}

	switch dataMap["created_by"].(type) {
	case string:
		t.createdBy = dataMap["created_by"].(string)
	default:
		t.createdBy = ""
	}

	switch dataMap["creation_date"].(type) {
	case int64:
		t.creationDate = dataMap["creation_date"].(int64)
	default:
		t.creationDate = 0
	}

	switch dataMap["encoding"].(type) {
	case string:
		t.encoding = dataMap["encoding"].(string)
	default:
		t.encoding = ""
	}

	switch dataMap["length"].(type) {
	case int64:
		t.length = dataMap["length"].(int64)
	default:
		t.length = 0
	}

	switch dataMap["info"].(type) {
	case map[string]interface{}:
		t.info = dataMap["info"].(map[string]interface{})
	default:
		return fmt.Errorf("There is no file information")
	}

	return nil
}

func (t *TorrentFile) getInfoHash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, (*t).info)
	infoHash := sha1.Sum(buf.Bytes())
	return infoHash, err
}
