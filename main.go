package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	wait                time.Duration
	noOfRequest         int64
	wg                  sync.WaitGroup
	config              string
	quiet, synch        bool
	printBody, printReq bool
	requestsChan        chan *http.Request
	successTimesChan    chan time.Duration
	successTimes        []time.Duration
)

type runFunc func([]*http.Request, *stat)

func init() {
	var (
		no   uint64
		help bool
	)
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	fs.DurationVar(&wait, "t", 100*time.Millisecond, "Time delay between requests")
	fs.StringVar(&config, "f", "data.json", "Path to JSON config file")
	fs.Uint64Var(&no, "n", 1, "Number of requests to send")
	fs.BoolVar(&quiet, "q", false, "Don't print response errors")
	fs.BoolVar(&printBody, "pb", false, "Print response body")
	fs.BoolVar(&printReq, "pr", false, "Print request details")
	fs.BoolVar(&synch, "s", false, "Run requests synchronous. One at time, wait for response, wait 't' before send next")
	fs.BoolVar(&help, "h", false, "Print help")
	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if help {
		fs.PrintDefaults()
		os.Exit(0)
	}
	noOfRequest = int64(no)
}

func main() {
	var (
		conf     *data
		err      error
		requests []*http.Request
		callback runFunc
	)
	start := time.Now()
	if conf, err = readConfig(config); err != nil {
		log.Panic(err)
	}
	if requests, err = conf.requests(); err != nil {
		log.Panic(err)
	}
	stats := newStat(noOfRequest, wait)
	approxTime := stats.CalcApprox()
	rp, rpUnit := stats.CalcRps()
	if synch {
		log.Printf("Send request synchronous")
		callback = runSync
	} else {
		log.Printf("Send request asynchronous")
		callback = runAsync
	}
	log.Printf("No of requests: %v, Delay: %v\n", noOfRequest, wait)
	log.Printf("Requests per %s: %d\n", rpUnit, rp)
	log.Printf("Approx test time : %v\n", approxTime)
	callback(requests, stats)
	log.Printf("Stats: Success requests: %d of %d\n", stats.successNo, stats.total)
	if stats.successNo > 0 {
		log.Printf("Stats: Avg time for request: %0.4f sec.\n", stats.AvgTime())
		log.Printf("Stats: Max request time: %v.\n", stats.MaxTime())
		log.Printf("Stats: Min request time: %v.\n", stats.MinTime())
	}
	log.Printf("Stats: Total time: %v\n", time.Now().Sub(start))
}
