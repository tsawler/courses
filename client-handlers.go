package clienthandlers

import (
	"fmt"
	"github.com/tsawler/goblender/client/clienthandlers/clientmodels"
	"github.com/tsawler/goblender/pkg/forms"
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
		err = dbModel.UpdateCourse(course)
		if err != nil {
			errorLog.Println(err)
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}
	} else {
		// inserting course
	}

	session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/courses/%d", course.ID), http.StatusSeeOther)
}
