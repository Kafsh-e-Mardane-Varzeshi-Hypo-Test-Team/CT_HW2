package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	add_problemBaseURL        = "http://localhost:8080"
	add_problemLoginPath      = "/login"
	add_problemNewProblemPath = "/newproblem"
	add_problemUsername       = "load_test_user_%d"
	add_problemPassword       = "12345678"
	totalAddProblemUsers      = 10000
	concurrencyAddProblem     = 100
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	startTime := time.Now()

	var successCount, failCount int64
	var totalResponseTime time.Duration
	var maxResponseTime time.Duration

	mu := sync.Mutex{}

	semaphore := make(chan struct{}, concurrencyAddProblem)

	for i := 0; i < totalAddProblemUsers; i++ {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(i int) {
			defer func() {
				wg.Done()
				<-semaphore
			}()

			username := fmt.Sprintf(add_problemUsername, i + 100)

			jar, _ := cookiejar.New(nil)
			client := &http.Client{Jar: jar}

			loginStart := time.Now()
			err := login(client, username, add_problemPassword)
			loginDuration := time.Since(loginStart)

			if err != nil {
				mu.Lock()
				log.Printf("login for user %d failed: %v", i + 100, err)
				failCount++
				mu.Unlock()
				return
			}

			problemTitle := fmt.Sprintf("Test Problem %d", i)
			problemStatement := fmt.Sprintf("This is a test statement for problem %d", i)

			addStart := time.Now()
			err = addProblem(client, problemTitle, problemStatement)
			if err != nil {
				mu.Lock()
				log.Printf("add problem for user %d failed: %v", i + 100, err)
				failCount++
				mu.Unlock()
				return
			}
			addDuration := time.Since(addStart)

			totalDuration := loginDuration + addDuration

			mu.Lock()
			successCount++
			totalResponseTime += totalDuration
			if totalDuration > maxResponseTime {
				maxResponseTime = totalDuration
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	totalTestDuration := time.Since(startTime)

	averageResponseTime := time.Duration(0)
	if successCount > 0 {
		averageResponseTime = totalResponseTime / time.Duration(successCount)
	}

	// Print Report
	fmt.Println("===========================")
	fmt.Printf("Total Users Tried: %d\n", totalAddProblemUsers)
	fmt.Printf("Successful Add Problem: %d\n", successCount)
	fmt.Printf("Failed Add Problem: %d\n", failCount)
	fmt.Printf("Average Response Time: %.2fs\n", averageResponseTime.Seconds())
	fmt.Printf("Max Response Time: %.2fs\n", maxResponseTime.Seconds())
	fmt.Printf("Total Test Duration: %.2fs\n", totalTestDuration.Seconds())
	if totalTestDuration.Seconds() > 0 {
		fmt.Printf("Requests per Second: %.2f\n", float64(totalAddProblemUsers)/totalTestDuration.Seconds())
	} else {
		fmt.Println("Requests per Second: N/A")
	}
	fmt.Println("===========================")
}

func login(client *http.Client, username, password string) error {
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	resp, err := client.PostForm(add_problemBaseURL+add_problemLoginPath, data)
	if err != nil {
		log.Printf("postform error: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusOK {
		log.Println(data)
		return fmt.Errorf("login failed with status code %d", resp.StatusCode)
	}

	return nil
}

func addProblem(client *http.Client, title, statement string) error {
	data := url.Values{}
	data.Set("title", title)
	data.Set("statement", statement)
	data.Set("time", strconv.Itoa(rand.Intn(2000)+1000))
	data.Set("memory", strconv.Itoa(rand.Intn(512)+128))
	data.Set("input", "Sample Input")
	data.Set("output", "Sample Output")

	req, err := http.NewRequest("POST", add_problemBaseURL+add_problemNewProblemPath, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("add problem failed with status code %d", resp.StatusCode)
	}

	return nil
}

/*
Total Users Tried: 10000
Successful Add Problem: 10000
Failed Add Problem: 0
Average Response Time: 4.21s
Max Response Time: 11.02s
Total Test Duration: 417.08s
Requests per Second: 23.98
*/