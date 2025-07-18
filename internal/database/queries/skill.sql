-- name: CreateSkill :one
INSERT INTO UserSkills (user_id, name, description) VALUES (?, ?, ?) RETURNING *;

-- name: GetSkillByID :one
SELECT * FROM UserSkills WHERE id = ? LIMIT 1;

-- name: GetSkillsByUserID :many
SELECT * FROM UserSkills WHERE user_id = ?;

-- name: UpdateSkill :one
UPDATE UserSkills SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeleteSkill :exec
DELETE FROM UserSkills WHERE id = ?;
