package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	submit_baseURL            = "http://localhost:8080"
	submit_loginPath          = "/signup"
	submit_submitPath         = "/submit"
	submit_usersToTest        = 100
	submit_submissionsPerUser = 500
)

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	var mu sync.Mutex

	successCount := 0
	failCount := 0
	totalDuration := time.Duration(0)
	maxDuration := time.Duration(0)

	for i := 0; i < submit_usersToTest; i++ {
		wg.Add(1)
		go func(userIdx int) {
			defer wg.Done()

			jar, _ := cookiejar.New(nil)
			client := &http.Client{Jar: jar}

			username := fmt.Sprintf("load_test_user2_%d", userIdx)
			password := "12345678"

			signupData := url.Values{}
			signupData.Set("username", username)
			signupData.Set("password", password)
			signupData.Set("confirm_password", password)

			resp, err := client.PostForm(submit_baseURL+submit_loginPath, signupData)
			if err != nil {
				log.Printf("[User %s] Signup request failed: %v\n", username, err)
				log.Println(signupData)
				mu.Lock()
				failCount += submit_submissionsPerUser
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Printf("[User %s] Login failed with status: %d\n", username, resp.StatusCode)
				mu.Lock()
				failCount += submit_submissionsPerUser
				mu.Unlock()
				return
			}

			// 2. 5 Submissions
			for j := 0; j < submit_submissionsPerUser; j++ {
				problemID := rand.Intn(10) + 1 // random problem id between 1 and 10
				codeContent := fmt.Sprintf("print('Hello from user%d, submission%d')", userIdx, j)

				var b bytes.Buffer
				w := multipart.NewWriter(&b)

				_ = w.WriteField("id", strconv.Itoa(problemID))
				_ = w.WriteField("language", "python")
				_ = w.WriteField("method", "code")
				_ = w.WriteField("code", codeContent)

				w.Close()

				req, err := http.NewRequest("POST", submit_baseURL+submit_submitPath, &b)
				if err != nil {
					log.Printf("[User %s] Failed to create submit request: %v\n", username, err)
					mu.Lock()
					failCount++
					mu.Unlock()
					continue
				}
				req.Header.Set("Content-Type", w.FormDataContentType())

				startTime := time.Now()
				resp, err := client.Do(req)
				duration := time.Since(startTime)

				mu.Lock()
				totalDuration += duration
				if duration > maxDuration {
					maxDuration = duration
				}
				if err != nil {
					log.Printf("[User %s] Submit request failed: %v\n", username, err)
					failCount++
					mu.Unlock()
					continue
				}

				if resp.StatusCode != http.StatusOK {
					log.Printf("[User %s] Submit failed with status: %d\n", username, resp.StatusCode)
					failCount++
				} else {
					successCount++
				}
				mu.Unlock()

				resp.Body.Close()
			}
		}(i)
	}

	wg.Wait()

	totalTime := time.Since(start)
	totalRequests := submit_usersToTest * submit_submissionsPerUser
	averageResponse := totalDuration / time.Duration(totalRequests)

	fmt.Println("========== Submit Load Test Result ==========")
	fmt.Printf("Total Users Tried: %d\n", submit_usersToTest)
	fmt.Printf("Total Submissions Tried: %d\n", totalRequests)
	fmt.Printf("Successful Submissions: %d\n", successCount)
	fmt.Printf("Failed Submissions: %d\n", failCount)
	fmt.Printf("Average Response Time: %v\n", averageResponse)
	fmt.Printf("Max Response Time: %v\n", maxDuration)
	fmt.Printf("Total Test Duration: %v\n", totalTime)
	fmt.Printf("Requests per Second: %.2f\n", float64(totalRequests)/totalTime.Seconds())
	fmt.Println("=============================================")
}

/*
Total Users Tried: 100
Total Submissions Tried: 500
Successful Submissions: 500
Failed Submissions: 0
Average Response Time: 438.705066ms
Max Response Time: 2.4719671s
Total Test Duration: 4.2295501s
Requests per Second: 118.22

Total Users Tried: 100       
Total Submissions Tried: 5000
Successful Submissions: 5000
Failed Submissions: 0
Average Response Time: 1.721620743s
Max Response Time: 7.2688152s
Total Test Duration: 1m34.0095574s
Requests per Second: 53.19

Total Users Tried: 100
Total Submissions Tried: 50000
Successful Submissions: 50000
Failed Submissions: 0
Average Response Time: 1.076032082s
Max Response Time: 8.2615066s
Total Test Duration: 9m7.336797s
Requests per Second: 91.35
*/