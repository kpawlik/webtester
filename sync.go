package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func runSync(requests []*http.Request, stats *stat) {
	i := int64(0)
	for {
		crf := ""
		for j := 0; j < len(requests); j++ {
			req := requests[j]
			req.Header.Del("Set-Cookie")
			if j == 0 {
				crf = sendSyncRequest(req, stats)
			} else {
				req.Header.Set("Set-Cookie", fmt.Sprintf("csrf_token=%s; Path=/", crf))
			}
			time.Sleep(wait)
		}
		i++
		if i == noOfRequest {
			break
		}
	}
}

func sendSyncRequest(req *http.Request, s *stat) (crf string) {
	if printReq {
		log.Printf("Request:\n%v\n", req)
	}
	start := time.Now()
	client := &http.Client{}
	resp, err := client.Do(req)
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
	crf = string(body)
	if printBody {
		log.Printf("Body:\n%s\n", string(body))
	}
	s.AddTime(time.Now().Sub(start))
	s.Success()
	return

}
