package client

import (
	// "crypto/sha1"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"torrent-dsp/constant"
	"torrent-dsp/model"
	"torrent-dsp/utils"
)


type PieceResult struct {
	Index int    `bencode:"index"`
	Begin int    `bencode:"begin"`
	Block []byte `bencode:"block"`
}


type PieceRequest struct {
	Index  int      `bencode:"index"`
	Hash   [20]byte `bencode:"hash"`
	Length int      `bencode:"length"`
}


func StartDownload(torrent model.Torrent, peers []model.Peer) {
	// create two channels for the download and upload
	piecesHashList := torrent.Info.PiecesToByteArray()
	downloadChannel := make(chan *PieceRequest, len(piecesHashList))
	resultChannel := make(chan *PieceResult)
	// fmt.Println("piecesHashList length ", len(piecesHashList))
	for idx, hash := range piecesHashList {
		pieceSize := int(torrent.Info.PieceLength)
		pieceStartIdx := idx * pieceSize
		pieceEndIdx := utils.CalcMin(pieceStartIdx + pieceSize, int(torrent.Info.Length))
		
		// TODO: there might be an off by one error here		
		downloadChannel <- &PieceRequest{Index: idx, Hash: hash, Length: pieceEndIdx - pieceStartIdx}
	}

	// start the download and upload goroutines
	for _, peer := range peers {
		go DownloadFromPeer(peer, torrent, downloadChannel, resultChannel)
	}


	// TODO: this needs to be changed
	// Collect results into a buffer until full
	buf := make([]byte, torrent.Info.Length)
	donePieces := 0
	for donePieces < len(torrent.Info.PiecesToByteArray()) {
		res := <- resultChannel
		pieceSize := int(torrent.Info.PieceLength)
		pieceStartIdx := res.Index * pieceSize
		pieceEndIdx := utils.CalcMin(pieceStartIdx + pieceSize, int(torrent.Info.Length))

		copy(buf[pieceStartIdx:pieceEndIdx], res.Block)
		donePieces++

		percent := float64(donePieces) / float64(len(torrent.Info.PiecesToByteArray())) * 100
		numWorkers := runtime.NumGoroutine() - 1
		log.Printf("Downloading... (%0.2f%%) Active Peers: %d\n", percent, numWorkers)
	}

	fmt.Println("Done downloading all pieces")
	close(downloadChannel)

	path := "downloaded_file.iso"
	outFile, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create file: %s", path)
		// return err
	}
	defer outFile.Close()

	_, err = outFile.Write(buf)
	if err != nil {
		log.Fatalf("Failed to write to file: %s", path)
		// return err
	}
	// return nil

	// return buf, nil

}


func DownloadFromPeer(peer model.Peer, torrent model.Torrent, downloadChannel chan *PieceRequest, resultChannel chan *PieceResult) {
	// create a client with the peer
	client, err := ClientFactory(peer, torrent)
	if err != nil {
		fmt.Printf("Failed to create a client with peer %s", peer.String())
		return
	}
	// fmt.Println("Printing bit field... ")
	// fmt.Println("bit field length ", client.BitField)


	// send un_choke message to the peer and then send interested message
	// fmt.Println("Sending unchoke message to peer: ", client.Peer.String())
	client.UnChoke()
	// time.Sleep(2 * time.Second)
	// fmt.Println("Sending interested message to peer: ", client.Peer.String())
	client.Interested()

	// iterate over the download channel and download the pieces by checking the bitfield
	for piece := range downloadChannel {
		// fmt.Println("Downloading piece: ", piece.Index)
		// check if the piece is available in the bit field
		if utils.BitOn(client.BitField, piece.Index) {
			// send request message to the peer
			// fmt.Println("Sending request message to peer: ", client.Peer.String())
			_, err = DownloadPiece(piece, client, downloadChannel, resultChannel, &torrent)
			if err != nil {
				return
			}
		} else {
			downloadChannel <- piece
		}
	}


	

	// start the download and upload goroutines
	// go DownloadPiece(peer, downloadChannel, resultChannel)
	// go UploadPiece(peer, resultChannel)
}


