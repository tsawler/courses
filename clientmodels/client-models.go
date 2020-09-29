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
	PostedDate  time.Time
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
	Description     string
	Processed       int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	User            models.User
	Course          Course
}

// CourseAccess records course access
type CourseAccess struct {
	ID        int
	UserID    int
	LectureID int
	CourseID  int
	Duration  int
	CreatedAt time.Time
	UpdatedAt time.Time
	Student   models.User
	Course    Course
	Lecture   Lecture
}

// Student holds a student
type Student struct {
	ID              int
	FirstName       string
	LastName        string
	UserActive      int
	AccessLevel     int
	Email           string
	Password        []byte
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Google2faSecret string
	UseTfa          int
	LoginTypesId    int
	DeletedAt       time.Time
	EmailVerifiedAt time.Time
	Roles           map[string]int
	Avatar          string
	Preferences     map[string]string
	TimeInCourse    int
	Assignments     []Assignment
	Courses         []Section
}

// CourseTraffic holds traffic data
type CourseTraffic struct {
	LectureName string `json:"y"`
	TotalTime   int    `json:"time"`
	TotalViews  int    `json:"views"`
}

// Section holds a section and associated course
type Section struct {
	ID          int
	CourseID    int
	SectionName string
	Active      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Course      Course
	Students    []Student
}
