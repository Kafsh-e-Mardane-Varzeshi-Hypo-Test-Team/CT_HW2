-- name: GetTestCaseById :one
SELECT * FROM test_cases
WHERE id = $1;

-- name: ListTestCases :many
SELECT * FROM test_cases
WHERE problem_id = $1
ORDER BY id;

-- name: CreateTestCase :one
INSERT INTO test_cases (
  problem_id, input, output
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateTestCase :one
UPDATE test_cases
SET input = $2, output = $3
WHERE id = $1
RETURNING *;

-- name: DeleteTestCase :exec
DELETE FROM test_cases
WHERE id = $1;