package main

import (
	"sync"
	"time"
	"torrent-dsp/leech"
	"torrent-dsp/seed"
)

func main() {
	var wg sync.WaitGroup
  	wg.Add(1)
  
  	go func() {
    	defer wg.Done()
    	seed.SeederMain()
  	}()


	time.Sleep(5 * time.Second)
	leech.Leech("./torrent-files/debian-11.6.0-amd64-netinst.iso.torrent")

	wg.Wait()
}
