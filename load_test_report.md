# Load Test Report for KeMV Online Judge Website

## Summary

The KeMV online judge platform was subjected to several load test scenarios, including user signup, login, problem submission, and problem addition after login.

Across all tests, **no failures** were recorded, and the system maintained stable behavior under varying loads.

- The platform successfully handled up to **10,000 concurrent users** for signup and login without errors.
- **Average response times** for signup and login remained under **3 seconds**, while submission response times were significantly faster, maintaining around **1 second** even with high submission volume.
- The **maximum observed response time** across all scenarios remained below **12 seconds**.
- During problem addition after login, performance degraded slightly (average response time ~4.21 seconds), suggesting optimization may be needed there.

Overall, the platform is **stable under current load expectations** but can benefit from minor optimizations for a better user experience, especially during content management actions.

---

## 1. Signup Users Load Test Results

| Total Users Tried | Successful Signups | Failed Signups | Avg Response Time | Max Response Time | Total Duration | Requests per second |
|:------------------|:-------------------|:---------------|:------------------|:------------------|:---------------|:-------------|
| 1000 | 1000 | 0 | 2.253s | 4.862s | 23.53s | 42.49 |
| 5000 | 5000 | 0 | 2.773s | 6.951s | 2m20s | 35.67 |
| 10000 | 10000 | 0 | 2.715s | 8.687s | 4m32s | 36.69 |

---

## 2. Login Users Load Test Results

| Total Users Tried | Successful Logins | Failed Logins | Avg Response Time | Max Response Time | Total Duration | Requests per second |
|:------------------|:------------------|:--------------|:------------------|:------------------|:---------------|:-------------|
| 10000 | 10000 | 0 | 2.399s | 6.653s | 4m1s | 41.50 |

---

## 3. Login and Add Problem Load Test Results

| Total Users Tried | Successful Add Problem | Failed Add Problem | Avg Response Time | Max Response Time | Total Duration | Requests per second |
|:------------------|:-----------------------|:-------------------|:------------------|:------------------|:---------------|:-------------|
| 10000 | 10000 | 0 | 4.210s | 11.020s | 417.08s (~6m57s) | 23.98 |

---

## 4. Signup + Submissions Load Test Results

| Total Users Tried | Total Submissions Tried | Successful Submissions | Failed Submissions | Avg Response Time | Max Response Time | Total Duration | Requests per second |
|:------------------|:------------------------|:-----------------------|:-------------------|:------------------|:------------------|:---------------|:-------------|
| 100 | 500 | 500 | 0 | 438.7ms | 2.472s | 4.23s | 118.22 |
| 100 | 5000 | 5000 | 0 | 1.722s | 7.269s | 1m34s | 53.19 |
| 100 | 50000 | 50000 | 0 | 1.076s | 8.262s | 9m7s | 91.35 |

---

## Final Notes

- **No failed operations** were recorded in any scenario.
- **Signup, login, and submission flows are robust**, even under high load.
- **Problem addition performance** could be optimized to reduce user-perceived latency.
- Future tests could involve distributed load generation and real-world network simulation to gather even more realistic results.

---

_This report is based on load tests executed **locally on a laptop environment** and might improve further under production-grade infrastructure._
