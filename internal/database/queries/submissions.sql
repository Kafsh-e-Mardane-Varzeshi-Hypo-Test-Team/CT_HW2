-- name: GetSubmissionById :one
SELECT * FROM submissions
WHERE id = $1;

-- name: ListProblemSubmissions :many
SELECT * FROM submissions
WHERE problem_id = $1
ORDER BY submitted_at DESC;

-- name: ListUserSubmissions :many
SELECT * FROM submissions
WHERE user_id = $1
ORDER BY submitted_at DESC;

-- name: CreateSubmission :one
INSERT INTO submissions (
  user_id, problem_id, source_code, status
) VALUES (
  @user_id, @problem_id, @source_code, @status::submission_status
)
RETURNING *;

-- name: UpdateSubmissionStatusTimeMemory :one
UPDATE submissions
SET status = @status::submission_status, 
    execution_time_ms = @execution_time_ms, 
    memory_used_mb = @memory_used_mb
WHERE id = @id
RETURNING *;