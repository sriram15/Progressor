-- name: AggregateWeekHours :one
SELECT COALESCE(SUM(te.duration), 0) AS totalTrackedMinsCurrentWeek
FROM TimeEntries te
JOIN Cards c ON te.cardId = c.id
JOIN Projects p ON c.projectId = p.id
WHERE p.id = ?
AND strftime('%Y-%W', te.startTime) = strftime('%Y-%W', 'now');

-- name: AggregateMonthHours :one
SELECT COALESCE(SUM(te.duration), 0) AS totalTrackedMinsCurrentMonth
FROM TimeEntries te
JOIN Cards c ON te.cardId = c.id
JOIN Projects p ON c.projectId = p.id
WHERE p.id = ?
AND strftime('%Y-%m', te.startTime) = strftime('%Y-%m', 'now');

-- name: AggregateYearHours :one
SELECT COALESCE(SUM(te.duration), 0) AS totalTrackedMinsCurrentYear
FROM TimeEntries te
JOIN Cards c ON te.cardId = c.id
JOIN Projects p ON c.projectId = p.id
WHERE p.id = ?
AND strftime('%Y', te.startTime) = strftime('%Y', 'now');

