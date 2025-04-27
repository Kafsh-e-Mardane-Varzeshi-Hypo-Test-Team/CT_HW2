-- name: GetSubmissionById :one
SELECT * FROM submissions
WHERE id = $1;

-- name: ListProblemSubmissions :many
SELECT * FROM submissions
WHERE problem_id = $1
ORDER BY submitted_at DESC;

-- TODO: May be better to add problem name as a column to avoid join
-- name: ListUserSubmissions :many
SELECT * FROM submissions AS s
JOIN problems AS p ON p.id = s.problem_id
WHERE s.user_id = $1
ORDER BY s.submitted_at DESC
LIMIT $2 OFFSET $3;

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