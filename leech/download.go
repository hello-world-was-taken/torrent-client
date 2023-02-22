package leech

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"torrent-dsp/common"
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

// parse the torrent metadata and get a list of peers from the tracker
func PrepareDownload(filename string) (model.Torrent, []model.Peer) {
	// open torrent file from the current directory and parse it
	torrent, err := common.ParseTorrentFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// get a list of peers from the tracker
	peers, err := GetPeersFromTrackers(&torrent)
	if err != nil {
		log.Fatal(err)
	}

	return torrent, peers
}

func StartDownload(filename string) {

	torrent, peers := PrepareDownload(filename)

	// create output file
	outFile, filename, err := common.CreateFile(&torrent)
	if err != nil {
		log.Fatalf("Error creating output file: ", err)
	}
	defer outFile.Close()

	// load the cache from a file
	piecesCache, err := LoadCache(outFile.Name() + ".json")
	if err != nil {
		// error loading file, assign pieces hash to an new map[int]bool
		piecesCache = &model.PiecesCache{}
		fmt.Println("Error loading file")
		return
	}

	// create two channels for the download and upload
	piecesHashList := torrent.Info.PiecesToByteArray()
	downloadChannel := make(chan *PieceRequest, len(piecesHashList))
	resultChannel := make(chan *PieceResult)

	for idx, hash := range piecesHashList {
		length := torrent.CalcRequestSize(idx)

		// TODO: there might be an off by one error here
		downloadChannel <- &PieceRequest{Index: idx, Hash: hash, Length: length}
	}

	// start the download and upload goroutines
	for _, peer := range peers {
		go DownloadFromPeer(peer, torrent, downloadChannel, resultChannel, piecesCache)
	}

	// TODO: this needs to be changed
	// Collect results into a buffer until full
	buf := make([]byte, torrent.Info.Length)
	donePieces := 0

	StoreDownloadedPieces(donePieces, torrent, resultChannel, err, outFile, piecesCache, buf)

	fmt.Println("Done downloading all pieces")
	close(downloadChannel)

}

func StoreDownloadedPieces(donePieces int, torrent model.Torrent, resultChannel chan *PieceResult, err error, outFile *os.File, piecesCache *model.PiecesCache, buf []byte) {

	for len(piecesCache.Pieces) < len(torrent.Info.PiecesToByteArray()) {
		res := <-resultChannel

		// calculate the start and end index of the piece
		pieceSize := int(torrent.Info.PieceLength)
		pieceStartIdx := res.Index * pieceSize
		pieceEndIdx := utils.CalcMin(pieceStartIdx+pieceSize, int(torrent.Info.Length))

		// prepare output file
		_, err = outFile.WriteAt(res.Block, int64(pieceStartIdx))
		if err != nil {
			log.Fatalf("Failed to write to file: %s", "downloaded_file.iso")
		}

		// update the cache
		piecesCache.Pieces[res.Index] = true
		SaveCache(outFile.Name()+".json", piecesCache)

		// TODO: do we need to store it on the buffer?
		copy(buf[pieceStartIdx:pieceEndIdx], res.Block)
		donePieces++

		// print the progress
		percent := float64(len(piecesCache.Pieces)) / float64(len(torrent.Info.PiecesToByteArray())) * 100
		numWorkers := runtime.NumGoroutine() - 1
		log.Printf("Downloading... (%0.2f%%) Active Peers: %d\n", percent, numWorkers)
	}
	return
}

func DownloadFromPeer(peer model.Peer, torrent model.Torrent, downloadChannel chan *PieceRequest, resultChannel chan *PieceResult, piecesCache *model.PiecesCache) {
	// create a client with the peer
	client, err := ClientFactory(peer, torrent)
	if err != nil {
		fmt.Printf("Failed to create a client with peer %s", peer.String())
		return
	}

	// prepare for download
	client.UnChoke()
	client.Interested()

	// download the pieces the peer has
	for piece := range downloadChannel {
		fmt.Println("Found from cache: ", !piecesCache.Pieces[piece.Index])
		_, ok := piecesCache.Pieces[piece.Index]
		
		if !ok && utils.BitOn(client.BitField, piece.Index) {

			// send request message to the peer
			_, err = DownloadPiece(piece, client, downloadChannel, resultChannel, &torrent)
			if err != nil {
				downloadChannel <- piece
				return
			}
		} else if !ok {
			downloadChannel <- piece
		}
	}
}

