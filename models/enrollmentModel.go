package models

type Enrollment struct {
	ID        int64 `json:"id,omitempty"`
	StudentID int64 `json:"student_id" validate:"required"`
	CourseID  int64 `json:"course_id" validate:"required"`
}
