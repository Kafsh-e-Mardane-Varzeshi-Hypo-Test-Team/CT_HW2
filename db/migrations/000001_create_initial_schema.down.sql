-- Drop tables
DROP TABLE IF EXISTS submissions;
DROP TABLE IF EXISTS test_cases;
DROP TABLE IF EXISTS problems;
DROP TABLE IF EXISTS users;

-- Drop triggers
DROP TRIGGER IF EXISTS update_problems_modtime ON problems;

-- Drop functions
DROP FUNCTION IF EXISTS update_modified_column();

-- Drop enum types
DROP TYPE IF EXISTS submission_status;
DROP TYPE IF EXISTS problem_status;
DROP TYPE IF EXISTS user_role;