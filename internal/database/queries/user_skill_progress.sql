-- name: GetUserSkillProgress :one
SELECT * FROM UserSkillProgress WHERE user_id = ? AND skill_id = ? LIMIT 1;

-- name: UpsertUserSkillProgress :one
INSERT INTO UserSkillProgress (user_id, skill_id, total_minutes_tracked)
VALUES (?, ?, ?)
ON CONFLICT(user_id, skill_id) DO UPDATE SET
    total_minutes_tracked = total_minutes_tracked + EXCLUDED.total_minutes_tracked,
    last_updated = CURRENT_TIMESTAMP
RETURNING *;
