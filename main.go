package main

import (
	"fmt"
	"log"
)

const version = "1"
const userID = "0"
const fileName = `/home/ahmed/Downloads/test.torrent` //`/home/ahmed/Downloads/Udemy - Certified Kubernetes Administrator (CKA) with Practice Tests.torrent`
const peerID = "12345678912345678901"                 //"TorrDown:v" + version + ":" + userID
const port = int32(3000)

func main() {
	var torrentFile TorrentFile
	err := torrentFile.Read(fileName)

	if err != nil {
		log.Fatal(err)
	}

	//	fmt.Print(torrentFile.info["files"].([]interface{})[0].(map[string]interface{})["path"].([]interface{})[1].(string))
	info_hash, err := torrentFile.getInfoHash()
	fmt.Println(info_hash)

	tempPeer := [20]byte{}

	for i := 0; i < 20; i++ {
		tempPeer[i] = peerID[i]
	}

	piece := map[string]interface{}{}
	piece["info_hash"] = [20]byte(info_hash)
	piece["peer_id"] = tempPeer
	piece["port"] = port
	piece["downloaded"] = int64(0)
	piece["uploaded"] = int64(0)
	piece["left"] = int64(torrentFile.length)
	piece["event"] = "started"
	piece["compact"] = true

	url := NewURLInfo(piece)

	for i := 0; i < len(torrentFile.announceList); i++ {
		url.SendTrackerRequest(torrentFile.announceList[i])
	}

	if err != nil {
		fmt.Println("err: ", err)
	}

}
