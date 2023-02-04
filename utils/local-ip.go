package utils

import (
	"net"
	"log"
)

// GetLocalIP returns the local IP address of the machine
func GetLocalIP() string {
	con, error := net.Dial("udp", "8.8.8.8:80")
	if error != nil {
		log.Fatal(error)
	}

	defer con.Close()
	
	localAddr := con.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}