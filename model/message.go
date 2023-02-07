package model

import "encoding/binary"

type Message struct {
	Length int    `bencode:"length"`
	MessageID int    `bencode:"message_id"`
	Payload []byte `bencode:"payload"`
}

const (
	CHOKE = 0
	UNCHOKE = 1
	INTERESTED = 2
	NOT_INTERESTED = 3
	HAVE = 4
	BITFIELD = 5
	REQUEST = 6
	PIECE = 7
	CANCEL = 8
)

// serialize the message into a byte array
// <msgLength 4 byte> <message id 1 byte> <payload>
func (message *Message) Serialize() []byte {
	msgLength := uint32(len(message.Payload) + 1) // +1 for the message id
	buffer := make([]byte, msgLength + 4) // +4 for the msgLength prefix
	binary.BigEndian.PutUint32(buffer[0:4], msgLength)
	buffer[4] = byte(message.MessageID)
	copy(buffer[5:], message.Payload)

	return buffer
}


// deserialize the message byte array into a Message struct
func DeserializeMessage(buffer []byte) *Message {
	message := &Message{}
	message.Length = int(binary.BigEndian.Uint32(buffer[0:4]))
	message.MessageID = int(buffer[4])
	message.Payload = buffer[5:]

	return message
}