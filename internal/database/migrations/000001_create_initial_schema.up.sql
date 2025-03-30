-- Create enum types
CREATE TYPE user_role AS ENUM ('normal', 'admin');
CREATE TYPE problem_status AS ENUM ('draft', 'published');
CREATE TYPE submission_status AS ENUM ('OK', 'Compile Error', 'Wrong Answer', 'Memory Limit', 'Time Limit', 'Runtime Error');

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    encrypted_password VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'normal',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Problems table
CREATE TABLE problems (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    statement TEXT NOT NULL,
    time_limit_ms INTEGER NOT NULL DEFAULT 1000,
    memory_limit_mb INTEGER NOT NULL DEFAULT 256,
    sample_input TEXT,
    sample_output TEXT,
    owner_id INTEGER NOT NULL REFERENCES users(id),
    status problem_status NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    modified_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add index on foreign key
CREATE INDEX idx_problems_owner ON problems(owner_id);

-- Test cases table
CREATE TABLE test_cases (
    id SERIAL PRIMARY KEY,
    problem_id INTEGER REFERENCES problems(id),
    input TEXT NOT NULL,
    output TEXT NOT NULL
);

-- Submissions table
CREATE TABLE submissions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    problem_id INTEGER REFERENCES problems(id),
    source_code TEXT NOT NULL,
    status submission_status NOT NULL,
    execution_time_ms INTEGER,
    memory_used_mb INTEGER,
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add index on foreign key
CREATE INDEX idx_submissions_user ON submissions(user_id);

-- Create function for automatic modified_at updates
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.modified_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update modified_at on problems table
CREATE TRIGGER update_problems_modtime
BEFORE UPDATE ON problems
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();