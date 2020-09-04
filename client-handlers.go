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

	rowSets := make(map[string]interface{})
	rowSets["lecture"] = lecture

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

	helpers.Render(w, r, "lecture.page.tmpl", &templates.TemplateData{
		RowSets: rowSets,
		Page:    pg,
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
	} else {
		lecture.CourseID = courseID
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
