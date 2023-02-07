package model

type HandShake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

// serialize the handshake request
// <pstrlen 1 byte><pstr 19 bytes i.e Bit Torrent><reserved 8 byte><info_hash><peer_id 20 bytes>
func (handShake *HandShake) Serialize() []byte {
	// create a buffer to store the handshake request. For bittorrent, request length is 49 + len(pstr) ==> 49 + 19 = 68
	buffer := make([]byte, 49 + len(handShake.Pstr))
	buffer[0] = byte(len(handShake.Pstr))
	copy(buffer[1:], []byte(handShake.Pstr))
	copy(buffer[1+len(handShake.Pstr):], make([]byte, 8))
	copy(buffer[1+len(handShake.Pstr)+8:], handShake.InfoHash[:])
	copy(buffer[1+len(handShake.Pstr)+8+20:], handShake.PeerID[:])

	return buffer
}


// deserialize the handshake response byte array into a HandShake struct
func DeserializeHandShake(buffer []byte) *HandShake {
	handShake := &HandShake{}
	pstrLength := int(buffer[0])
	handShake.Pstr = string(buffer[1:pstrLength+1])
	copy(handShake.InfoHash[:], buffer[pstrLength+28:pstrLength+48])
	copy(handShake.PeerID[:], buffer[pstrLength+48:pstrLength+68])

	return handShake
}