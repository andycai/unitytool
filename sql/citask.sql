-- 任务表
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL,  -- script, http
    script TEXT,
    url VARCHAR(255),
    method VARCHAR(10),
    headers TEXT,
    body TEXT,
    timeout INTEGER DEFAULT 300,
    status VARCHAR(20) DEFAULT 'inactive',  -- active, inactive
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 任务日志表
CREATE TABLE IF NOT EXISTS task_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL,  -- success, failed, running
    output TEXT,
    error TEXT,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    duration INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
); 