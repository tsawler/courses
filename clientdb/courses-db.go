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
	stmt := `SELECT id, course_name, active, 
		prof_name, prof_email, teams_link,
		created_at, updated_at FROM Courses ORDER BY course_name`

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []clientmodels.Course

	for rows.Next() {
		var s clientmodels.Course
		err = rows.Scan(
			&s.ID,
			&s.CourseName,
			&s.Active,
			&s.ProfName,
			&s.ProfEmail,
			&s.TeamsLink,
			&s.CreatedAt,
			&s.UpdatedAt)
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
	stmt := `SELECT id, course_name, active, 
		prof_name, prof_email, teams_link,
		created_at, updated_at FROM courses where active = 1 ORDER BY course_name`

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []clientmodels.Course

	for rows.Next() {
		var s clientmodels.Course
		err = rows.Scan(
			&s.ID,
			&s.CourseName,
			&s.Active,
			&s.ProfName,
			&s.ProfEmail,
			&s.TeamsLink,
			&s.CreatedAt,
			&s.UpdatedAt)
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

	query := `select id, course_name, active, description, 
		prof_name, prof_email, teams_link,
		created_at, updated_at from courses where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&course.ID,
		&course.CourseName,
		&course.Active,
		&course.Description,
		&course.ProfName,
		&course.ProfEmail,
		&course.TeamsLink,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	if err != nil {
		fmt.Println("Error getting course")
		fmt.Println(err)
		return course, err
	}

	// get lectures, if any
	query = `select l.id, l.course_id, l.lecture_name, coalesce(l.video_id, 0), l.active, l.sort_order, l.notes, l.created_at,
			l.updated_at, coalesce(v.video_name, ''), coalesce(v.file_name, ''), coalesce(v.thumb, '')
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
			&l.Notes,
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

	query := `select id, course_name, active, 
		prof_name, prof_email, teams_link,
		description, created_at, updated_at from courses where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&course.ID,
		&course.CourseName,
		&course.Active,
		&course.Description,
		&course.ProfName,
		&course.ProfEmail,
		&course.TeamsLink,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	if err != nil {
		fmt.Println("Error getting course")
		fmt.Println(err)
		return course, err
	}

	// get lectures, if any
	query = `select l.id, l.course_id, l.lecture_name, coalesce(l.video_id, 0), l.active, l.sort_order, l.notes, l.created_at,
			l.updated_at, coalesce(v.video_name, ''), coalesce(v.file_name, ''), coalesce(v.thumb, '')
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
			&l.Notes,
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

	query := `select l.id, l.course_id, l.lecture_name, coalesce(l.video_id, 0), l.active, l.sort_order, l.notes, l.created_at,
			l.updated_at, coalesce(v.video_name, ''), coalesce(v.file_name, ''), coalesce(v.thumb, '')
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
		&l.Notes,
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

	query := `update courses set course_name = $1, active = $2, description = $3, 
		prof_name = $4, prof_email = $5, teams_link = $6, updated_at = $7 where id = $8`

	_, err := m.DB.ExecContext(ctx, query,
		c.CourseName,
		c.Active,
		c.Description,
		c.ProfName,
		c.ProfEmail,
		c.TeamsLink,
		time.Now(),
		c.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// InsertCourse inserts a course and returns new id
func (m *DBModel) InsertCourse(c clientmodels.Course) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	query := `insert into courses (course_name, active, description, prof_name, prof_email, teams_link, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	err := m.DB.QueryRowContext(ctx, query,
		c.CourseName,
		c.Active,
		c.Description,
		c.ProfName,
		c.ProfEmail,
		c.TeamsLink,
		time.Now(),
		time.Now()).Scan(&newID)

	if err != nil {
		fmt.Println("Error inserting new course")
		fmt.Println(err)
		return 0, err
	}

	return newID, nil
}

// InsertCourse inserts a course lecture and returns new id
func (m *DBModel) InsertLecture(c clientmodels.Lecture) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	if c.VideoID > 0 {

		query := `insert into lectures (course_id, lecture_name, video_id, active, sort_order, notes, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`

		err := m.DB.QueryRowContext(ctx, query, c.CourseID, c.LectureName, c.VideoID, c.Active, c.SortOrder, c.Notes, time.Now(), time.Now()).Scan(&newID)

		if err != nil {
			fmt.Println("Error inserting new course lecture")
			fmt.Println(err)
			return 0, err
		}
	} else {
		query := `insert into lectures (course_id, lecture_name, active, sort_order, notes, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7) returning id`

		err := m.DB.QueryRowContext(ctx, query, c.CourseID, c.LectureName, c.Active, c.SortOrder, c.Notes, time.Now(), time.Now()).Scan(&newID)

		if err != nil {
			fmt.Println("Error inserting new course lecture")
			fmt.Println(err)
			return 0, err
		}
	}

	return newID, nil
}

// UpdateLecture updates a course lecture
func (m *DBModel) UpdateLecture(c clientmodels.Lecture) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if c.VideoID > 0 {
		query := `update lectures set lecture_name = $1, video_id = $2, active = $3, notes = $4,
			updated_at = $5 where id = $6`

		_, err := m.DB.ExecContext(ctx, query, c.LectureName, c.VideoID, c.Active, c.Notes, time.Now(), c.ID)

		if err != nil {
			fmt.Println("Error updating course lecture")
			fmt.Println(err)
			return err
		}
	} else {
		query := `update lectures set lecture_name = $1, video_id = null, active = $2, notes = $3,
			updated_at = $4 where id = $5`

		_, err := m.DB.ExecContext(ctx, query, c.LectureName, c.Active, c.Notes, time.Now(), c.ID)

		if err != nil {
			fmt.Println("Error updating course lecture")
			fmt.Println(err)
			return err
		}
	}

	return nil
}

// UpdateLectureSortOrder updates sort order
func (m *DBModel) UpdateLectureSortOrder(id, order int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		update lectures set 
			sort_order = $1, 
			updated_at = $2
		where 
			id = $3`

	_, err := m.DB.ExecContext(ctx, stmt,
		order,
		time.Now(),
		id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateLectureContent updates a course lecture content (notes)
func (m *DBModel) UpdateLectureContent(c clientmodels.Lecture) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update lectures set notes = $1, updated_at = $2 where id = $3`

	_, err := m.DB.ExecContext(ctx, query, c.Notes, time.Now(), c.ID)

	if err != nil {
		fmt.Println("Error updating course lecture")
		fmt.Println(err)
		return err
	}

	return nil
}

// UpdateCourseContent updates a course content (description)
func (m *DBModel) UpdateCourseContent(c clientmodels.Course) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update courses set description = $1, updated_at = $2 where id = $3`

	_, err := m.DB.ExecContext(ctx, query, c.Description, time.Now(), c.ID)

	if err != nil {
		fmt.Println("Error updating course html!")
		fmt.Println(err)
		return err
	}

	return nil
}
