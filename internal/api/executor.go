package api

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerExecutor struct {
	cli *client.Client
}

func NewDockerExecutor() (*DockerExecutor, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
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
	if err := os.WriteFile(codePath, []byte(req.Code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write code file: %v", err)
	}

	resp, err := e.cli.ContainerCreate(ctx, &container.Config{
		Image: "golang:latest",
		Cmd:   []string{"sh", "-c", "go run /app/main.go"},
		Tty:   false,
	}, &container.HostConfig{
		Binds: []string{tempDir + ":/app"},
		Resources: container.Resources{
			Memory:     int64(req.MemoryLimit * 1024 * 1024), // Convert MB to bytes
			MemorySwap: int64(req.MemoryLimit * 1024 * 1024),
			NanoCPUs:   int64(1e9), // 1 CPU
		},
	}, nil, nil, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %v", err)
	}

	defer e.cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})

	timeout := time.Duration(req.TimeLimit) * time.Millisecond
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := e.cli.ContainerStart(ctxWithTimeout, resp.ID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container: %v", err)
	}

	statusCh, errCh := e.cli.ContainerWait(ctxWithTimeout, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return &ExecuteResponse{
				Status: "runtime_error",
				Error:  err.Error(),
			}, nil
		}
	case <-statusCh:
	}

	out, err := e.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get container logs: %v", err)
	}
	defer out.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	output := buf.String()

	status := "success"
	if output != req.ExpectedOutput {
		status = "wrong_answer"
	}

	return &ExecuteResponse{
		Status:     status,
		TimeUsed:   req.TimeLimit, // Simplified - should get actual usage
		MemoryUsed: 0,             // Should get actual usage
	}, nil
}
