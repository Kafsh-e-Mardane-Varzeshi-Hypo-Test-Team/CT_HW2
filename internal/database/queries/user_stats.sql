-- name: UserHasSolvedProblem :one
SELECT EXISTS (
    SELECT 1 FROM submissions 
    WHERE user_id = @user_id
    AND problem_id = @problem_id 
    AND status = 'OK'
);

-- name: GetUserStatsById :one
SELECT 
    total_submissions,
    total_accepted,
    CASE 
        WHEN total_submissions > 0 
        THEN ROUND((total_accepted::NUMERIC / total_submissions::NUMERIC) * 100, 2)
        ELSE 0 
    END as acceptance_rate
FROM user_stats 
WHERE user_id = @user_id;

-- name: GetUserRankingById :many
SELECT 
    u.username,
    us.total_accepted,
    us.total_submissions,
    CASE 
        WHEN us.total_submissions > 0 
        THEN ROUND((us.total_accepted::NUMERIC / us.total_submissions::NUMERIC) * 100, 2)
        ELSE 0 
    END as acceptance_rate
FROM user_stats us
JOIN users u ON u.id = us.user_id
ORDER BY us.total_accepted DESC, acceptance_rate DESC
LIMIT $1;