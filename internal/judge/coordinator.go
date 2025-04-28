package judge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/services"
)

type RunnerNode struct {
	ID       string
	Endpoint string
	Healthy  bool
	LastPing time.Time
}

type RunnerCoordinator struct {
	mu     sync.RWMutex
	nodes  []RunnerNode
	client *http.Client
}

func StartCoordinator(ctx context.Context, database *services.DBService) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	runnerCoordinator := NewRunnerCoordinator()

	for {
		select {
		case <-ctx.Done():
			log.Println("Coordinator shutting down...")
			return

		case <-ticker.C:
			submission, err := database.Queries.LockNextPending(ctx)
			if err != nil {
				log.Printf("error fetching submission: %v", err)
				continue
			}

			if submission == nil {
				// No pending submissions
				continue
			}

			// Start processing the submission
			go func(sub *generated.Submission) {
				resp, err := runnerCoordinator.ExecuteSubmission(ctx, sub)
				if err != nil {
					log.Printf("error processing submission %d: %v", sub.ID, err)
					// You might want to mark submission as error
				}
			}(submission)
		}
	}
}

func NewRunnerCoordinator() *RunnerCoordinator {
	return &RunnerCoordinator{
		client: &http.Client{Timeout: 30 * time.Second},
		nodes:  make([]RunnerNode, 0),
	}
}

func (rc *RunnerCoordinator) AddRunner(endpoint string) string {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	id := fmt.Sprintf("runner-%d", len(rc.nodes))
	rc.nodes = append(rc.nodes, RunnerNode{
		ID:       id,
		Endpoint: endpoint,
		Healthy:  true,
		LastPing: time.Now(),
	})
	return id
}

func (rc *RunnerCoordinator) GetAvailableRunner() (*RunnerNode, error) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	for _, node := range rc.nodes {
		if node.Healthy {
			return &node, nil
		}
	}
	return nil, fmt.Errorf("no available runners")
}

func (rc *RunnerCoordinator) ExecuteSubmission(ctx context.Context, submission generated.Submission) (*ExecuteResponse, error) {
	runner, err := rc.GetAvailableRunner()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/execute", runner.Endpoint)
	reqBody, err := json.Marshal(submission)
	if err != nil {
		return nil, fmt.Errorf("error marshaling submission: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := rc.client.Do(req)
	if err != nil {
		rc.markUnhealthy(runner.ID)
		return nil, fmt.Errorf("error sending to runner: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("runner returned status %d", resp.StatusCode)
	}

	var result ExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

func (rc *RunnerCoordinator) markUnhealthy(runnerID string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	for i, node := range rc.nodes {
		if node.ID == runnerID {
			rc.nodes[i].Healthy = false
			rc.nodes[i].LastPing = time.Now()
			return
		}
	}
}

func getRequestFromSubmission(ctx context.Context, database services.DBService, submission generated.Submission) ExecuteRequest {
	// get problme
	problem, err := database.Queries.GetProblemById(ctx, submission.ProblemID.Int32)
	if err != nil {
		log.Printf("error fetching problem: %v", err)
		return ExecuteRequest{}
	}

	return ExecuteRequest{
		Code:           submission.SourceCode,
		Input:          problem.SampleInput.String,
		ExpectedOutput: problem.SampleOutput.String,
		MemoryLimit:    int(problem.MemoryLimitMb),
		TimeLimit:      int(problem.TimeLimitMs),
	}
}