func DownloadPiece(piece *PieceRequest, client *model.Client, downloadChannel chan *PieceRequest, resultChannel chan *PieceResult, torrent *model.Torrent) (PieceResult, error) {

	// set the deadline for the connection
	client.Conn.SetDeadline(time.Now().Add(constant.PIECE_DOWNLOAD_TIMEOUT))
	defer client.Conn.SetDeadline(time.Time{})

	totalDownloaded := 0
	requested := 0
	blockDownloadCount := 0
	blockLength := constant.MAX_BLOCK_LENGTH
	fmt.Println("Downloading piece: ", piece.Index, piece.Length, torrent.Info.PieceLength)
	buffer := make([]byte, piece.Length)

	for totalDownloaded < piece.Length {

		if client.ChokedState != constant.CHOKE {
			for blockDownloadCount < constant.MAX_BATCH_DOWNLOAD {
				length := blockLength
				// Last block might be shorter than the typical block
				if piece.Length-requested < blockLength {
					length = piece.Length - requested
				}

				// send request message to the peer
				err := client.Request(uint32(piece.Index), uint32(requested), uint32(length))
				if err != nil {
					downloadChannel <- piece
					return PieceResult{}, err
				}
				requested += length
				blockDownloadCount++
			}
		}

		// collect the response
		message, err := model.DeserializeMessage(client.Conn)
		if err != nil {
			downloadChannel <- piece
			return PieceResult{}, err
		}

		// keep alive
		if message == nil {
			downloadChannel <- piece
			return PieceResult{}, err
		}

		switch message.MessageID {
		case constant.CHOKE:
			client.ChokedState = constant.CHOKE
		case constant.UN_CHOKE:
			client.ChokedState = constant.UN_CHOKE
		case constant.INTERESTED:
			ParseInterested(message)
		case constant.NOT_INTERESTED:
			ParseNotInterested(message)
		case constant.HAVE:
			index, err := ParseHave(message)
			if err != nil {
				fmt.Println("Error parsing have message from peer: ", client.Peer.String())
				return PieceResult{}, err
			}
			utils.TurnBitOn(client.BitField, index)
		case constant.REQUEST:
			ParseRequest(message)
		case constant.PIECE:
			n, err := ParsePiece(piece.Index, buffer, message)
			if err != nil {
				fmt.Println("Error parsing piece message from peer: ", client.Peer.String())
				downloadChannel <- piece
				return PieceResult{}, err
			}
			totalDownloaded += n
			blockDownloadCount--
		case constant.CANCEL:
			ParseCancel(message)
		}

	}

	// verify the piece
	if !utils.BitHashChecker(buffer, piece.Hash) {
		return PieceResult{}, fmt.Errorf("Piece hash verification failed for piece: %d", piece.Index)
	}

	// send the piece to the result channel
	resultChannel <- &PieceResult{Index: piece.Index, Block: buffer}

	return PieceResult{}, nil
}

func ParseInterested(msg *model.Message) {
	// Check that the message is a INTERESTED message.
	if msg.MessageID != constant.INTERESTED {
		fmt.Errorf("Expected INTERESTED (ID %d), got ID %d", constant.INTERESTED, msg.MessageID)
	}
}

func ParseNotInterested(msg *model.Message) {
	// Check that the message is a NOT_INTERESTED message.
	if msg.MessageID != constant.NOT_INTERESTED {
		fmt.Errorf("Expected NOT_INTERESTED (ID %d), got ID %d", constant.NOT_INTERESTED, msg.MessageID)
	}
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

func ParseRequest(msg *model.Message) {
	// TODO: spawn goroutine to handle request
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

func ParseCancel(msg *model.Message) {
	// Check that the message is a CANCEL message.
	if msg.MessageID != constant.CANCEL {
		fmt.Errorf("Expected CANCEL (ID %d), got ID %d", constant.CANCEL, msg.MessageID)
	}
}
