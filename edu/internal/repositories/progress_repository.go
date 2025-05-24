package repositories

import (
	"context"
	"course2/internal/models"
	"database/sql"

	"github.com/google/uuid"
)

type ProgressRepository struct {
	db *sql.DB
}

func NewProgressRepository(db *sql.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

func (r *ProgressRepository) CreateLessonProgress(ctx context.Context, progress *models.LessonProgress) error {
	query := `
		INSERT INTO lesson_progress (
			id, user_id, lesson_id, viewed_at, test_score,
			passed_test, completed_at, last_attempt_at, is_completed,
			attempts_count, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW()
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		progress.ID, progress.UserID, progress.LessonID,
		progress.ViewedAt, progress.TestScore, progress.PassedTest,
		progress.CompletedAt, progress.LastAttemptAt, progress.IsCompleted,
		progress.AttemptsCount,
	)

	return err
}

func (r *ProgressRepository) GetLessonProgress(ctx context.Context, userID, lessonID uuid.UUID) (*models.LessonProgress, error) {
	query := `
		SELECT id, user_id, lesson_id, viewed_at, test_score,
			   passed_test, completed_at, last_attempt_at, is_completed,
			   attempts_count, created_at, updated_at
		FROM lesson_progress
		WHERE user_id = $1 AND lesson_id = $2
	`

	progress := &models.LessonProgress{}
	err := r.db.QueryRowContext(ctx, query, userID, lessonID).Scan(
		&progress.ID, &progress.UserID, &progress.LessonID,
		&progress.ViewedAt, &progress.TestScore, &progress.PassedTest,
		&progress.CompletedAt, &progress.LastAttemptAt, &progress.IsCompleted,
		&progress.AttemptsCount, &progress.CreatedAt, &progress.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return progress, nil
}

func (r *ProgressRepository) UpdateLessonProgress(ctx context.Context, progress *models.LessonProgress) error {
	query := `
		UPDATE lesson_progress
		SET viewed_at = $1, test_score = $2, passed_test = $3,
			completed_at = $4, last_attempt_at = $5, is_completed = $6,
			attempts_count = attempts_count + 1, updated_at = NOW()
		WHERE id = $7 AND user_id = $8 AND lesson_id = $9
	`

	result, err := r.db.ExecContext(ctx, query,
		progress.ViewedAt, progress.TestScore, progress.PassedTest,
		progress.CompletedAt, progress.LastAttemptAt, progress.IsCompleted,
		progress.ID, progress.UserID, progress.LessonID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *ProgressRepository) GetCourseProgress(ctx context.Context, userID, courseID uuid.UUID) (*models.CourseProgress, error) {
	query := `
		WITH course_stats AS (
			SELECT 
				COUNT(*) as total_lessons,
				COUNT(CASE WHEN lp.viewed_at IS NOT NULL AND 
					(NOT l.requires_test OR lp.passed_test) THEN 1 END) as completed_lessons,
				COALESCE(SUM(CASE WHEN xe.type = 'lesson_view' THEN xe.amount ELSE 0 END), 0) +
				COALESCE(SUM(CASE WHEN xe.type = 'test_pass' THEN xe.amount ELSE 0 END), 0) as xp_earned,
				MAX(lp.completed_at) as last_completed_at
			FROM lessons l
			LEFT JOIN lesson_progress lp ON l.id = lp.lesson_id AND lp.user_id = $1
			LEFT JOIN xp_entries xe ON xe.user_id = $1 AND 
				(xe.lesson_id = l.id OR xe.course_id = l.course_id)
			WHERE l.course_id = $2
		)
		SELECT 
			$2 as course_id,
			total_lessons,
			completed_lessons,
			CASE 
				WHEN total_lessons = 0 THEN 0
				ELSE ROUND((completed_lessons::float / total_lessons::float) * 100)
			END as percentage,
			xp_earned,
			CASE 
				WHEN completed_lessons = total_lessons AND total_lessons > 0 THEN 
					last_completed_at
				ELSE NULL
			END as completed_at
		FROM course_stats
	`

	progress := &models.CourseProgress{}
	err := r.db.QueryRowContext(ctx, query, userID, courseID).Scan(
		&progress.CourseID, &progress.TotalLessons, &progress.CompletedLessons,
		&progress.Percentage, &progress.XPEarned, &progress.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return progress, nil
}

func (r *ProgressRepository) AddXPEntry(ctx context.Context, entry *models.XPEntry) error {
	query := `
		INSERT INTO xp_entries (
			id, user_id, course_id, lesson_id, type,
			amount, earned_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, NOW()
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		entry.ID, entry.UserID, entry.CourseID, entry.LessonID,
		entry.Type, entry.Amount,
	)

	return err
}

func (r *ProgressRepository) GetTotalXP(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM xp_entries
		WHERE user_id = $1
	`

	var totalXP int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&totalXP)
	if err != nil {
		return 0, err
	}

	return totalXP, nil
}
