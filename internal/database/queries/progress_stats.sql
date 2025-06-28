-- name: AggregateWeekHours :one
SELECT CAST(IFNULL(SUM(te.duration), 0) AS FLOAT) AS totalTrackedMinsCurrentWeek
FROM TimeEntries te
JOIN Cards c ON te.cardId = c.id
JOIN Projects p ON c.projectId = p.id
WHERE p.id = ?
AND strftime('%Y-%W', te.startTime) = strftime('%Y-%W', 'now');

-- name: AggregatePreviousWeekHours :one
SELECT CAST(IFNULL(SUM(te.duration), 0) AS FLOAT) AS totalTrackedMinsPreviousWeek
FROM TimeEntries te
JOIN Cards c ON te.cardId = c.id
JOIN Projects p ON c.projectId = p.id
WHERE p.id = ?
AND strftime('%Y-%W', te.startTime) = strftime('%Y-%W', 'now', '-7 days');

-- name: AggregateMonthHours :one
SELECT CAST(IFNULL(SUM(te.duration), 0) AS FLOAT) AS totalTrackedMinsCurrentMonth
FROM TimeEntries te
JOIN Cards c ON te.cardId = c.id
JOIN Projects p ON c.projectId = p.id
WHERE p.id = ?
AND strftime('%Y-%m', te.startTime) = strftime('%Y-%m', 'now');

-- name: AggregatePreviousMonthHours :one
SELECT CAST(IFNULL(SUM(te.duration), 0) AS FLOAT) AS totalTrackedMinsPreviousMonth
FROM TimeEntries te
JOIN Cards c ON te.cardId = c.id
JOIN Projects p ON c.projectId = p.id
WHERE p.id = ?
AND strftime('%Y-%m', te.startTime) = strftime('%Y-%m', 'now', 'start of month', '-1 month');

-- name: GetDailyTotalMinutes :many
SELECT
    DATE(startTime) AS date,
    SUM(duration) AS total_minutes
FROM
    TimeEntries
WHERE
    startTime >= DATE('now', '-1 year')
GROUP BY
    DATE(startTime)
ORDER BY
    DATE(startTime);

-- name: GetWeeklyProgress :one
SELECT COUNT(DISTINCT DATE(startTime)) AS progress_days
FROM TimeEntries
WHERE strftime('%Y-%W', startTime) = strftime('%Y-%W', 'now');

-- name: GetPreviousWeeklyProgress :one
SELECT COUNT(DISTINCT DATE(startTime)) AS progress_days
FROM TimeEntries
WHERE strftime('%Y-%W', startTime) = strftime('%Y-%W', 'now', '-7 days');

-- name: GetMonthlyProgress :one
SELECT COUNT(DISTINCT DATE(startTime)) AS progress_days
FROM TimeEntries
WHERE strftime('%Y-%m', startTime) = strftime('%Y-%m', 'now');

-- name: GetPreviousMonthlyProgress :one
SELECT COUNT(DISTINCT DATE(startTime)) AS progress_days
FROM TimeEntries
WHERE strftime('%Y-%m', startTime) = strftime('%Y-%m', 'now', 'start of month', '-1 month');
