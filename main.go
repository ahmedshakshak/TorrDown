package main

import (
	"fmt"
	"log"
)

const version = "1"
const userID = "0"
const fileName = `C:\Users\Ahmed\Downloads\6201484321_f1a88ca2cb_b_archive.torrent`
const peerID = "12345678912345678901" //"TorrDown:v" + version + ":" + userID
const port = int32(3000)

func main() {
	var torrentFile TorrentFile
	err := torrentFile.Read(fileName)

	if err != nil {
		log.Fatal(err)
	}

	//	fmt.Print(torrentFile.info["files"].([]interface{})[0].(map[string]interface{})["path"].([]interface{})[1].(string))
	info_hash, err := torrentFile.getInfoHash()
	fmt.Println(torrentFile.numOfFiles)
	tempPeer := [20]byte{}

	for i := 0; i < 20; i++ {
		tempPeer[i] = peerID[i]
	}

	piece := map[string]interface{}{}
	piece["info_hash"] = info_hash
	piece["peer_id"] = tempPeer
	piece["port"] = port
	piece["downloaded"] = int64(0)
	piece["uploaded"] = int64(0)
	piece["left"] = torrentFile.downloadLength
	piece["event"] = "started"
	piece["compact"] = true

	tcp := tracker.NewTCPTracker(piece)
	udp := tracker.NewUDPTracker(piece)

	for i := 0; i < len(torrentFile.announceList); i++ {
		fmt.Println(torrentFile.announceList[i])
		if torrentFile.announceList[i][:3] == "udp" {
			peerlist, err := udp.GetPeerList(torrentFile.announceList[i])
			fmt.Println(err)
			fmt.Println(peerlist)
		} else {
			peerlist, err := tcp.GetPeerList(torrentFile.announceList[i])
			fmt.Println(err)
			fmt.Println(peerlist)
		}
	}

	if err != nil {
		fmt.Println("err: ", err)
	}
}
