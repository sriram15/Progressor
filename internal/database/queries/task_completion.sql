-- name: CreateTaskCompletion :one
INSERT INTO TaskCompletions (
    cardId,
    userId,
    baseExp,
    timeBonusExp,
    streakBonusExp,
    totalExp
) VALUES (
    ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetTaskCompletion :one
SELECT * FROM TaskCompletions
WHERE cardId = ? AND userId = ?;

-- name: ListTaskCompletionsByUser :many
SELECT * FROM TaskCompletions
WHERE userId = ?
ORDER BY completionTime DESC;

-- name: TotalUserExp :one
SELECT CAST(IFNULL(SUM(totalExp), 0) AS FLOAT) as total_exp FROM TaskCompletions WHERE userId = ?;
