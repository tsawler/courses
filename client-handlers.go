package clienthandlers

import (
	"encoding/json"
	"fmt"
	"github.com/tsawler/goblender/client/clienthandlers/clientmodels"
	"github.com/tsawler/goblender/pkg/forms"
	"github.com/tsawler/goblender/pkg/handlers"
	"github.com/tsawler/goblender/pkg/helpers"
	"github.com/tsawler/goblender/pkg/models"
	"github.com/tsawler/goblender/pkg/templates"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

// AllCourses lists all active courses with link to overview
func AllCourses(w http.ResponseWriter, r *http.Request) {

	pg, err := repo.DB.GetPageBySlug("courses")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	courses, err := dbModel.AllActiveCourses()
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	rowSets := make(map[string]interface{})
	rowSets["courses"] = courses

	helpers.Render(w, r, "courses.page.tmpl", &templates.TemplateData{
		Page:    pg,
		RowSets: rowSets,
	})
}

// ShowCourse shows one course
func ShowCourse(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.Atoi(r.URL.Query().Get(":ID"))
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	course, err := dbModel.GetCourse(courseID)
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	rowSets := make(map[string]interface{})
	rowSets["course"] = course

	pg := models.Page{
		ID:          course.ID,
		AccessLevel: 2,
		Active:      1,
		SEOImage:    0,
		MenuId:      0,
		MenuColor:   "navbar-light",
		HasSlider:   0,
		Immutable:   1,
		Content:     course.Description,
		PageTitle:   course.CourseName,
	}

	helpers.Render(w, r, "course.page.tmpl", &templates.TemplateData{
		RowSets: rowSets,
		Page:    pg,
	})
}

// ShowLecture shows one lecture
func ShowLecture(w http.ResponseWriter, r *http.Request) {
	lectureID, err := strconv.Atoi(r.URL.Query().Get(":ID"))
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	lecture, err := dbModel.GetLecture(lectureID)
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	course, err := dbModel.GetCourse(lecture.CourseID)
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	rowSets := make(map[string]interface{})
	rowSets["lecture"] = lecture
	rowSets["course"] = course

	pg := models.Page{
		ID:          lecture.ID,
		AccessLevel: 2,
		Active:      1,
		SEOImage:    0,
		MenuId:      0,
		MenuColor:   "navbar-light",
		HasSlider:   0,
		Immutable:   1,
		Content:     lecture.Notes,
		PageTitle:   lecture.LectureName,
	}

	next, prev, err := dbModel.GetNextPreviousLectures(course.ID, lecture.ID)
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	intMap := make(map[string]int)
	intMap["next"] = next
	intMap["previous"] = prev

	helpers.Render(w, r, "lecture.page.tmpl", &templates.TemplateData{
		RowSets: rowSets,
		Page:    pg,
		IntMap:  intMap,
	})
}

// AdminAllCourses shows list of all courses for admin
func AdminAllCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := dbModel.AllCourses()
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	rowSets := make(map[string]interface{})
	rowSets["courses"] = courses

	helpers.Render(w, r, "courses-all-admin.page.tmpl", &templates.TemplateData{
		RowSets: rowSets,
	})
}

// AdminCourse shows course for add/edit
func AdminCourse(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.Atoi(r.URL.Query().Get(":ID"))
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	var course clientmodels.Course

	if courseID > 0 {
		c, err := dbModel.GetCourse(courseID)
		if err != nil {
			errorLog.Println(err)
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}

		course = c
	}

	rowSets := make(map[string]interface{})
	rowSets["course"] = course

	helpers.Render(w, r, "courses-admin.page.tmpl", &templates.TemplateData{
		RowSets: rowSets,
		Form:    forms.New(nil),
	})
}

// SortOrder is type for sorting
type SortOrder struct {
	ID    string `json:"id"`
	Order int    `json:"order"`
}

