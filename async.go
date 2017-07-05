package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func runAsync(requests []*http.Request, stats *stat) {
	requestsChan = make(chan *http.Request, len(requests))
	for _, req := range requests {
		requestsChan <- req
	}
	for i := int64(0); i < noOfRequest; i++ {
		go sendAsyncRequest(<-requestsChan, stats)
		wg.Add(1)
		time.Sleep(wait)
	}
	wg.Wait()
}

func sendAsyncRequest(req *http.Request, s *stat) {
	defer wg.Done()
	if printReq {
		log.Printf("Request:\n%v\n", req)
	}
	start := time.Now()
	client := &http.Client{}
	resp, err := client.Do(req)
	requestsChan <- req
	if err != nil {
		if !quiet {
			log.Printf("Error 1: %v\n", err)
		}
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error 2: Status: %s (%d). (%s)\n", resp.Status, resp.StatusCode, req.URL)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if !quiet {
			log.Printf("Error 3: %v\n", err)
		}
		return
	}
	if printBody {
		log.Printf("Body:\n%s\n", string(body))
	}
	s.AddTime(time.Now().Sub(start))
	s.Success()

}
