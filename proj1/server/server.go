package server

import (
	"encoding/json"
	"log"
	L "proj1/feed"
	Q "proj1/queue"
	S "proj1/semaphore"
	"runtime"
	"sync"
)

type Config struct {
	Encoder *json.Encoder // Represents the buffer to encode Responses
	Decoder *json.Decoder // Represents the buffer to decode Requests
	Mode    string        // Represents whether the server should execute
	// sequentially or in parallel
	// If Mode == "s"  then run the sequential version
	// If Mode == "p"  then run the parallel version
	// These are the only values for Version
	ConsumersCount int // Represents the number of consumers to spawn
}

//Run starts up the twitter server based on the configuration
//information provided and only returns when the server is fully
// shutdown.
func Run(config Config) {
	feed := L.NewFeed()
	if config.Mode == "s" {
		ctr := 0
		for {

			var v Q.Request
			if err := config.Decoder.Decode(&v); err != nil {
				log.Println(err)
				return
			}
			if v.Command == "DONE" {
				ctr++
				break
			}
			if ctr > 0 {
				break
			}
			encHelper(v, feed, config)
		}
	} else {
		Queue := Q.NewLockFreeQueue()
		sema := S.NewSemaphore(0)
		jobDone := 0
		var wg sync.WaitGroup
		for i := 0; i < config.ConsumersCount; i++ {
			wg.Add(1)
			go Consumer(config, Queue, sema, feed, &jobDone)
		}
		Producer(config, Queue, sema, &jobDone)
		for runtime.NumGoroutine() > 1 {
			sema.Up()
		}
	}
}

func Consumer(config Config, Queue *Q.LockFreeQueue, sema *S.Semaphore, feed L.Feed, jobDone *int) {
	for {
		sema.Down()
		v := Queue.Dequeue()
		if *jobDone == 1 && v == nil {
			break
		}
		encHelper(*v, feed, config)
	}
}

func Producer(config Config, Queue *Q.LockFreeQueue, sema *S.Semaphore, jobDone *int) {
	for {
		var v Q.Request
		if err := config.Decoder.Decode(&v); err != nil {
			log.Println(err)
			return
		}
		if v.Command == "DONE" {
			*jobDone++
			break
		}
		Queue.Enqueue(&v)
		sema.Up()
	}
}

func encHelper(v Q.Request, feed L.Feed, config Config) {
	resp := true
	r := make(map[string]interface{})
	if v.Command == "ADD" {
		feed.Add(v.Body, v.TimeStamp)
	} else if v.Command == "REMOVE" {
		resp = feed.Remove(v.TimeStamp)
	} else if v.Command == "CONTAINS" {
		resp = feed.Contains(v.TimeStamp)
	} else if v.Command == "FEED" {
		r["id"] = v.Id
		r["feed"] = feed.GetFeed()
	}
	if v.Command != "FEED" {
		r["success"] = resp
		r["id"] = v.Id
	}
	if err := config.Encoder.Encode(&r); err != nil {
		log.Println(err)
	}
}
