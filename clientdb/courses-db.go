package clientdb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/tsawler/goblender/client/clienthandlers/clientmodels"
	"strings"
	"time"
)

// DBModel holds the database
type DBModel struct {
	DB *sql.DB
}

// AllSections returns slice of all sections
func (m *DBModel) AllSections() ([]clientmodels.Section, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT s.id, s.section_name, s.active, 
		s.course_id, s.created_at, s.updated_at, c.id as course_id, c.course_name, c.active,
		c.prof_name, c.prof_email, c.teams_link,
		c.created_at as course_created_at, c.updated_at as course_updated_at
		FROM course_sections s
		left join courses c on (s.course_id = c.id)
		ORDER BY s.created_at desc`

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []clientmodels.Section

	for rows.Next() {
		var s clientmodels.Section
		err = rows.Scan(
			&s.ID,
			&s.SectionName,
			&s.Active,
			&s.CourseID,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.Course.ID,
			&s.Course.CourseName,
			&s.Course.Active,
			&s.Course.ProfName,
			&s.Course.ProfEmail,
			&s.Course.TeamsLink,
			&s.Course.CreatedAt,
			&s.Course.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sections = append(sections, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sections, nil
}

// UpdateSection updates a course section
func (m *DBModel) UpdateSection(c clientmodels.Section) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update course_sections set section_name = $1, active = $2, course_id = $3, updated_at = $4 where id = $5`

	_, err := m.DB.ExecContext(ctx, query,
		c.SectionName,
		c.Active,
		c.CourseID,
		time.Now(),
		c.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// InsertSection inserts a course section
func (m *DBModel) InsertSection(c clientmodels.Section) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into course_sections (section_name, course_id, active, created_at, updated_at) 
			values ($1, $2, $3, $4, $5)
			returning id`

	var id int
	err := m.DB.QueryRowContext(ctx, query,
		c.SectionName,
		c.CourseID,
		c.Active,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return id, nil
}

// DeleteSection deletes a course section
func (m *DBModel) DeleteSection(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from course_sections where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// GetSection gets a section
func (m *DBModel) GetSection(id int) (clientmodels.Section, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var s clientmodels.Section

	query := `SELECT s.id, s.section_name, s.active, 
		s.course_id, s.created_at, s.updated_at, c.id as course_id, c.course_name, c.active,
		c.prof_name, c.prof_email, c.teams_link,
		c.created_at as course_created_at, c.updated_at as course_updated_at
		FROM course_sections s
		left join courses c on (s.course_id = c.id)
		where s.id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&s.ID,
		&s.SectionName,
		&s.Active,
		&s.CourseID,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.Course.ID,
		&s.Course.CourseName,
		&s.Course.Active,
		&s.Course.ProfName,
		&s.Course.ProfEmail,
		&s.Course.TeamsLink,
		&s.Course.CreatedAt,
		&s.Course.UpdatedAt,
	)

	if err != nil {
		return s, err
	}

	// get students, if any
	var students []clientmodels.Student
	query2 := `select u.id, u.first_name, u.last_name 
		from users u 
		where u.id in (select user_id from section_students where section_id = $1)`

	rows, err := m.DB.QueryContext(ctx, query2, s.ID)
	if err != nil {
		return s, err
	}
	defer rows.Close()

	for rows.Next() {
		var u clientmodels.Student
		err = rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
		)
		students = append(students, u)
	}
	s.Students = students

	return s, nil

}

