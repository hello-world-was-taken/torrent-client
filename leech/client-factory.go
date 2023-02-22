package leech

import (
	"bytes"
	"net"
	"time"

	"torrent-dsp/constant"
	"torrent-dsp/model"
	"torrent-dsp/utils"
)

func ClientFactory(peer model.Peer, torrent model.Torrent) (*model.Client, error) {
	client, err := createClient(peer, torrent)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// to create a client we need to:
// 1. connect to the peer
// 2. shake hands with the peer
// 3. receive bit field message from the peer
func createClient(peer model.Peer, torrent model.Torrent) (*model.Client, error) {
	conn, err := connectToPeer(peer, torrent)
	if err != nil {
		return nil, err
	}

	// shake hands with the peer
	err = ShakeHandWithPeer(torrent, peer, constant.CLIENT_ID, conn)
	if err != nil {
		return nil, err
	}

	// receive bit field message from the peer
	bitFieldMessage, err := ReceiveBitFieldMessage(conn)
	if err != nil {
		return &model.Client{}, err
	}

	// create a new client
	client := &model.Client{
		Peer:        peer,
		BitField:    bitFieldMessage.Payload,
		Conn:        conn,
		ChokedState: constant.CHOKE,
	}

	return client, nil
}

// connect to peer. If not possible then return error
func connectToPeer(peer model.Peer, torrent model.Torrent) (net.Conn, error) {

	conn, err := net.DialTimeout("tcp", peer.String(), constant.CONNECTION_TIMEOUT)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func ShakeHandWithPeer(torrent model.Torrent, peer model.Peer, clientID string, conn net.Conn) error {

	conn.SetDeadline(time.Now().Add(constant.CONNECTION_TIMEOUT))
	defer conn.SetDeadline(time.Time{})
	
	// convert the client id to a byte array
	clientIDByte := [20]byte{}
	copy(clientIDByte[:], []byte(clientID))

	// create a handshake request
	handshakeRequest := model.HandShake{
		Pstr:     "BitTorrent protocol",
		InfoHash: torrent.InfoHash,
		PeerID:   clientIDByte,
	}

	// send the handshake request
	handshakeResponse, err := handshakeRequest.Send(conn)
	if err != nil {
		return err
	}

	// check that the infohash in the response matches the infohash of the torrent
	if !bytes.Equal(handshakeResponse.InfoHash[:], torrent.InfoHash[:]) {
		return err
	}

	// check that the peer id in the response is different from ours
	if bytes.Equal(handshakeResponse.PeerID[:], utils.ConvertStringToByteArray(constant.CLIENT_ID)[:]) {
		return err
	}

	return nil
}

func ReceiveBitFieldMessage(conn net.Conn) (*model.Message, error) {
	conn.SetDeadline(time.Now().Add(constant.CONNECTION_TIMEOUT))
	defer conn.SetDeadline(time.Time{}) // Disable the deadline


	// receive the bitField message
	// fmt.Println("Deserializing bit field message...")
	bitFieldMessageResponse, err := model.DeserializeMessage(conn)
	if err != nil {
		// fmt.Println("Error receiving bit field message")
		return nil, err
	}

	// check that the message is a bit field message
	if bitFieldMessageResponse.MessageID != constant.BIT_FIELD {
		// fmt.Println("Expected bit field message")
		return &model.Message{}, nil
		// log.Fatal("expected bit field message")
	}

	return bitFieldMessageResponse, nil
}
