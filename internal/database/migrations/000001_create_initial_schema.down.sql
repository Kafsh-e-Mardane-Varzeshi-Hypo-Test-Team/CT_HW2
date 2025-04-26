-- Drop triggers first
DROP TRIGGER IF EXISTS after_submission_insert ON submissions;
DROP TRIGGER IF EXISTS after_submission_update ON submissions;
DROP TRIGGER IF EXISTS create_user_stats_after_user_insert ON users;
DROP TRIGGER IF EXISTS update_problems_modtime ON problems;

-- Drop function
DROP FUNCTION IF EXISTS update_submission_stats();
DROP FUNCTION IF EXISTS initialize_user_stats();
DROP FUNCTION IF EXISTS update_modified_column();

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