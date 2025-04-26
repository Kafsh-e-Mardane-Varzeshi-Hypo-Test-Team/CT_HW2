-- name: GetProblemById :one
SELECT * FROM problems
WHERE id = $1;

-- name: ListPublishedProblems :many
SELECT * FROM problems
WHERE status = 'published'
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: ListUserProblems :many
SELECT * FROM problems
WHERE owner_id = $1
ORDER BY id;

-- TODO: optimize this query
-- name: GetPublishedProblemCount :one
SELECT COUNT(*) FROM problems
WHERE status = 'published';

-- name: CreateProblem :one
INSERT INTO problems (
  title, statement, time_limit_ms, memory_limit_mb, 
  sample_input, sample_output, owner_id, status
) VALUES (
  @title, @statement, @time_limit_ms, @memory_limit_mb, 
  @sample_input, @sample_output, @owner_id, @status::problem_status
)
RETURNING *;

-- name: UpdateProblem :one
UPDATE problems
SET title = @title, 
    statement = @statement, 
    time_limit_ms = @time_limit_ms, 
    memory_limit_mb = @memory_limit_mb, 
    sample_input = @sample_input, 
    sample_output = @sample_output, 
    status = @status::problem_status
WHERE id = @id
RETURNING *;

-- name: DeleteProblem :exec
DELETE FROM problems
WHERE id = $1;