// PostAdminCourse updates or adds a course
func PostAdminCourse(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.Atoi(r.URL.Query().Get(":ID"))
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	var course clientmodels.Course
	if courseID > 0 {
		c, err := dbModel.GetCourse(courseID)
		if err != nil {
			errorLog.Println(err)
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}
		course = c
		course.CourseName = r.Form.Get("course_name")
		active, _ := strconv.Atoi(r.Form.Get("active"))
		course.Active = active
		course.ProfName = r.Form.Get("prof_name")
		course.ProfEmail = r.Form.Get("prof_email")
		course.TeamsLink = r.Form.Get("teams_link")
		err = dbModel.UpdateCourse(course)
		if err != nil {
			errorLog.Println(err)
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}
	} else {
		// inserting course
		course.CourseName = r.Form.Get("course_name")
		active, _ := strconv.Atoi(r.Form.Get("active"))
		course.Active = active
		course.Description = ""
		newID, err := dbModel.InsertCourse(course)
		if err != nil {
			errorLog.Println(err)
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}
		course.ID = newID
	}

	//now do sort order
	var sorted []SortOrder
	sortList := r.Form.Get("sort_list")

	err = json.Unmarshal([]byte(sortList), &sorted)
	if err != nil {
		app.ErrorLog.Println(err)
	}

	for _, v := range sorted {
		lectureID, _ := strconv.Atoi(v.ID)
		err := dbModel.UpdateLectureSortOrder(lectureID, v.Order)
		if err != nil {
			app.ErrorLog.Println(err)
		}
	}

	action, _ := strconv.Atoi(r.Form.Get("action"))
	session.Put(r.Context(), "flash", "Changes saved")

	if action == 1 {
		http.Redirect(w, r, "/admin/courses/all", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/admin/courses/%d", course.ID), http.StatusSeeOther)

}

// AdminLecture shows form for lecture
func AdminLecture(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.Atoi(r.URL.Query().Get(":courseID"))
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	lectureID, err := strconv.Atoi(r.URL.Query().Get(":ID"))
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	course, err := dbModel.GetCourse(courseID)
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	var lecture clientmodels.Lecture
	if lectureID > 0 {
		lecture, err = dbModel.GetLecture(lectureID)
		if err != nil {
			errorLog.Println(err)
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}
	} else {
		lecture.CourseID = courseID
		lecture.PostedDate = time.Now()
	}

	rowSets := make(map[string]interface{})
	rowSets["course"] = course
	rowSets["lecture"] = lecture

	helpers.Render(w, r, "lecture-admin.page.tmpl", &templates.TemplateData{
		RowSets: rowSets,
		Form:    forms.New(nil),
	})

}

// PostAdminLecture posts a lecture
func PostAdminLecture(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.Atoi(r.URL.Query().Get(":courseID"))
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	lectureID, err := strconv.Atoi(r.URL.Query().Get(":ID"))
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	videoID, _ := strconv.Atoi(r.Form.Get("video_id"))
	active, _ := strconv.Atoi(r.Form.Get("active"))
	lectureName := r.Form.Get("lecture_name")
	pd := r.Form.Get("posted_date")
	layout := "2006-01-02 15:04"
	t, err := time.Parse(layout, pd)
	if err != nil {
		fmt.Println(err)
	}

	var lecture clientmodels.Lecture
	if lectureID > 0 {
		lecture, err = dbModel.GetLecture(lectureID)
		if err != nil {
			errorLog.Print(err)
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}
	}

	lecture.CourseID = courseID
	lecture.LectureName = lectureName
	lecture.Active = active
	lecture.VideoID = videoID
	lecture.PostedDate = t

	if lectureID == 0 {
		lecture.Notes = ""
		_, err := dbModel.InsertLecture(lecture)
		if err != nil {
			errorLog.Print(err)
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}

	} else {
		err := dbModel.UpdateLecture(lecture)
		if err != nil {
			errorLog.Print(err)
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}
	}
	session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/courses/%d", courseID), http.StatusSeeOther)
}

// GetLectureContentJSON gets html (notes) for lecture on edit page
func GetLectureContentJSON(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":ID"))
	lecture, err := dbModel.GetLecture(id)
	if err != nil {
		errorLog.Println(err)
		return
	}

	theData := handlers.PageContentJSON{
		OK:      true,
		Content: template.HTML(lecture.Notes),
	}

	out, err := json.MarshalIndent(theData, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		app.ErrorLog.Println(err)
	}
}

// SaveLecture saves lecture html (notes)
func SaveLecture(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(r.Form.Get("page_id"))
	pageContent := r.Form.Get("thedata")

	lecture, err := dbModel.GetLecture(id)
	if err != nil {
		app.ErrorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	lecture.Notes = pageContent

	err = dbModel.UpdateLectureContent(lecture)
	if err != nil {
		app.ErrorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	app.Session.Put(r.Context(), "flash", "Lecture successfully updated!")
	http.Redirect(w, r, fmt.Sprintf("/courses/lecture/%d", id), http.StatusSeeOther)
}

// GetCourseContentJSON gets html (description) for course on edit page
func GetCourseContentJSON(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":ID"))
	course, err := dbModel.GetCourse(id)
	if err != nil {
		errorLog.Println(err)
		return
	}

	theData := handlers.PageContentJSON{
		OK:      true,
		Content: template.HTML(course.Description),
	}

	out, err := json.MarshalIndent(theData, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		app.ErrorLog.Println(err)
	}
}

