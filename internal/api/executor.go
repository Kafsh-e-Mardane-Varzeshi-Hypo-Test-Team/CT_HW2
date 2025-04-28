package api

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerExecutor struct {
	cli *client.Client
}

func NewDockerExecutor() (*DockerExecutor, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerExecutor{cli: cli}, nil
}

func (e *DockerExecutor) ExecuteGoCode(req ExecuteRequest) (*ExecuteResponse, error) {
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "goexec")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	codePath := filepath.Join(tempDir, "main.go")
	inputPath := filepath.Join(tempDir, "input.txt")
	outputPath := filepath.Join(tempDir, "output.txt")
	compilePath := filepath.Join(tempDir, "compile.txt")
	statsPath := filepath.Join(tempDir, "stats.txt")

	if err := os.WriteFile(codePath, []byte(req.Code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write code: %w", err)
	}
	if err := os.WriteFile(inputPath, []byte(req.Input), 0644); err != nil {
		return nil, fmt.Errorf("failed to write input: %w", err)
	}

	// Modified command to measure time and memory using /usr/bin/time
	containerConfig := &container.Config{
		Image:      "golang:tip-alpine3.21",
		WorkingDir: "/app",
		Cmd: []string{"sh", "-c",
			"go build -o /app/main /app/main.go 2> /app/compile.txt && " +
				"/usr/bin/time -f 'TIME:%e\\nMEM:%M' " + // Format: TIME in seconds, MEM in KB
				"timeout " + fmt.Sprintf("%d", req.TimeLimit) + "s /app/main < /app/input.txt > /app/output.txt 2> /app/stats.txt || echo $? > /app/exitcode",
		},
		AttachStdout: true,
		AttachStderr: true,
	}

	hostConfig := &container.HostConfig{
		AutoRemove:  false,
		NetworkMode: "none",
		Resources: container.Resources{
			Memory:     int64(req.MemoryLimit) * 1024 * 1024,
			MemorySwap: int64(req.MemoryLimit) * 1024 * 1024,
			NanoCPUs:   int64(1e9),
		},
		Binds: []string{tempDir + ":/app"},
	}

	resp, err := e.cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %v", err)
	}
	defer e.cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})

	startTime := time.Now()
	if err := e.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(req.TimeLimit+5)*time.Second)
	defer cancel()

	statusCh, errCh := e.cli.ContainerWait(timeoutCtx, resp.ID, container.WaitConditionNotRunning)

	var containerExitCode int64
	select {
	case err := <-errCh:
		if err != nil {
			return nil, fmt.Errorf("container wait error: %w", err)
		}
	case status := <-statusCh:
		containerExitCode = status.StatusCode
	}

	duration := time.Since(startTime)

	// Parse time and memory from stats
	timeUsed := int(duration.Seconds())
	memoryUsed := 0

	statsBytes, err := os.ReadFile(statsPath)
	if err == nil {
		stats := string(statsBytes)
		if strings.Contains(stats, "TIME:") && strings.Contains(stats, "MEM:") {
			// Extract time (in seconds)
			timeStr := strings.Split(strings.Split(stats, "TIME:")[1], "\n")[0]
			if timeSec, err := strconv.ParseFloat(timeStr, 64); err == nil {
				timeUsed = int(timeSec * 1000) // Convert to milliseconds
			}

			// Extract memory (in KB)
			memStr := strings.Split(strings.Split(stats, "MEM:")[1], "\n")[0]
			if memKB, err := strconv.Atoi(memStr); err == nil {
				memoryUsed = memKB // Memory in KB
			}
		}
	}

	// Check for time limit exceeded
	exitCodePath := filepath.Join(tempDir, "exitcode")
	if exitCodeBytes, err := os.ReadFile(exitCodePath); err == nil {
		if strings.TrimSpace(string(exitCodeBytes)) == "124" {
			return &ExecuteResponse{
				Status:     "TimeLimit",
				TimeUsed:   req.TimeLimit * 1000, // Convert to milliseconds
				MemoryUsed: memoryUsed,
				Error:      "time limit exceeded",
			}, nil
		}
	}

	// Check for memory limit
	inspection, err := e.cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	if inspection.State.OOMKilled {
		return &ExecuteResponse{
			Status:     "MemoryLimit",
			TimeUsed:   timeUsed,
			MemoryUsed: req.MemoryLimit * 1024, // Convert MB to KB
			Error:      "memory limit exceeded",
		}, nil
	}

	// Check compile error
	compileBytes, _ := os.ReadFile(compilePath)
	if len(compileBytes) > 0 {
		return &ExecuteResponse{
			Status:     "CompileError",
			TimeUsed:   0,
			MemoryUsed: 0,
			Error:      string(compileBytes),
		}, nil
	}

	// Check runtime error
	if containerExitCode != 0 {
		outputBytes, _ := os.ReadFile(outputPath)
		return &ExecuteResponse{
			Status:     "RuntimeError",
			TimeUsed:   timeUsed,
			MemoryUsed: memoryUsed,
			Error:      string(outputBytes),
		}, nil
	}

	// Compare output
	outputBytes, _ := os.ReadFile(outputPath)
	actualOutput := string(bytes.TrimSpace(outputBytes))
	expectedOutput := string(bytes.TrimSpace([]byte(req.ExpectedOutput)))

	var status string
	if actualOutput == expectedOutput {
		status = "Accepted"
	} else {
		status = "WrongAnswer"
	}

	return &ExecuteResponse{
		Status:     status,
		TimeUsed:   timeUsed,
		MemoryUsed: memoryUsed,
		Error:      "",
	}, nil
}
