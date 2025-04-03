-- Drop triggers first
DROP TRIGGER IF EXISTS after_submission_insert ON submissions;
DROP TRIGGER IF EXISTS after_submission_update ON submissions;

-- Drop function
DROP FUNCTION IF EXISTS update_submission_stats();

-- Drop tables
DROP TABLE IF EXISTS submissions;
DROP TABLE IF EXISTS test_cases;
DROP TABLE IF EXISTS problems;
DROP TABLE IF EXISTS user_stats;
DROP TABLE IF EXISTS users;

-- Drop enum types
DROP TYPE IF EXISTS submission_status;
DROP TYPE IF EXISTS problem_status;
DROP TYPE IF EXISTS user_role;