// SaveCourse saves course html (description)
func SaveCourse(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(r.Form.Get("page_id"))
	pageContent := r.Form.Get("thedata")

	course, err := dbModel.GetCourse(id)
	if err != nil {
		app.ErrorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	course.Description = pageContent

	err = dbModel.UpdateCourseContent(course)
	if err != nil {
		app.ErrorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	app.Session.Put(r.Context(), "flash", "Lecture successfully updated!")
	http.Redirect(w, r, fmt.Sprintf("/courses/overview/%d", id), http.StatusSeeOther)
}

// SubmitAssignment displays page to submit an assignment
func SubmitAssignment(w http.ResponseWriter, r *http.Request) {
	pg, err := repo.DB.GetPageBySlug("submit-assignment")
	if err == models.ErrNoRecord {
		helpers.NotFound(w)
		return
	} else if err != nil {
		helpers.ServerError(w, err)
		return
	}

	courses, err := dbModel.AllActiveCourses()
	if err != nil {
		errorLog.Println(err)
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	rowSets := make(map[string]interface{})
	rowSets["courses"] = courses
	helpers.Render(w, r, "submit-assignment.page.tmpl", &templates.TemplateData{
		Page:    pg,
		Form:    forms.New(nil),
		RowSets: rowSets,
	})
}

// PostSubmitAssignment handles assignment submission
func PostSubmitAssignment(w http.ResponseWriter, r *http.Request) {
	userID := app.Session.GetInt(r.Context(), "userID")

	courseID, _ := strconv.Atoi(r.Form.Get("course_id"))
	description := r.Form.Get("description")

	helpers.CreateDirIfNotExist("./ui/static/site-content/assignments/")
	helpers.CreateDirIfNotExist(fmt.Sprintf("./ui/static/site-content/assignments/%d", userID))

	fileName, displayName, err := helpers.UploadOneFileReturnSlugName(r, fmt.Sprintf("./ui/static/site-content/assignments/%d/", userID))
	if err != nil {
		errorLog.Println(err)
	}

	assignment := clientmodels.Assignment{
		FileNameDisplay: displayName,
		FileName:        fileName,
		Description:     description,
		UserID:          userID,
		CourseID:        courseID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	_, err = dbModel.InsertAssignment(assignment)
	if err != nil {
		session.Put(r.Context(), "error", "Error: assignment NOT received!")
		http.Redirect(w, r, "/courses/assignments/submit-an-assignment", http.StatusNotAcceptable)
		return
	}

	session.Put(r.Context(), "flash", "Assignment received!")
	http.Redirect(w, r, "/courses/assignments/submit-an-assignment", http.StatusSeeOther)
}

// Assignments displays assignments in admin tool
func Assignments(w http.ResponseWriter, r *http.Request) {
	a, err := dbModel.AllAssignments(0)
	if err != nil {
		errorLog.Print(err)
	}

	rowSets := make(map[string]interface{})
	rowSets["assignments"] = a

	helpers.Render(w, r, "assignments-admin.page.tmpl", &templates.TemplateData{
		RowSets: rowSets,
	})
}

// Assignment displays assignment in admin tool
func Assignment(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":ID"))
	a, err := dbModel.GetAssignment(id)
	if err != nil {
		errorLog.Print(err)
	}

	rowSets := make(map[string]interface{})
	rowSets["assignment"] = a

	helpers.Render(w, r, "assignment-admin.page.tmpl", &templates.TemplateData{
		RowSets: rowSets,
		Form:    forms.New(nil),
	})
}

// GradeAssignment grades an assignment
func GradeAssignment(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":ID"))
	a, err := dbModel.GetAssignment(id)
	if err != nil {
		errorLog.Print(err)
	}

	a.Mark, _ = strconv.Atoi(r.Form.Get("mark"))
	a.TotalValue, _ = strconv.Atoi(r.Form.Get("total_value"))

	err = dbModel.GradeAssignment(a)

	app.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, "/admin/courses/assignments", http.StatusSeeOther)
}

// StudentAssignments displays assignments in admin tool for a given student
func StudentAssignments(w http.ResponseWriter, r *http.Request) {
	userID := app.Session.GetInt(r.Context(), "userID")
	a, err := dbModel.AllAssignments(userID)
	if err != nil {
		errorLog.Print(err)
	}

	rowSets := make(map[string]interface{})
	rowSets["assignments"] = a

	helpers.Render(w, r, "student-assignments.page.tmpl", &templates.TemplateData{
		RowSets: rowSets,
	})
}

// StudentStartedLecture records student starting lecture
func StudentStartedLecture(w http.ResponseWriter, r *http.Request) {
	lectureID, _ := strconv.Atoi(r.Form.Get("lecture_id"))
	userID := app.Session.GetInt(r.Context(), "userID")

	access := clientmodels.CourseAccess{
		UserID:    userID,
		LectureID: lectureID,
		IsEntered: 1,
		IsLeft:    0,
		Duration:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_ = dbModel.RecordStartLecture(access)
}

// StudentLeftLecture records student leaving lecture
func StudentLeftLecture(w http.ResponseWriter, r *http.Request) {
	lectureID, _ := strconv.Atoi(r.Form.Get("lecture_id"))
	duration, _ := strconv.Atoi(r.Form.Get("duration"))
	userID := app.Session.GetInt(r.Context(), "userID")

	access := clientmodels.CourseAccess{
		UserID:    userID,
		LectureID: lectureID,
		IsEntered: 0,
		IsLeft:    1,
		Duration:  duration,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_ = dbModel.RecordLeaveLecture(access)
}