// AllCourses returns slice of courses (without lectures)
func (m *DBModel) AllCourses() ([]clientmodels.Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `SELECT id, course_name, active, 
		prof_name, prof_email, teams_link,
		created_at, updated_at FROM courses ORDER BY course_name`

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
			l.updated_at, coalesce(v.video_name, ''), coalesce(v.file_name, ''), coalesce(v.thumb, ''), coalesce(v.duration, 0), l.posted_date at time zone 'America/Halifax'
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
			&l.Video.Duration,
			&l.PostedDate,
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
			l.updated_at, coalesce(v.video_name, ''), coalesce(v.file_name, ''), coalesce(v.thumb, ''), coalesce(v.duration, ''), l.posted_date at time zone 'America/Halifax'
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
			&l.Video.Duration,
			&l.PostedDate,
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
			l.updated_at, coalesce(v.video_name, ''), coalesce(v.file_name, ''), coalesce(v.thumb, ''), 
			coalesce(v.duration, 0), l.posted_date
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
		&l.Video.Duration,
		&l.PostedDate,
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

// InsertLecture inserts a lecture lecture and returns new id
func (m *DBModel) InsertLecture(c clientmodels.Lecture) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	if c.VideoID > 0 {

		query := `insert into lectures (course_id, lecture_name, video_id, active, sort_order, notes, created_at, updated_at, posted_date)
			values ($1, $2, $3, $4, (select max(sort_order) + 1 from lectures where course_id = $1), $5, $6, $7, $8) returning id`

		err := m.DB.QueryRowContext(ctx, query, c.CourseID, c.LectureName, c.VideoID, c.Active, c.Notes, time.Now(), time.Now(), c.PostedDate).Scan(&newID)

		if err != nil {
			fmt.Println("Error inserting new course lecture")
			fmt.Println(err)
			return 0, err
		}
	} else {
		query := `insert into lectures (course_id, lecture_name, active, sort_order, notes, created_at, updated_at, posted_date)
			values ($1, $2, $3, (select max(sort_order) + 1 from lectures where course_id = $1), $4, $5, $6, $7) returning id`

		err := m.DB.QueryRowContext(ctx, query, c.CourseID, c.LectureName, c.Active, c.Notes, time.Now(), time.Now(), c.PostedDate).Scan(&newID)

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
			updated_at = $5, posted_date = $6 where id = $7`

		_, err := m.DB.ExecContext(ctx, query, c.LectureName, c.VideoID, c.Active, c.Notes, time.Now(), c.PostedDate, c.ID)

		if err != nil {
			fmt.Println("Error updating course lecture")
			fmt.Println(err)
			return err
		}
	} else {
		query := `update lectures set lecture_name = $1, video_id = null, active = $2, notes = $3,
			updated_at = $4, posted_date = $5 where id = $6`

		_, err := m.DB.ExecContext(ctx, query, c.LectureName, c.Active, c.Notes, time.Now(), c.PostedDate, c.ID)

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

// GetNextPreviousLectures gets ids for next/previous buttons
func (m *DBModel) GetNextPreviousLectures(courseID, lectureID int) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var next, previous int

	query := `select coalesce(
				(select id from lectures where sort_order > (select sort_order from lectures where id = $1) 
				and course_id = $2 and active = 1 order by sort_order asc limit 1), 0);`

	row := m.DB.QueryRowContext(ctx, query, lectureID, courseID)

	err := row.Scan(&next)
	if err != nil {
		fmt.Println("Error getting course next lecture")
		fmt.Println(err)
		return 0, 0, err
	}

	query = `select coalesce(
				(select id from lectures where sort_order < (select sort_order from lectures where id = $1) 
				and course_id = $2 and active = 1 order by sort_order desc limit 1), 0);`

	row = m.DB.QueryRowContext(ctx, query, lectureID, courseID)

	err = row.Scan(&previous)
	if err != nil {
		fmt.Println("Error getting course previous lecture")
		fmt.Println(err)
		return 0, 0, err
	}

	return next, previous, nil
}

// InsertAssignment inserts an assignment and returns new id
func (m *DBModel) InsertAssignment(c clientmodels.Assignment) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	query := `insert into assignments (file_name_display, file_name, user_id, course_id, description, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err := m.DB.QueryRowContext(ctx, query,
		c.FileNameDisplay,
		c.FileName,
		c.UserID,
		c.CourseID,
		c.Description,
		time.Now(),
		time.Now()).Scan(&newID)

	if err != nil {
		fmt.Println("Error inserting new assignment")
		fmt.Println(err)
		return 0, err
	}

	return newID, nil
}

// UpdateAssignment updates an assignment (grading)
func (m *DBModel) UpdateAssignment(a clientmodels.Assignment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		update assignments set 
			mark = $1, 
			total_value = $2,
			updated_at = $3
		where 
			id = $3`

	_, err := m.DB.ExecContext(ctx, stmt,
		a.Mark,
		a.TotalValue,
		time.Now(),
		a.ID)
	if err != nil {
		return err
	}
	return nil
}

// AllAssignments gets assignments
func (m *DBModel) AllAssignments(id int) ([]clientmodels.Assignment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	where := ""
	if id > 0 {
		where = fmt.Sprintf("where user_id = %d", id)
	}

	var a []clientmodels.Assignment

	stmt := fmt.Sprintf(`SELECT a.id, a.file_name_display, a.file_name, a.user_id, a.course_id, 
		a.mark, a.total_value, a.processed, a.created_at, a.updated_at,
		u.id, u.first_name, u.last_name, u.email,
		c.id, c.course_name, a.description
		FROM 
			assignments a 
			left join users u on (a.user_id = u.id)
			left join courses c on (a.course_id = c.id)
		%s
		ORDER BY updated_at desc`, where)

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s clientmodels.Assignment
		err = rows.Scan(
			&s.ID,
			&s.FileNameDisplay,
			&s.FileName,
			&s.UserID,
			&s.CourseID,
			&s.Mark,
			&s.TotalValue,
			&s.Processed,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.User.ID,
			&s.User.FirstName,
			&s.User.LastName,
			&s.User.Email,
			&s.Course.ID,
			&s.Course.CourseName,
			&s.Description)
		if err != nil {
			return nil, err
		}
		a = append(a, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return a, nil
}

// GetAssignment gets one assignment
func (m *DBModel) GetAssignment(id int) (clientmodels.Assignment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var s clientmodels.Assignment

	stmt := `SELECT a.id, a.file_name_display, a.file_name, a.user_id, a.course_id, 
		a.mark, a.total_value, a.processed, a.created_at, a.updated_at,
		u.id, u.first_name, u.last_name, u.email,
		c.id, c.course_name, a.description
		FROM 
			assignments a 
			left join users u on (a.user_id = u.id)
			left join courses c on (a.course_id = c.id)
		where a.id = $1`

	row := m.DB.QueryRowContext(ctx, stmt, id)

	err := row.Scan(
		&s.ID,
		&s.FileNameDisplay,
		&s.FileName,
		&s.UserID,
		&s.CourseID,
		&s.Mark,
		&s.TotalValue,
		&s.Processed,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.User.ID,
		&s.User.FirstName,
		&s.User.LastName,
		&s.User.Email,
		&s.Course.ID,
		&s.Course.CourseName,
		&s.Description)
	if err != nil {
		return s, err
	}

	return s, nil
}

// GradeAssignment assigns a grade
func (m *DBModel) GradeAssignment(a clientmodels.Assignment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		update assignments set 
			mark = $1, 
			total_value = $2,
			processed = 1,
			updated_at = $3
		where 
			id = $4`

	_, err := m.DB.ExecContext(ctx, stmt,
		a.Mark,
		a.TotalValue,
		time.Now(),
		a.ID)
	if err != nil {
		return err
	}
	return nil
}

// RecordCourseAccess records a student starting/leaving a lecture
func (m *DBModel) RecordCourseAccess(a clientmodels.CourseAccess) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// if the lecture has a video, duration cannot be longer than the video duration
	lecture, _ := m.GetLecture(a.LectureID)
	if lecture.VideoID > 0 {
		// has a video
		if a.Duration > lecture.Video.Duration {
			a.Duration = lecture.Video.Duration
		}
	} else {
		// if the duration > 300 seconds, and the lecture does not have a video,
		// the set duration to 300 seconds
		if a.Duration > 300 {
			a.Duration = 300
		}
	}

	query := `insert into course_accesses (user_id, lecture_id, course_id, duration, created_at,
			updated_at) values ($1, $2, (select course_id from lectures where id = $3), $4, $5, $6)`

	_, err := m.DB.ExecContext(ctx, query,
		a.UserID,
		a.LectureID,
		a.LectureID,
		a.Duration,
		a.CreatedAt,
		a.UpdatedAt,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// CourseAccessHistory gets access history for a course by id
func (m *DBModel) CourseAccessHistory(courseID int) ([]clientmodels.CourseAccess, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var a []clientmodels.CourseAccess

	query := `select ca.id, ca.user_id, ca.lecture_id, ca.course_id, ca.duration, ca.created_at, ca.updated_at,
		u.first_name, u.last_name, l.lecture_name, u.id
		from course_accesses ca 
		left join users u on (ca.user_id = u.id)
		left join lectures l on (ca.lecture_id = l.id)
		where ca.course_id = $1 and ca.user_id <> 1 order by ca.created_at desc`

	rows, err := m.DB.QueryContext(ctx, query, courseID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s clientmodels.CourseAccess
		err = rows.Scan(
			&s.ID,
			&s.UserID,
			&s.LectureID,
			&s.CourseID,
			&s.Duration,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.Student.FirstName,
			&s.Student.FirstName,
			&s.Student.LastName,
			&s.Lecture.LectureName,
			&s.Student.ID,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		a = append(a, s)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return a, nil
}

// CourseAccessHistoryForStudent gets access history for a course by id
func (m *DBModel) CourseAccessHistoryForStudent(userID int) ([]clientmodels.CourseAccess, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var a []clientmodels.CourseAccess

	query := `select ca.id, ca.user_id, ca.lecture_id, ca.course_id, ca.duration, ca.created_at, ca.updated_at,
		u.first_name, u.last_name, l.lecture_name, c.course_name
		from course_accesses ca 
		left join users u on (ca.user_id = u.id)
		left join lectures l on (ca.lecture_id = l.id)
		left join courses c on (ca.course_id = c.id)
		where ca.user_id = $1 order by ca.created_at desc`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s clientmodels.CourseAccess
		err = rows.Scan(
			&s.ID,
			&s.UserID,
			&s.LectureID,
			&s.CourseID,
			&s.Duration,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.Student.FirstName,
			&s.Student.LastName,
			&s.Lecture.LectureName,
			&s.Course.CourseName,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		a = append(a, s)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return a, nil
}

// AllStudents returns all students
func (m *DBModel) AllStudents() ([]clientmodels.Student, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var a []clientmodels.Student

	stmt := `
		SELECT id, last_name, first_name, email, user_active, created_at, updated_at, 
		(select coalesce(sum(duration), 0) from course_accesses where user_id = u.id)
		FROM users u
		where access_level < 3
    	ORDER BY last_name
		`

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s clientmodels.Student
		err = rows.Scan(&s.ID, &s.LastName, &s.FirstName, &s.Email, &s.UserActive, &s.CreatedAt, &s.UpdatedAt, &s.TimeInCourse)
		if err != nil {
			return nil, err
		}

		// get assignments, if any
		assignments, _ := m.AllAssignments(s.ID)
		s.Assignments = assignments

		// Append it to the slice
		a = append(a, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return a, nil
}

// GetTrafficForCourse gets chart data
func (m *DBModel) GetTrafficForCourse(id int) ([]clientmodels.CourseTraffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var a []clientmodels.CourseTraffic

	stmt := `
		select l.lecture_name, 
		coalesce ((select sum(duration) from course_accesses where lecture_id = l.id and user_id > 1), 0) as total_time,
		coalesce ((select count(id) from course_accesses where lecture_id = l.id and user_id > 1), 0) as total_views
		from lectures l
		where l.course_id = $1
		order by sort_order
		`

	rows, err := m.DB.QueryContext(ctx, stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s clientmodels.CourseTraffic
		err = rows.Scan(
			&s.LectureName,
			&s.TotalTime,
			&s.TotalViews)
		if err != nil {
			return nil, err
		}

		// Append it to the slice
		a = append(a, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return a, nil
}

// GetTrafficForCourseForStudent gets chart data
func (m *DBModel) GetTrafficForCourseForStudent(id, studentID int) ([]clientmodels.CourseTraffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var a []clientmodels.CourseTraffic

	stmt := `
		select l.lecture_name, 
		coalesce ((select sum(duration) from course_accesses where lecture_id = l.id and user_id = $1), 0) as total_time,
		coalesce ((select count(id) from course_accesses where lecture_id = l.id and user_id = $1), 0) as total_views
		from lectures l
		where l.course_id = $2
		order by sort_order
		`

	rows, err := m.DB.QueryContext(ctx, stmt, studentID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s clientmodels.CourseTraffic
		err = rows.Scan(
			&s.LectureName,
			&s.TotalTime,
			&s.TotalViews)
		if err != nil {
			return nil, err
		}

		// Append it to the slice
		a = append(a, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return a, nil
}

// StudentsForSection returns slice of all students, and whether registered or not
func (m *DBModel) StudentsForSection(id int) ([]clientmodels.Student, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var students []clientmodels.Student

	query := `select id, first_name, last_name, 
				case when (select count(user_id) from section_students where section_id = $1 and user_id = u.id) = 0 then 0
				else 1 end as is_registered 
				from users u order by last_name`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s clientmodels.Student
		err = rows.Scan(
			&s.ID,
			&s.FirstName,
			&s.LastName,
			&s.IsRegistered,
		)
		if err != nil {
			return nil, err
		}

		students = append(students, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// UpdateEnrollmentForSection manages student enrollment
func (m *DBModel) UpdateEnrollmentForSection(id int, students []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if len(students) == 0 {
		students = append(students, "0")
	}

	// first delete all unchecked students (they will not be in the slice)
	stmt := fmt.Sprintf(`delete from section_students where user_id not in (%s) and section_id = $1`, strings.Join(students, ","))

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		fmt.Println("Error deleting", err)
		return err
	}

	// add in new students
	for _, x := range students {
		stmt = fmt.Sprintf(`insert into section_students (section_id, user_id, created_at, updated_at)
		select $1, %s, $2, $3 where not exists (select id from section_students 
		where user_id = %s and section_id = $4)`, x, x)

		_, err := m.DB.ExecContext(ctx, stmt, id, time.Now(), time.Now(), id)
		if err != nil {
			fmt.Println("Error inserting", err)
			return err
		}
	}

	return nil
}
