-- +goose Up
CREATE TABLE TaskCompletions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cardId INTEGER NOT NULL,
    userId INTEGER NOT NULL,
    completionTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    baseExp INTEGER NOT NULL,
    timeBonusExp INTEGER NOT NULL,
    streakBonusExp INTEGER NOT NULL,
    totalExp INTEGER NOT NULL,
    FOREIGN KEY (cardId) REFERENCES Cards(id) ON DELETE CASCADE,
    FOREIGN KEY (userId) REFERENCES UserProfile(id) ON DELETE CASCADE
);

-- Add columns to UserProfile table
ALTER TABLE UserProfile ADD COLUMN progressionPoints INTEGER NOT NULL DEFAULT 0;
ALTER TABLE UserProfile ADD COLUMN archerLevel INTEGER NOT NULL DEFAULT 1;
ALTER TABLE UserProfile ADD COLUMN archerExperience INTEGER NOT NULL DEFAULT 0;
ALTER TABLE UserProfile ADD COLUMN idleResources INTEGER NOT NULL DEFAULT 0;
ALTER TABLE UserProfile ADD COLUMN lastLoginTimeStamp TIMESTAMP DEFAULT NULL;

-- Create ArcherStats table
CREATE TABLE ArcherStats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userId INTEGER NOT NULL,
    statName TEXT NOT NULL,
    statValue INTEGER NOT NULL,
    FOREIGN KEY (userId) REFERENCES UserProfile(id) ON DELETE CASCADE,
    UNIQUE(userId, statName)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS TaskCompletions;

ALTER TABLE UserProfile DROP COLUMN progressionPoints;
ALTER TABLE UserProfile DROP COLUMN archerLevel;
ALTER TABLE UserProfile DROP COLUMN archerExperience;
ALTER TABLE UserProfile DROP COLUMN idleResources;
ALTER TABLE UserProfile DROP COLUMN lastLoginTimeStamp;

DROP TABLE IF EXISTS ArcherStats;
