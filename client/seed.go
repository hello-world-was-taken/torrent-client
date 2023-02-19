package client


import (
	"net"
	"log"
	"fmt"

	"torrent-dsp/model"
	"torrent-dsp/client"
	// "time"
    // "encoding/binary"
    // "encoding/hex"
)

func SeederMain() {
	// start a server listening on port 6881
    ln, err := net.Listen("tcp", ":8080")

    if err != nil {
        log.Fatalf("Failed to listen: %s", err)
    }

    for {
        if conn, err := ln.Accept(); err == nil {
            go handleConnection(conn)
        }
    }
}


func handleConnection(conn net.Conn) {
	handShake, err := ReceiveHandShake(conn)
	if err != nil {
		fmt.Println("Error receiving handshake")
		return
	}
}


func ReceiveHandShake(conn net.Conn) (*model.HandShake, error) {
	// read the handshake response
	// handshake response is 68 bytes long <length of Protocol> + <Protocol> + <Reserved> + <InfoHash> + <PeerID>
	// TODO: increase time out
	buffer := make([]byte, 68)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading handshake response")
		return &model.HandShake{}, err
	}

	// deserialize the handshake response
	handShake, err := client.DeserializeHandShake(buffer)
	if err != nil {
		fmt.Println("Error deserializing handshake response")
		return &model.HandShake{}, err
	}

	fmt.Println("Handshake sent successfully")
	return handShake, nil
}


func SendBitField(conn net.Conn, bitField []byte) error {
	// send the bitfield
	_, err := conn.Write(bitField)
	if err != nil {
		return err
	}
	return nil
}


// func maintainUpload() {
//     var timePassed int
//     for {
//         time.Sleep(time.Second)
//         timePassed++
//         uploadedPieces := 0
//         for _, value := range uploadRates {
//             uploadedPieces += value
//         }
//         currentUploadSpeed := float64(uploadedPieces * pieceLen) / float64(timePassed)
//         if allowedUpload > currentUploadSpeed {
//             if globalSleepUpload > 0 {
//                 lock.Lock()
//                 globalSleepUpload--
//                 lock.Unlock()
//             }
//             continue
//         } else {
//             lock.Lock()
//             globalSleepUpload++
//             lock.Unlock()
//         }
//         if hasAllPieces(recievedData) {
//             break
//         }
//     }
// }


// func chokeUnchokeMechanism() {
//     count := 0
//     for {
//         for _, peer := range choked {
//             if downloadRates[peer] > 0 {
//                 sortedUnchoked := make([]string, len(unchoked))
//                 copy(sortedUnchoked, unchoked)
//                 sort.Slice(sortedUnchoked, func(i, j int) bool {
//                     return downloadRates[sortedUnchoked[i]] < downloadRates[sortedUnchoked[j]]
//                 })
//                 // comparing the slowest peer that is unchoked with our choked peer
//                 if downloadRates[peer] > downloadRates[sortedUnchoked[0]] {
//                     mutex.Lock()
//                     unchoked = append(unchoked, peer)
//                     choked = append(choked, sortedUnchoked[0])
//                     unchoked = removePeer(unchoked, sortedUnchoked[0])
//                     mutex.Unlock()
//                 }
//             }
//         }
//         time.Sleep(10 * time.Second)
//         count += 1
//         if count % 3 == 0 {
//             // every 30 seconds, optimistically unchoking one peer
//             if len(choked) > 0 {
//                 mutex.Lock()
//                 unchoked = append(unchoked, choked[0])
//                 choked = append(choked[:0], choked[1:]...)
//                 mutex.Unlock()
//                 if len(unchoked) > 0 {
//                     mutex.Lock()
//                     choked = append(choked, unchoked[0])
//                     unchoked = append(unchoked[:0], unchoked[1:]...)
//                     mutex.Unlock()
//                 }
//             }
//         }
//         if endAllThreads {
//             return
//         }
//     }
// }

// func removePeer(peers []string, peer string) []string {
//     for i, p := range peers {
//         if p == peer {
//             return append(peers[:i], peers[i+1:]...)
//         }
//     }
//     return peers
// }
