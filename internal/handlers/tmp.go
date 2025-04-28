package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	admin = gin.H{
		"Username": "admin",
		"IsAdmin":  true,
	}
	problems = []gin.H{
		{
			"ID":          1,
			"Title":       "Problem1",
			"Owner":       "admin",
			"Status":      "Draft",
			"TimeLimit":   1000,
			"MemoryLimit": 256,
			"Statement": `
## Problem Statement

Given an array of integers 'nums' and an integer 'target', return indices of the two numbers such that they add up to 'target'.

You may assume that each input would have exactly one solution, and you may not use the same element twice.

## Input

- The first line contains an integer n (2 ≤ n ≤ 10^4), the size of the array
- The second line contains n space-separated integers representing the array elements (-10⁹ ≤ nums[i] ≤ 10⁹)
- The third line contains the target integer (-10⁹ ≤ target ≤ 10⁹)

## Output

Print two space-separated integers representing the indices of the two numbers.

## Constraints

- Only one valid answer exists
- You cannot use the same element twice

## Example:

**Input:**  

4

2 7 11 15

9

**Output:**  

0 1
`,
		},
		{
			"ID":          2,
			"Title":       "Problem2",
			"Owner":       "admin",
			"Status":      "Published",
			"TimeLimit":   1000,
			"MemoryLimit": 256,
			"Statement":   "statement",
		},
		{
			"ID":          3,
			"Title":       "Problem3",
			"Owner":       "admin",
			"Status":      "Published",
			"TimeLimit":   1000,
			"MemoryLimit": 256,
			"Statement":   "statement",
		},
		{
			"ID":          4,
			"Title":       "Problem4",
			"Owner":       "admin",
			"Status":      "Draft",
			"TimeLimit":   1000,
			"MemoryLimit": 256,
			"Statement":   "statement",
		},
	}
	submissions = []gin.H{
		{
			"ID":      22401,
			"When":    "03-30-2025 20:30",
			"Problem": problems[1],
			"Status":  "Pending",
			"Time":    121,
			"Memory":  19,
		},
		{
			"ID":      22349,
			"When":    "03-29-2025 15:43",
			"Problem": problems[0],
			"Status":  "Accepted",
			"Time":    234,
			"Memory":  42,
		},
		{
			"ID":      22320,
			"When":    "03-29-2025 15:34",
			"Problem": problems[0],
			"Status":  "Wrong Answer",
			"Time":    223,
			"Memory":  42,
		},
		{
			"ID":      22311,
			"When":    "03-29-2025 15:32",
			"Problem": problems[0],
			"Status":  "Compile Error",
			"Time":    0,
			"Memory":  0,
		},
	}
)

func (h *Handler) SubmissionsPage(c *gin.Context) {
	currentPage, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.Redirect(http.StatusFound, "/submissions?page=1")
		return
	}
	c.HTML(http.StatusOK, "submissions.html", gin.H{
		"Submissions": submissions,
		"CurrentPage": currentPage,
		"TotalPages":  3,
	})
}

func (h *Handler) SubmitPage(c *gin.Context) {
	problemID := c.Param("id")
	c.HTML(http.StatusOK, "submit.html", gin.H{
		"ID": problemID,
	})
}
