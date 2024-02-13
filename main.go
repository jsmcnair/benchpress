package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	numClients := flag.Int("c", 1, "Number of client to create")
	numRequests := flag.Int("n", 1, "Number of requests to make per client")
	sleep := flag.String("s", "0s", "Time to sleep between requests")
	url := flag.String("u", "", "URL to make requests to")
	flag.Parse()
	totalRequests := *numClients * *numRequests

	fmt.Println("Number of clients: ", *numClients)
	fmt.Println("Number of requests per client: ", *numRequests)
	fmt.Println("Sleep time: ", *sleep)
	fmt.Println("Total requests: ", totalRequests)

	if *url == "" {
		fmt.Println("URL flag not passed, using built-in server.")
		go server()
		*url = "http://localhost:8080"
	}

	sleepTime, err := time.ParseDuration(*sleep)
	if err != nil {
		fmt.Println("Error parsing sleep time: ", err)
		os.Exit(1)
	}

	var counter atomic.Uint64
	var wg sync.WaitGroup

	fmt.Println("Making requests...\n")

	startTime := time.Now()
	for range *numClients {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client(*url, *numRequests, sleepTime, &counter)
		}()
	}
	wg.Wait()
	finishTime := time.Now()

	timeTaken := finishTime.Sub(startTime)

	successfulRequests := counter.Load()
	fmt.Printf("Success: %d/%d\n", successfulRequests, totalRequests)
	fmt.Printf("Time taken: %v\n", timeTaken)
	fmt.Printf("Requests per second: %f", float64(successfulRequests)/finishTime.Sub(startTime).Seconds())
}

func client(url string, numRequests int, sleep time.Duration, counter *atomic.Uint64) {

	// shared HTTP transport and client for efficient connection reuse
	tr := &http.Transport{
		MaxIdleConns:          10,
		IdleConnTimeout:       14 * time.Second,
		ResponseHeaderTimeout: 14 * time.Second,
		DisableKeepAlives:     false,
	}

	httpClient := &http.Client{
		Transport: tr,
		// do not follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for i := 0; i < numRequests; i++ {
		var resp *http.Response
		req, _ := http.NewRequest("GET", url, nil)
		resp, err := httpClient.Do(req)

		if err == nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			counter.Add(1)
		}

		time.Sleep(sleep)
	}
}

func server() {
	fmt.Println("Starting server on port 8080")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "")
	})
	http.ListenAndServe(":8080", nil)
	fmt.Print("Server started on port 8080")
}
