package clientdb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/tsawler/goblender/client/clienthandlers/clientmodels"
	"time"
)

// DBModel holds the database
type DBModel struct {
	DB *sql.DB
}

// AllCourses returns slice of courses (without lectures)
func (m *DBModel) AllCourses() ([]clientmodels.Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := "SELECT id, course_name, active, created_at, updated_at FROM Courses ORDER BY course_name"

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []clientmodels.Course

	for rows.Next() {
		var s clientmodels.Course
		err = rows.Scan(&s.ID, &s.CourseName, &s.Active, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		courses = append(courses, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

// AllActiveCourses returns slice of active courses (without lectures)
func (m *DBModel) AllActiveCourses() ([]clientmodels.Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := "SELECT id, course_name, active, created_at, updated_at FROM courses where active = 1 ORDER BY course_name"

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []clientmodels.Course

	for rows.Next() {
		var s clientmodels.Course
		err = rows.Scan(&s.ID, &s.CourseName, &s.Active, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		courses = append(courses, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

// GetCourse gets a course (for admin) with all lectures
func (m *DBModel) GetCourse(id int) (clientmodels.Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var course clientmodels.Course

	query := "select id, course_name, active, created_at, updated_at from courses where id = $1"

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&course.ID,
		&course.CourseName,
		&course.Active,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	if err != nil {
		fmt.Println("Error getting course")
		fmt.Println(err)
		return course, err
	}

	// get lectures, if any
	query = `select l.id, l.course_id, l.lecture_name, l.video_id, l.active, l.sort_order, l.created_at,
			l.updated_at, v.video_name, v.file_name, v.thumb
			from lectures l
			left join videos v on (l.video_id = v.id)
			where l.course_id = $1 order by l.sort_order`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var lectures []clientmodels.Lecture
	for rows.Next() {
		var l clientmodels.Lecture
		err = rows.Scan(
			&l.ID,
			&l.CourseID,
			&l.LectureName,
			&l.VideoID,
			&l.Active,
			&l.SortOrder,
			&l.CreatedAt,
			&l.UpdatedAt,
			&l.Video.VideoName,
			&l.Video.FileName,
			&l.Video.Thumb,
		)
		if err != nil {
			return course, err
		}

		lectures = append(lectures, l)
	}

	course.Lectures = lectures

	return course, nil
}

// GetCourseForPublic gets a course with only active lectures (for students)
func (m *DBModel) GetCourseForPublic(id int) (clientmodels.Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var course clientmodels.Course

	query := "select id, course_name, active, created_at, updated_at from courses where id = $1"

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&course.ID,
		&course.CourseName,
		&course.Active,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	if err != nil {
		fmt.Println("Error getting course")
		fmt.Println(err)
		return course, err
	}

	// get lectures, if any
	query = `select l.id, l.course_id, l.lecture_name, l.video_id, l.active, l.sort_order, l.created_at,
			l.updated_at, v.video_name, v.file_name, v.thumb
			from lectures l
			left join videos v on (l.video_id = v.id)
			where l.course_id = $1 and l.active = 1 order by l.sort_order`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var lectures []clientmodels.Lecture
	for rows.Next() {
		var l clientmodels.Lecture
		err = rows.Scan(
			&l.ID,
			&l.CourseID,
			&l.LectureName,
			&l.VideoID,
			&l.Active,
			&l.SortOrder,
			&l.CreatedAt,
			&l.UpdatedAt,
			&l.Video.VideoName,
			&l.Video.FileName,
			&l.Video.Thumb,
		)
		if err != nil {
			return course, err
		}

		lectures = append(lectures, l)
	}

	course.Lectures = lectures

	return course, nil
}

// GetLecture returns one lecture
func (m *DBModel) GetLecture(id int) (clientmodels.Lecture, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var l clientmodels.Lecture

	query := `select l.id, l.course_id, l.lecture_name, l.video_id, l.active, l.sort_order, l.created_at,
			l.updated_at, v.video_name, v.file_name, v.thumb
			from lectures l
			left join videos v on (l.video_id = v.id)
			where l.id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&l.ID,
		&l.CourseID,
		&l.LectureName,
		&l.VideoID,
		&l.Active,
		&l.SortOrder,
		&l.CreatedAt,
		&l.UpdatedAt,
		&l.Video.VideoName,
		&l.Video.FileName,
		&l.Video.Thumb,
	)

	if err != nil {
		return l, err
	}

	return l, nil
}

// UpdateCourse updates a course
func (m *DBModel) UpdateCourse(c clientmodels.Course) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update courses set course_name = $1, active = $2, updated_at = $3 where id = $4`

	_, err := m.DB.ExecContext(ctx, query, c.CourseName, c.Active, time.Now(), c.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
