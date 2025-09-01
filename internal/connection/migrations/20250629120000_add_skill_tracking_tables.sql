-- +goose Up
CREATE TABLE UserSkills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES UserProfile(id) ON DELETE CASCADE
);

CREATE TABLE ProjectSkill (
    project_id INTEGER NOT NULL,
    skill_id INTEGER NOT NULL,
    PRIMARY KEY (project_id, skill_id),
    FOREIGN KEY (project_id) REFERENCES Projects(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES UserSkills(id) ON DELETE CASCADE
);

CREATE TABLE UserSkillProgress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    skill_id INTEGER NOT NULL,
    total_minutes_tracked INTEGER DEFAULT 0,
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, skill_id),
    FOREIGN KEY (user_id) REFERENCES UserProfile(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES UserSkills(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS UserSkillProgress;
DROP TABLE IF EXISTS ProjectSkill;
DROP TABLE IF EXISTS UserSkills;
