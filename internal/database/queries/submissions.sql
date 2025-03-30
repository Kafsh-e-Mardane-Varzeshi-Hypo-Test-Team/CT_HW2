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
  $1, $2, $3, $4::submission_status
)
RETURNING *;

-- name: UpdateSubmissionStatus :one
UPDATE submissions
SET status = $2::submission_status, 
    execution_time_ms = $3, 
    memory_used_mb = $4
WHERE id = $1
RETURNING *;