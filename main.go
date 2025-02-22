package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	numClients := flag.Int("c", 1, "Number of clients to create.")
	numRequests := flag.Int("n", 1000, "Number of requests to make per client.")
	sleep := flag.Duration("s", time.Millisecond, "Time to sleep between requests.")
	rps := flag.Int("r", 1000, "Requests per second to attempt to make per client.")
	url := flag.String("u", "", "URL to make requests to. If not passed, a local built-in server is created and requests are sent to it.")
	flag.Parse()
	totalRequests := *numClients * *numRequests

	fmt.Println("Number of clients: ", *numClients)
	fmt.Println("Number of requests per client: ", *numRequests)
	fmt.Println("Total requests: ", totalRequests)


	var tickTime time.Duration
	if *rps == 0 {	
		tickTime = *sleep
	} else {
		fmt.Println("Requests per second per client: ", *rps)
		tickTime = time.Duration(1e9 / *rps)
	}

	if *url == "" {
		fmt.Println("URL flag not passed, using built-in server.")
		go server()
		*url = "http://localhost:8080"
	}
	
	var statusCounts = make(map[int]*atomic.Uint64)
	var wg sync.WaitGroup

	fmt.Println("Making requests...")
	fmt.Println("")

	startTime := time.Now()
	for range *numClients {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client(*url, *numRequests, tickTime, statusCounts)
		}()
	}
	wg.Wait()
	finishTime := time.Now()

	timeTaken := finishTime.Sub(startTime)
	

	successfulRequests := summariseStatusCounts(statusCounts, totalRequests)
	fmt.Println()
	fmt.Printf("Success: %d/%d\n", successfulRequests, totalRequests)
	fmt.Printf("Time taken: %s\n", timeTaken.String())
	fmt.Printf("Successful requests per second: %f\n", float64(successfulRequests)/finishTime.Sub(startTime).Seconds())
	fmt.Printf("Total requests per second: %f", float64(totalRequests)/finishTime.Sub(startTime).Seconds())
}

func client(url string, numRequests int, sleep time.Duration, statusCounts map[int]*atomic.Uint64) {
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
	defer httpClient.CloseIdleConnections()
	
	t := time.NewTicker(sleep)

	for i := 0; i < numRequests; i++ {
		var resp *http.Response
		req, _ := http.NewRequest("GET", url, nil)
		resp, err := httpClient.Do(req)

		if err == nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			countResponseStatusCode(resp.StatusCode, statusCounts)
		}
		<-t.C
	}
	t.Stop()
}

func countResponseStatusCode(code int, statusCounts map[int]*atomic.Uint64) {
	if _, ok := statusCounts[code]; !ok {
		statusCounts[code] = &atomic.Uint64{}
	}
	statusCounts[code].Add(1)
}

func summariseStatusCounts(statusCounts map[int]*atomic.Uint64, totalRequests int) uint64 {
	var successfulRequests uint64
	
	fmt.Println("Response counts by status code:")
	for code, count := range statusCounts {
		loaded := count.Load()
		if code == 200 {
			successfulRequests = loaded
		}
		fmt.Printf("\t%d: %d/%d (%.2f%%)\n", code, loaded, totalRequests, float64(loaded)/float64(totalRequests)*100)
	}
	return successfulRequests
}

func server() {
	fmt.Println("Starting server on port 8080")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
