package main

import (
	// "fmt"

	"torrent-dsp/client"
	// "torrent-dsp/utils"
	// "torrent-dsp/server"
)

func main() {
	// go server.StartTracker()
	client.ConnectToTracker()
	// fmt.Println(utils.GetLocalIP())

	// infinite loop to keep the program running
	for { }
}