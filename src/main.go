package main

import (
	"fmt"
	"log"
	"p2p"
	"tf_manager"
	"tracker"
)

const version = "1"
const userID = "0"
const fileName = `../homealone.torrent`
const peerID = "12345678901234567890" //"TorrDown:v" + version + ":" + userID
const port = int32(3000)

func main() {
	var torrentFile tf_manager.TorrentFile
	err := torrentFile.Read(fileName)

	if err != nil {
		log.Fatal(err)
	}

	infoHash := torrentFile.InfoHash
	fmt.Println("Number of files:", len(torrentFile.FilesLengths))
	fmt.Printf("Number of pieces:%v => %v\n", len(torrentFile.Info["pieces"].(string)), len(torrentFile.Info["pieces"].(string))/20)
	tempPeer := [20]byte{}

	for i := 0; i < 20; i++ {
		tempPeer[i] = peerID[i]
	}

	piece := map[string]interface{}{}
	piece["info_hash"] = infoHash
	piece["peer_id"] = tempPeer
	piece["port"] = port
	piece["downloaded"] = int64(0)
	piece["uploaded"] = int64(0)
	piece["left"] = torrentFile.DownloadLength
	piece["event"] = "started"
	piece["compact"] = true

	tcp := tracker.NewTCPTracker(piece)
	udp := tracker.NewUDPTracker(piece)
	peerChan := make(chan *p2p.Peer)
	counter := 0

	for i := 0; i < len(torrentFile.AnnounceList); i++ {
		peerlist := []string{}

		if torrentFile.AnnounceList[i][:3] == "udp" {
			peerlist, err = udp.GetPeerList(torrentFile.AnnounceList[i])
		} else {
			peerlist, err = tcp.GetPeerList(torrentFile.AnnounceList[i])
		}

		for _, address := range peerlist {
			p2p, err := p2p.NewP2P(&address, infoHash[:], []byte(peerID[:]))

			if err != nil {
				fmt.Printf("Error(P2P): error in creating new P2P(%v)\n", err)
				continue
			}

			go p2p.SendHandShake(peerChan)
			counter++
		}
	}

	fmt.Println("counter size:", counter)
	for counter > 0 {
		peer := <-peerChan
		counter--

		if peer != nil && len(peer.Pieces) > 0 {
			go func(peer *p2p.Peer) {
				ret, err := peer.SendUnchokeMessage(&peer.Conn)
				if err != nil {
					return
				}

				ret, err = peer.SendInterestedMessage(&peer.Conn)
				if err != nil || len(ret) != 5 || ret[4] != 1 {
					return
				}

				for pieceIdx, downloaded := range torrentFile.PieceIsBeingDownloaded {
					if !downloaded && contains(peer.Pieces, int32(pieceIdx)) {
						torrentFile.PieceIsBeingDownloaded[pieceIdx] = true
						requestedSize := 512
						pieceLength := torrentFile.PieceLength
						if (pieceIdx+1 < len(torrentFile.PieceIsBeingDownloaded) &&
							torrentFile.PieceFileIndex[pieceIdx+1] != torrentFile.PieceFileIndex[pieceIdx]) ||
							pieceIdx == len(torrentFile.PieceIsBeingDownloaded)-1 {
							pieceLength = torrentFile.FilesLengths[torrentFile.PieceFileIndex[pieceIdx]] % torrentFile.PieceLength
						}

						buff := make([]byte, pieceLength)
						done := true
						for blockOffset := int64(0); blockOffset < pieceLength; blockOffset += int64(requestedSize) {
							blockSize := requestedSize
							if pieceLength-blockOffset < int64(requestedSize) {
								blockSize = int(pieceLength - blockOffset)
							}

							ret, err := peer.SendRequestMessage(&peer.Conn, int32(pieceIdx), int32(blockOffset), int32(blockSize))
							if err != nil {
								fmt.Println("request error:", err)
								torrentFile.PieceIsBeingDownloaded[pieceIdx] = false
								done = false
								break
							} else {
								for i, j := blockOffset, 13; j < len(ret); i, j = i+1, j+1 {
									buff[i] = ret[j]
								}
							}
						}
						fmt.Println("Piece with idx:", pieceIdx, done)
						// TODO: writhe the pieces to the hard drive and bingooooooo
					}
				}
			}(peer)
		}
	}

	/*	fmt.Println(err)
		err = ioutil.WriteFile(`C:\Users\Ahmed\testGO\test.txt`, []byte("hi\n do u here mer ? ;D xD"), 0666)
		fmt.Println(err)
		os.Create()*/
}

func contains(l []int32, val int32) bool {
	for _, ele := range l {
		if ele == val {
			return true
		}
	}

	return false
}
