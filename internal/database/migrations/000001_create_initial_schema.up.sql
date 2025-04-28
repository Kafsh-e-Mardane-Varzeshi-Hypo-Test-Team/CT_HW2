-- Create enum types
CREATE TYPE user_role AS ENUM ('normal', 'admin');
CREATE TYPE problem_status AS ENUM ('draft', 'published');
CREATE TYPE submission_status AS ENUM ('Pending', 'Running', 'Accepted', 'Compile Error', 'Wrong Answer', 'Memory Limit', 'Time Limit', 'Runtime Error');

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

-- User_stats table
CREATE TABLE user_stats (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    total_submissions INTEGER NOT NULL DEFAULT 0,
    total_accepted INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION update_submission_stats() 
RETURNS TRIGGER AS $$
BEGIN
    -- Handle total_submissions
    IF (TG_OP = 'INSERT') THEN
        UPDATE user_stats 
        SET total_submissions = total_submissions + 1,
            updated_at = NOW()
        WHERE user_id = NEW.user_id;
    END IF;

    -- Handle status changes
    IF (TG_OP = 'INSERT' AND NEW.status = 'Accepted') OR
       (TG_OP = 'UPDATE' AND OLD.status != 'Accepted' AND NEW.status = 'Accepted') THEN
        -- Increment total_accepted only if it's their first Accept for this problem
        UPDATE user_stats 
        SET total_accepted = CASE 
                WHEN NOT EXISTS (
                    SELECT 1 FROM submissions 
                    WHERE user_id = NEW.user_id 
                    AND problem_id = NEW.problem_id 
                    AND status = 'Accepted'
                    AND id != NEW.id
                ) THEN total_accepted + 1
                ELSE total_accepted
            END,
            updated_at = NOW()
        WHERE user_id = NEW.user_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for both INSERT and UPDATE
CREATE TRIGGER after_submission_insert
    AFTER INSERT ON submissions
    FOR EACH ROW
    EXECUTE FUNCTION update_submission_stats();

CREATE TRIGGER after_submission_update
    AFTER UPDATE OF status ON submissions
    FOR EACH ROW
    EXECUTE FUNCTION update_submission_stats();

-- Create trigger function to initialize user stats
CREATE OR REPLACE FUNCTION initialize_user_stats() 
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO user_stats (user_id, total_submissions, total_accepted)
    VALUES (NEW.id, 0, 0);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically create user stats when a user is created
CREATE TRIGGER create_user_stats_after_user_insert
    AFTER INSERT ON users
    FOR EACH ROW
    EXECUTE FUNCTION initialize_user_stats();

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