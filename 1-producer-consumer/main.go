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
	"time"
	"sync"
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
	for {
		select {

			case t:= <-datach: {
				if t.IsTalkingAboutGo() {
					fmt.Println(t.Username, "\ttweets about golang")
					wg.Add(1)
					go consumer(datach, wg)
				} else {
					fmt.Println(t.Username, "\tdoes not tweet about golang")
					wg.Add(1)
					go consumer(datach, wg)
				}
			}
			default:
				return
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	datach := make(chan Tweet, 200)

	// Producer
	producer(stream, datach)

	// Consumer
	consumer(datach, wg)

	fmt.Printf("Process took %s\n", time.Since(start))
}
