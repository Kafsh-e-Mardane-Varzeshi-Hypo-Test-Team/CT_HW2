# Load Test Report for CodeRunner Online Judge

## 1. Test Overview
| Item | Description |
|:-----|:------------|
| **System Under Test (SUT)** | CodeRunner, an online programming contest and learning platform |
| **Objective** | Evaluate system stability and responsiveness under contest-level concurrent usage and high-volume submissions. |
| **Test Date** | 2025-04-27 |
| **Tested By** | QA Team Alpha |

---

## 2. Test Environment
| Component | Description |
|:----------|:------------|
| **Application Server(s)** | 4x AWS EC2 c6i.xlarge (4 vCPU, 8 GB RAM, Ubuntu 22.04) |
| **Database Server(s)** | AWS RDS PostgreSQL (db.m6g.large) |
| **Judging Backend** | 20 containerized workers (Docker) with 2 GB RAM each |
| **Load Testing Tool** | k6 (open-source load testing tool) |
| **Network Environment** | Public AWS VPC (Low latency, ~10 ms between components) |

---

## 3. Test Scenarios
| Scenario | Description | Users/Requests Simulated |
|:---------|:------------|:--------------------------|
| **Login/Signup** | 2000 users logged in within 5 minutes. | 2000 users |
| **Problem Browsing** | 800 users browsing problems simultaneously. | 800 users |
| **Code Submission** | Burst of 300 submissions per minute sustained for 10 minutes. | 300 submissions/min |
| **Judging Workload** | Continuous judging of 250 concurrent submissions. | 250 evaluations |
| **Leaderboard Refresh** | 120 requests per second to leaderboard APIs. | 120 RPS |

---

## 4. Metrics Collected
| Metric | Value | Notes |
|:-------|:------|:------|
| **Max Concurrent Users Handled** | 3000 users |
| **Average Response Time** | 480 ms |
| **Peak Response Time** | 2100 ms (during submission spike) |
| **Throughput** | 700 requests/sec |
| **Error Rate** | 0.8% (mostly 503 Service Unavailable) |
| **CPU Usage** | Min: 25%, Avg: 68%, Max: 92% |
| **Memory Usage** | Min: 3 GB, Avg: 5.7 GB, Max: 7.5 GB |
| **Database Query Latency** | Avg: 45 ms, P95: 110 ms, P99: 180 ms |
| **Judging Queue Wait Time** | Avg: 12 sec, Max: 40 sec |

---

## 5. Key Findings
- **Performance Bottlenecks**: Submission API layer slowed under high load; PostgreSQL connection pool saturation observed.
- **System Failures/Errors Observed**: 503 errors started at 85% CPU utilization; login retries increased at peak.
- **Scaling Behavior**: Auto-scaling of judge workers occurred correctly but database throughput became a bottleneck.

---

## 6. Recommendations
- Increase database connection pool size and optimize slow queries on the submission log table.
- Add horizontal scaling for application servers during peak times.
- Implement caching (e.g., Redis) for frequently accessed problem metadata.
- Prewarm leaderboard cache before contest start.

---

## 7. Graphs & Charts
_(See attached figures.)_
- **Figure 1:** Response time vs number of concurrent users
- **Figure 2:** Error rates over time
- **Figure 3:** CPU/Memory utilization over test duration
- **Figure 4:** Judging queue wait times during submission burst

---

## 8. Conclusion
> **Summary:** CodeRunner successfully handled up to 3000 concurrent users under heavy load with a tolerable error rate and acceptable judging delays. However, performance degradation was noted past 85% CPU utilization, requiring optimizations before hosting events larger than 2500 simultaneous users.

---

### Appendix
- **Test Plan Document**: [Link to Test Plan PDF]
- **Load Generator Configuration**: k6 script and parameters (available in repository)
- **Logs and Raw Data**: Full data sets attached (AWS CloudWatch logs + k6 output)

