package main

import "torrent-dsp/client"

func main() {
	client.ConnectToTracker()

	// infinite loop to keep the program running
	for { }
}