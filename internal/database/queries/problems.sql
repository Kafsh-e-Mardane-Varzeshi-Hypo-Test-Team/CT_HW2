-- name: GetProblemById :one
SELECT * FROM problems
WHERE id = $1;

-- name: ListPublishedProblems :many
SELECT * FROM problems
WHERE status = 'published'
ORDER BY id;

-- name: ListUserProblems :many
SELECT * FROM problems
WHERE owner_id = $1
ORDER BY id;

-- name: CreateProblem :one
INSERT INTO problems (
  title, statement, time_limit_ms, memory_limit_mb, 
  sample_input, sample_output, owner_id, status
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8::problem_status
)
RETURNING *;

-- name: UpdateProblem :one
UPDATE problems
SET title = $2, 
    statement = $3, 
    time_limit_ms = $4, 
    memory_limit_mb = $5, 
    sample_input = $6, 
    sample_output = $7, 
    status = $8::problem_status
WHERE id = $1
RETURNING *;

-- name: DeleteProblem :exec
DELETE FROM problems
WHERE id = $1;