package clienthandlers

import (
	"github.com/tsawler/goblender/pkg/helpers"
	"github.com/tsawler/goblender/pkg/templates"
	"net/http"
	"strconv"
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

	helpers.Render(w, r, "course.page.tmpl", &templates.TemplateData{
		RowSets: rowSets,
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

	rowSets := make(map[string]interface{})
	rowSets["lecture"] = lecture

	helpers.Render(w, r, "lecture.page.tmpl", &templates.TemplateData{
		RowSets: rowSets,
	})
}
