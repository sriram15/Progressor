-- name: AddProjectSkill :exec
INSERT INTO ProjectSkill (project_id, skill_id) VALUES (?, ?);

-- name: RemoveProjectSkill :exec
DELETE FROM ProjectSkill WHERE project_id = ? AND skill_id = ?;

-- name: GetSkillsForProject :many
SELECT s.* FROM UserSkills s JOIN ProjectSkill ps ON s.id = ps.skill_id WHERE ps.project_id = ?;
