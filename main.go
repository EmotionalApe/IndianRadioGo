package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

type Station struct {
	Name string
	URL  string
}

var stations = []Station{
	{Name: "Radio Mirchi", URL: "https://stream-161.zeno.fm/wgaznsmt92quv?zs=I3sEd_C7ToqwlXEfOuH5Dg"},
	{Name: "Mirchi Top20", URL: "https://drive.uber.radio/uber/bollywoodnow/icecast.audio"},
	{Name: "Mirchi Love", URL: "https://stream-142.zeno.fm/pxc55r5uyc9uv?zs=ZHx1l_tVQRKBuRaR5Tmmiw"},
	{Name: "Non Stop Hindi", URL: "https://www.liveradio.es/http://198.178.123.14:8216/stream"},
	{Name: "Bollywood Hits ", URL: "https://stream-142.zeno.fm/rqqps6cbe3quv?zs=qV1P76HuQwCcMNbro4nllQ&rj-ttl=5&rj-tok=AAABccQyYG8AGJ3KW_903"},
	{Name: "Big FM ", URL: "https://funasia.streamguys1.com/bigdallas"},
	{Name: "Radio City ", URL: "https://stream-142.zeno.fm/pxc55r5uyc9uv?zs=ZHx1l_tVQRKBuRaR5Tmmiw"},
	{Name: "Fox FM ", URL: "https://cp11.serverse.com/proxy/foxfm/stream/;stream.mp3"},
	{Name: "Radio City Metal ", URL: "https://stream-144.zeno.fm/apxnpe2zrf9uv?zs=PDAh_CNpQfCnQqGvtGMZEw"},
	{Name: "Virgin Radio", URL: "https://stream-155.zeno.fm/agtp9c146qzuv?zs=lxaXquchRCqswfTwdCqLbg"},
}

func main() {
	fmt.Println("Welcome to Online Radio Player ♪ ♫ ")

	for i, station := range stations {
		fmt.Printf("%d. %s\n", i+1, station.Name)
	}

	fmt.Print("Enter Station Number ")
	var stationNumber int
	_, err := fmt.Scanln(&stationNumber)
	if err != nil || stationNumber < 1 || stationNumber > len(stations) {
		log.Fatal("Invalid station number.")
	}

	selectedStation := stations[stationNumber-1]

	resp, err := http.Get(selectedStation.URL)
	if err != nil {
		log.Fatalf("Failed to connect to the station: %v", err)
	}
	defer resp.Body.Close()

	streamer, format, err := mp3.Decode(resp.Body)
	if err != nil {
		log.Fatalf("Failed to decode the station's stream: %v", err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	fmt.Printf("Playing %s...\n", selectedStation.Name)

	// Channel for playback completion signal
	done := make(chan bool)

	// Play stream with callback for completion signal
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true // Signal completion
	})))

	// Wait for playback completion
	<-done
}
