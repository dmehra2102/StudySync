DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS study_sessions;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS notifications;

DROP INDEX IF EXISTS idx_study_sessions_user_id;
DROP INDEX IF EXISTS idx_study_sessions_planned_for;
DROP INDEX IF EXISTS idx_tasks_user_id;
DROP INDEX IF EXISTS idx_tasks_due_date;
DROP INDEX IF EXISTS idx_notifications_user_id;
DROP INDEX IF EXISTS idx_notifications_read;