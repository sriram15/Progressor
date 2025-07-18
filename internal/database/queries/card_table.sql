-- name: GetCard :one
SELECT 
    c.id AS card_id,
    c.title,
    c.description,
    c.createdAt,
    c.updatedAt,
    c.status,
    c.completedAt,
    c.isactive,
    c.estimatedMins,
    c.trackedMins,
    c.projectId,
    te.id AS time_entry_id,
    te.startTime,
    te.endTime
FROM 
    Cards c
LEFT JOIN 
    TimeEntries te ON c.id = te.cardId
WHERE 
    c.id = ? AND c.projectId = ?;

-- name: ListCards :many
SELECT *, id AS card_id FROM Cards WHERE projectId = ? AND status = ?;

-- name: CreateCard :exec
INSERT INTO Cards (title, description, status, projectId, estimatedMins) VALUES (?, ?, ?, ?, ?);

-- name: UpdateCard :exec
UPDATE Cards SET title = ?, description = ?, status = ?, completedAt = ?, estimatedMins = ?, trackedMins = ? WHERE id = ?;

-- name: DeleteCard :exec
DELETE FROM Cards WHERE projectId = ? AND id = ?;

-- name: GetActiveCard :one
SELECT id, title, status, projectId FROM Cards WHERE isactive == true;

-- name: UpdateCardActive :exec
UPDATE Cards SET isactive = ?, trackedMins = ? WHERE id = ?;

-- name: CreateTimeEntry :one
INSERT INTO TimeEntries (cardId, startTime, endTime) 
VALUES (?, ?, ?) 
RETURNING *;

-- name: GetActiveTimeEntry :one
 SELECT * FROM TimeEntries WHERE cardId = ? AND startTime == endTime;

-- name: UpdateActiveTimeEntry :exec
UPDATE TimeEntries SET endTime = ?, duration = ? WHERE id = ?;
