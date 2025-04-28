package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	signup_baseURL      = "http://localhost:8080/signup" // Local signup URL
	signup_totalUsers   = 10000                          // Total number of users to create
	signup_concurrency  = 100                            // Number of concurrent goroutines
	signup_usernameBase = "load_test_user"               // Base username
	signup_password     = "12345678"                     // Password for all users
)

type SignupResult struct {
	success  bool
	duration time.Duration
}

func signupUser(i int) SignupResult {
	form := url.Values{}
	form.Add("username", fmt.Sprintf("%s_%d", signup_usernameBase, i))
	form.Add("password", signup_password)
	form.Add("confirm_password", signup_password)

	start := time.Now()

	resp, err := http.PostForm(signup_baseURL, form)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Println("Error on request:", err)
		return SignupResult{success: false, duration: elapsed}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusOK {
		return SignupResult{success: true, duration: elapsed}
	}
	return SignupResult{success: false, duration: elapsed}
}

func SignupLoadTest() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]SignupResult, 0, signup_totalUsers)

	sem := make(chan struct{}, signup_concurrency)

	startAll := time.Now()

	for i := 0; i < signup_totalUsers; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(i int) {
			defer wg.Done()
			result := signupUser(i)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()

			<-sem
		}(i)
	}

	wg.Wait()
	totalElapsed := time.Since(startAll)

	// Calculate and print results
	var successCount, failCount int
	var totalTime time.Duration
	var maxTime time.Duration

	for _, res := range results {
		if res.success {
			successCount++
		} else {
			failCount++
		}
		totalTime += res.duration
		if res.duration > maxTime {
			maxTime = res.duration
		}
	}

	fmt.Println("---- Load Test Results ----")
	fmt.Printf("Total Users Tried: %d\n", signup_totalUsers)
	fmt.Printf("Successful Signups: %d\n", successCount)
	fmt.Printf("Failed Signups: %d\n", failCount)
	fmt.Printf("Average Response Time: %v\n", totalTime/time.Duration(signup_totalUsers))
	fmt.Printf("Max Response Time: %v\n", maxTime)
	fmt.Printf("Total Test Duration: %v\n", totalElapsed)
	fmt.Printf("Requests per Second: %.2f\n", float64(signup_totalUsers)/totalElapsed.Seconds())
}

/*
---- Load Test Results ----
Total Users Tried: 1000
Successful Signups: 1000
Failed Signups: 0
Average Response Time: 2.253281547s
Max Response Time: 4.861686s
Total Test Duration: 23.5331112s
Requests per Second: 42.49

---- Load Test Results ----
Total Users Tried: 5000
Successful Signups: 5000
Failed Signups: 0
Average Response Time: 2.772791335s
Max Response Time: 6.9510831s
Total Test Duration: 2m20.175024s
Requests per Second: 35.67

---- Load Test Results ----
Total Users Tried: 10000
Successful Signups: 10000
Failed Signups: 0
Average Response Time: 2.714650333s
Max Response Time: 8.6867526s
Total Test Duration: 4m32.5249488s
Requests per Second: 36.69
*/
