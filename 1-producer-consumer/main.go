//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func producer(stream Stream, datach chan Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			return
		}

		datach <- *tweet
	}
}

func consumer(datach chan Tweet, wg sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)
	for {
		select {

		case t := <-datach:
			{
				if t.IsTalkingAboutGo() {
					fmt.Println(t.Username, "\ttweets about golang")
				} else {
					fmt.Println(t.Username, "\tdoes not tweet about golang")
				}
			}
		default:
			return
		}
	}
}

const RoutinesNum = 4
const ChannelSize = 200

func main() {
	start := time.Now()
	stream := GetMockStream()

	datach := make(chan Tweet, ChannelSize)

	// Producer
	producer(stream, datach)

	// Consumer
	for i := 0; i < RoutinesNum; i++ {
		go consumer(datach, wg)
	}

	fmt.Printf("Process took %s\n", time.Since(start))
}