func DownloadPiece(piece *PieceRequest, client *model.Client, downloadChannel chan *PieceRequest, resultChannel chan *PieceResult, torrent *model.Torrent) (PieceResult, error) {
	
	client.Conn.SetDeadline(time.Now().Add(constant.PIECE_DOWNLOAD_TIMEOUT))
    defer client.Conn.SetDeadline(time.Time{})

	totalDownloaded := 0
	blockDownloadCount := 0
	blockLength := constant.MAX_BLOCK_LENGTH
	buffer := make([]byte, piece.Length)


	for totalDownloaded < piece.Length {

		// fmt.Println("Client choked state: ", client.ChokedState)
		if client.ChokedState != constant.CHOKE {
			for blockDownloadCount < constant.MAX_BATCH_DOWNLOAD {
				length := utils.CalcMin(blockLength, piece.Length - ( blockDownloadCount * blockLength ))

				// send request message to the peer
				err := client.Request(uint32(piece.Index), uint32(blockDownloadCount * blockLength), uint32(length))
				if err != nil {
					// fmt.Println("Error sending request message to peer: ", client.Peer.String())
					downloadChannel <- piece
					return PieceResult{}, err
				}

				blockDownloadCount++
			}
		}

		// collect the response
		// fmt.Println("Waiting for response from peer: ", client.Peer.String())
		message, err := model.DeserializeMessage(client.Conn)
		if err != nil {
			// fmt.Println("Error deserializing message from peer: ", err)
			downloadChannel <- piece
			return PieceResult{}, err
		}
	
		// keep alive
		if message == nil {
			downloadChannel <- piece
			return PieceResult{}, err
		}

		switch message.MessageID {
		case constant.UN_CHOKE:
			client.ChokedState = constant.UN_CHOKE
		case constant.CHOKE:
			client.ChokedState = constant.CHOKE
		case constant.HAVE:
			index, err := ParseHave(message)
			if err != nil {
				fmt.Println("Error parsing have message from peer: ", client.Peer.String())
				return PieceResult{}, err
			}
			utils.TurnBitOn(client.BitField, index)
			// client.BitField[index] = 1
		case constant.PIECE:
			n, err := ParsePiece(piece.Index, buffer, message)
			if err != nil {
				fmt.Println("Error parsing piece message from peer: ", client.Peer.String())
				downloadChannel <- piece
				return PieceResult{}, err
			}
			totalDownloaded += n
			blockDownloadCount--
		}

	}

	// verify the piece
	// if !utils.BitHashChecker(buffer, piece.Hash) {
	// 	fmt.Println("----> sha1",sha1.Sum(buffer))
	// 	fmt.Println("----> piece hash",piece.Hash)
	// 	fmt.Println("Piece hash verification failed for piece: ", piece.Index)
	// 	return PieceResult{}, fmt.Errorf("Piece hash verification failed for piece: %d", piece.Index)
	// }

	// send the piece to the result channel
	resultChannel <- &PieceResult{Index: piece.Index, Block: buffer}

	return PieceResult{}, nil
}


func ParsePiece(index int, buf []byte, msg *model.Message) (int, error) {

	// Check that the message is a PIECE message.
	if msg.MessageID != constant.PIECE {
		return 0, fmt.Errorf("Expected PIECE (ID %d), got ID %d", constant.PIECE, msg.MessageID)
	}

	// Check that the payload is long enough.
	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("Payload too short. %d < 8", len(msg.Payload))
	}

	// Extract the begin offset from the payload.
	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if begin >= len(buf) {
		fmt.Println("begin problem")
		return 0, fmt.Errorf("Begin offset too high. %d >= %d", begin, len(buf))
	}

	// Copy the data from the payload to the buffer.
	data := msg.Payload[8:]
	if begin+len(data) > len(buf) {
		fmt.Println("data problem: ", begin+len(data), " - ", len(buf))
		return 0, fmt.Errorf("Data too long [%d] for offset %d with length %d", len(data), begin, len(buf))
	}
	// fmt.Println("Successfully parsed piece")
	copy(buf[begin:], data)

	// Return the length of the data and no error.
	return len(data), nil
}


func ParseHave(msg *model.Message) (int, error) {
	if msg.MessageID != constant.HAVE {
		return 0, fmt.Errorf("Expected HAVE (ID %d), got ID %d", constant.HAVE, msg.MessageID)
	}

	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("Expected payload length 4, got length %d", len(msg.Payload))
	}

	index := int(binary.BigEndian.Uint32(msg.Payload))
	
	return index, nil
}
