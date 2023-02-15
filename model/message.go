package model

import (
	"encoding/binary"
	"log"
	"net"
)

type Message struct {
	Length uint32    `bencode:"length"`
	MessageID uint8    `bencode:"message_id"`
	Payload []byte `bencode:"payload"`
}

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
func DeserializeMessage(conn net.Conn) *Message {
	length := make([]byte, 4)
	_, err := conn.Read(length)
	if err != nil {
		log.Fatalf("Error reading message length: %s", err)
	}

	msgLength := binary.BigEndian.Uint32(length)
	buffer := make([]byte, msgLength)
	_, err = conn.Read(buffer)
	if err != nil {
		log.Fatalf("Error reading message: %s", err)
	}


	message := &Message{}
	message.Length = msgLength
	message.MessageID = buffer[4]
	message.Payload = buffer[5:]

	return message
}