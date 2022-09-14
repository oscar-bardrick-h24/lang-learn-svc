package domain

import (
	"fmt"
	"time"
)

type Course struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	LessonIDs []string  `json:"lessons"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type coursePatchableAttributes struct {
	Title     string   `json:"title"`
	LessonIDs []string `json:"lessons"`
}

func (c *Course) getPatchableAttributes() coursePatchableAttributes {
	return coursePatchableAttributes{
		Title:     c.Title,
		LessonIDs: c.LessonIDs,
	}
}

func (c *Course) patchAttributes(p coursePatchableAttributes) error {
	if p.Title == c.Title && equalSlices(p.LessonIDs, c.LessonIDs) {
		return fmt.Errorf("patch effects no change")
	}

	c.Title = p.Title
	c.LessonIDs = p.LessonIDs

	return nil
}

func (c *Course) Validate(idValidator IDValidator, pwordValidator PasswordValidator) error {
	if !idValidator.IsValid(c.ID) {
		return fmt.Errorf("id format is invalid")
	}

	// Title
	switch {
	case c.Title == "":
		return fmt.Errorf("course title must not be empty")
	case len(c.Title) > 255:
		return fmt.Errorf("course title must not be longer than 255 characters")
	}

	// CreatedBy
	if !idValidator.IsValid(c.CreatedBy) {
		return fmt.Errorf("course createdBy ID format is invalid")
	}

	// LessonIDs
	for i, id := range c.LessonIDs {
		if !idValidator.IsValid(id) {
			return fmt.Errorf("course lessonIDs[%d] - lesson ID format is invalid", i)
		}
	}

	return nil
}

func equalSlices(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	for i, v := range x {
		if v != y[i] {
			return false
		}
	}
	return true
}
