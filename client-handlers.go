package clienthandlers

import (
	"github.com/tsawler/goblender/pkg/helpers"
	"github.com/tsawler/goblender/pkg/templates"
	"net/http"
	"strconv"
)

// SomeHandler is an example handler
func SomeHandler(w http.ResponseWriter, r *http.Request) {
	helpers.Render(w, r, "client-sample.page.tmpl", &templates.TemplateData{})
}

// CustomShowHome is a sample handler which returns the home page using our local page template for the client,
// and is called from client-routes.go using a route that overrides the one in goBlender. This allows us
// to build custom functionality without having to use non-standard routes.
func CustomShowHome(w http.ResponseWriter, r *http.Request) {
	// do something interesting here, and then render the template
	helpers.Render(w, r, "client-sample.page.tmpl", &templates.TemplateData{})
}

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

	for _, x := range course.Lectures {
		w.Write([]byte(x.LectureName))
	}
}
