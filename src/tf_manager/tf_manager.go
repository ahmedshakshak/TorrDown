package tf_manager

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jackpal/bencode-go"
)

// torrent file information
type TorrentFile struct {
	AnnounceList           []string
	Comment                string
	CreatedBy              string
	CreationDate           int64
	Encoding               string
	Info                   map[string]interface{}
	DownloadLength         int64
	InfoHash               [20]byte
	FilesLengths           []int64
	FilesPaths             []string
	PieceIsBeingDownloaded []bool
	PieceFileIndex         []int32
	PieceLength            int64
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
	switch dataMap["announce"].(type) {
	case string:
		t.AnnounceList = []string{dataMap["announce"].(string)}
	default:
		t.AnnounceList = []string{}
	}

	switch dataMap["announce-list"].(type) {
	case []interface{}:
		for i := 0; i < len(dataMap["announce-list"].([]interface{})); i++ {
			t.AnnounceList = append(t.AnnounceList, dataMap["announce-list"].([]interface{})[i].([]interface{})[0].(string))
		}
	default:
		// do nothing
	}

	if len(t.AnnounceList) < 0 {
		return fmt.Errorf("There is no tracker in torrent file")
	}

	switch dataMap["comment"].(type) {
	case string:
		t.Comment = dataMap["comment"].(string)
	default:
		t.Comment = ""
	}

	switch dataMap["created by"].(type) {
	case string:
		t.CreatedBy = dataMap["created by"].(string)
	default:
		t.CreatedBy = ""
	}

	switch dataMap["creation date"].(type) {
	case int64:
		t.CreationDate = dataMap["creation date"].(int64)
	default:
		t.CreationDate = 0
	}

	switch dataMap["encoding"].(type) {
	case string:
		t.Encoding = dataMap["encoding"].(string)
	default:
		t.Encoding = ""
	}

	switch dataMap["info"].(type) {
	case map[string]interface{}:
		t.Info = dataMap["info"].(map[string]interface{})
		t.PieceLength = t.Info["piece length"].(int64)
		t.FilesPaths = []string{}
		t.FilesLengths = []int64{}
		t.PieceIsBeingDownloaded = []bool{}
		directoryName := t.Info["name"].(string)
		if t.Info["files"] != nil {
			for i := 0; i < len(t.Info["files"].([]interface{})); i++ {
				t.DownloadLength += t.Info["files"].([]interface{})[i].(map[string]interface{})["length"].(interface{}).(int64)
				t.FilesLengths = append(t.FilesLengths, t.Info["files"].([]interface{})[i].(map[string]interface{})["length"].(interface{}).(int64))
				path := directoryName

				for _, file := range t.Info["files"].([]interface{})[i].(map[string]interface{})["path"].([]interface{}) {
					path += "/" + file.(string)
				}
				t.FilesPaths = append(t.FilesPaths, path)
				t.FilesPaths[len(t.FilesPaths)-1] = directoryName + "/" + t.FilesPaths[len(t.FilesPaths)-1]

				for i := int64(0); i < t.FilesLengths[len(t.FilesLengths)-1]; i += int64(t.PieceLength) {
					t.PieceIsBeingDownloaded = append(t.PieceIsBeingDownloaded, false)
					t.PieceFileIndex = append(t.PieceFileIndex, int32(len(t.FilesLengths)-1))
				}
			}
		} else {
			t.DownloadLength += t.Info["length"].(interface{}).(int64)
			t.FilesLengths = append(t.FilesLengths, t.Info["length"].(interface{}).(int64))
			t.FilesPaths = append(t.FilesPaths, directoryName)
			for i := int64(0); i < t.FilesLengths[len(t.FilesLengths)-1]; i += int64(t.PieceLength) {
				t.PieceIsBeingDownloaded = append(t.PieceIsBeingDownloaded, false)
				t.PieceFileIndex = append(t.PieceFileIndex, int32(len(t.FilesLengths)-1))
			}
		}

	default:
		return fmt.Errorf("There is no file information")
	}

	var buf bytes.Buffer
	bencode.Marshal(&buf, t.Info)
	t.InfoHash = sha1.Sum(buf.Bytes())
	return nil
}
