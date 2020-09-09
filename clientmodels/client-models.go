package clientmodels

import (
	"github.com/tsawler/goblender/pkg/models"
	"time"
)

// Course describes course model
type Course struct {
	ID          int
	CourseName  string
	Active      int
	Lectures    []Lecture
	Description string
	ProfName    string
	ProfEmail   string
	TeamsLink   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Lecture describes a lecture
type Lecture struct {
	ID          int
	LectureName string
	CourseID    int
	VideoID     int
	Video       models.Video
	SortOrder   int
	Active      int
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Assignment holds an assignment
type Assignment struct {
	ID              int
	FileNameDisplay string
	FileName        string
	UserID          int
	CourseID        int
	Mark            int
	TotalValue      int
	Processed       int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	User            models.User
	Course          Course
}
