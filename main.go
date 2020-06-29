package main

import (
	"fmt"
	"log"

	"TorrDown/tracker"
)

const version = "1"
const userID = "0"
const fileName = `C:\Users\Ahmed\Downloads\Justice League Dark Apokolips War (2020) [720p] [BluRay] [YTS.MX].torrent`
const peerID = "12345678901234567890" //"TorrDown:v" + version + ":" + userID
const port = int32(3000)

func main() {
	var torrentFile TorrentFile
	err := torrentFile.Read(fileName)

	if err != nil {
		log.Fatal(err)
	}

	//	fmt.Print(torrentFile.info["files"].([]interface{})[0].(map[string]interface{})["path"].([]interface{})[1].(string))
	info_hash, err := torrentFile.getInfoHash()
	fmt.Println("Number of files:", torrentFile.numOfFiles)
	fmt.Printf("Number of pieces:%v => %v\n", len(torrentFile.info["pieces"].(string)), len(torrentFile.info["pieces"].(string))/20)
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
	peerChan := make(chan *Peer)
	counter := 0

	for i := 0; i < len(torrentFile.announceList); i++ {
		fmt.Println(torrentFile.announceList[i])
		peerlist := []string{}
		if torrentFile.announceList[i][:3] == "udp" {
			peerlist, err = udp.GetPeerList(torrentFile.announceList[i])
		} else {
			peerlist, err = tcp.GetPeerList(torrentFile.announceList[i])
		}

		for _, address := range peerlist {
			p2p, err := NewP2P(&address, info_hash[:], []byte(peerID[:]))
			if err != nil {
				fmt.Printf("Error(P2P): error in creating new P2P(%v)\n", err)
				continue
			}

			go p2p.sendHandShake(peerChan)
			counter++
		}
	}

	fmt.Println("counter size:", counter)
	for counter > 0 {
		peer := <-peerChan
		counter--

		if peer != nil {
			fmt.Println("handshake ret:", peer.pieces)
		}
	}

	if err != nil {
		fmt.Println("err: ", err)
	}
}
