package domain

import "time"

type UserCourse struct {
	UserID         string    `json:"user_id"`
	CourseID       string    `json:"course_id"`
	ActiveLessonID string    `json:"active_lesson_id"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}
