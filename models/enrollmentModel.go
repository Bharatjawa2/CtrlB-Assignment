package models

type Enrollment struct {
    ID        int64 `json:"id"`
    StudentID int64 `json:"student_id"`
    CourseID  int64 `json:"course_id"`
}
