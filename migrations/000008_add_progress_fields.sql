-- +goose Up
ALTER TABLE lesson_progress
    ADD COLUMN completed_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN last_attempt_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN is_completed BOOLEAN NOT NULL DEFAULT false;

-- Обновляем существующие записи
UPDATE lesson_progress 
SET is_completed = CASE 
    WHEN viewed_at IS NOT NULL AND (NOT requires_test OR passed_test) THEN true 
    ELSE false 
END,
completed_at = CASE 
    WHEN viewed_at IS NOT NULL AND (NOT requires_test OR passed_test) THEN viewed_at 
    ELSE NULL 
END
FROM lessons 
WHERE lesson_progress.lesson_id = lessons.id;

-- +goose Down
ALTER TABLE lesson_progress
    DROP COLUMN completed_at,
    DROP COLUMN last_attempt_at,
    DROP COLUMN is_completed; 