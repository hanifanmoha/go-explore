DROP TABLE IF EXISTS worker_logs;
DROP TABLE IF EXISTS workers;

CREATE TABLE IF NOT EXISTS workers (
    id SERIAL PRIMARY KEY,
    status VARCHAR(50) NOT NULL DEFAULT 'pending'
);

CREATE TABLE IF NOT EXISTS worker_logs (
    id SERIAL PRIMARY KEY,
    worker_id INT NOT NULL,
    finished_at TIMESTAMPTZ,
    worker_name varchar(100),
    FOREIGN KEY (worker_id) REFERENCES workers(id) ON DELETE CASCADE
);

-- Insert 50 rows of pending workers
INSERT INTO workers (status) VALUES
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending'),
    ('pending');