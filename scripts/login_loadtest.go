package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	login_baseURL     = "http://localhost:8080/login" // Local login URL
	login_totalUsers  = 10000                         // Total number of users to login
	login_concurrency = 100                           // Number of concurrent goroutines
	login_username    = "load_test_user"                // Base test username
	login_password    = "12345678"                    // Base test password
)

type LoginResult struct {
	success  bool
	duration time.Duration
}

// randomizeCredentials randomly modifies the username or password to simulate errors
func randomizeCredentials(i int) (string, string) {
	// 90% chance of using correct username/password, 10% chance of error
	if rand.Float32() < 0.1 {
		// Generate incorrect credentials
		if rand.Float32() < 0.5 {
			// Incorrect username
			return fmt.Sprintf("wronguser_%d", i), login_password
		}
		// Incorrect password
		return fmt.Sprintf("%s_%d", login_username, i), "wrongpassword"
	}
	// Correct username/password
	return fmt.Sprintf("%s_%d", login_username, i), login_password
}

// loginUser sends a login request and measures the response time
func loginUser(i int) LoginResult {
	// Randomize username and password
	username, password := randomizeCredentials(i)

	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)

	start := time.Now()

	resp, err := http.PostForm(login_baseURL, form)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Println("Error on request:", err)
		return LoginResult{success: false, duration: elapsed}
	}
	defer resp.Body.Close()

	// Successful login should return a redirect or a successful response
	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusOK {
		return LoginResult{success: true, duration: elapsed}
	}
	return LoginResult{success: false, duration: elapsed}
}

func LoginLoadTest() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]LoginResult, 0, login_totalUsers)

	sem := make(chan struct{}, login_concurrency)

	startAll := time.Now()

	// Launch goroutines
	for i := 0; i < login_totalUsers; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(i int) {
			defer wg.Done()
			result := loginUser(i)

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
	fmt.Printf("Total Users Tried: %d\n", login_totalUsers)
	fmt.Printf("Successful Logins: %d\n", successCount)
	fmt.Printf("Failed Logins: %d\n", failCount)
	fmt.Printf("Average Response Time: %v\n", totalTime/time.Duration(login_totalUsers))
	fmt.Printf("Max Response Time: %v\n", maxTime)
	fmt.Printf("Total Test Duration: %v\n", totalElapsed)
	fmt.Printf("Requests per Second: %.2f\n", float64(login_totalUsers)/totalElapsed.Seconds())
}

/*
Total Users Tried: 10000
Successful Logins: 10000
Failed Logins: 0
Average Response Time: 2.39934013s
Max Response Time: 6.6534222s
Total Test Duration: 4m0.9381192s
Requests per Second: 41.50
*/