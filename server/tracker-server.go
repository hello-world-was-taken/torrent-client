// create a tracker server listening at port 6869 using net package
package server

import (
	"log"
	"net"
	"bufio"
	"io"
	"fmt"
)


func StartTracker() {
	ln, err := net.Listen("tcp", ":6869")
	fmt.Println("Listening on port 6869...")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		fmt.Println("Found a connection")
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}


// handleConnection handles the connection from the client and reads the data using bufio package
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// read the data from the connection
	data, err := bufio.NewReader(conn).ReadString('-')
	if err != nil {
		if err == io.EOF {
			fmt.Println("Connection closed by client")
		} else {
			log.Fatal(err)
		}
	}

	// print the data to the console
	fmt.Print(data)
}