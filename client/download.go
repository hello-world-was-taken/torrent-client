package client

// import (
// 	"torrent-dsp/model"
// )

// type PieceResult struct {
// 	Index int `bencode:"index"`
// 	Begin int `bencode:"begin"`
// 	Block []byte `bencode:"block"`
// }

// type PieceRequest struct {
// 	Index int `bencode:"index"`
// 	Hash [20]byte `bencode:"hash"`
// 	Length int `bencode:"length"`
// }

// func StartDownload(torrent model.Torrent) {
// 	// create two channels for the download and upload
// 	piecesHashList := torrent.Info.PiecesToByteArray()
// 	downloadChannel := make(chan *PieceRequest, len(piecesHashList))
// 	resultChannel := make(chan *PieceResult)

// 	for idx, hash := range piecesHashList {
// 		downloadChannel <- &PieceRequest{Index: idx, Hash: hash, Length: 0}
// 	}

// }


// func DownloadPiece(piece PieceRequest, downloadChannel chan *PieceRequest, resultChannel chan *PieceResult) {
// 	// download the piece
	
// }


// func UploadPiece(piece PieceResult) {
// 	// upload the piece
// }


