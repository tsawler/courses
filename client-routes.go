package clienthandlers

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	mw "github.com/tsawler/goblender/pkg/middleware"
	"net/http"
)

// ClientRoutes is used to handle custom routes for specific clients
func ClientRoutes(mux *pat.PatternServeMux, standardMiddleWare, dynamicMiddleware alice.Chain) (*pat.PatternServeMux, error) {

	mux.Get("/courses/all", dynamicMiddleware.Append(mw.Auth).ThenFunc(AllCourses))
	mux.Get("/courses/overview/:ID", dynamicMiddleware.Append(mw.Auth).ThenFunc(ShowCourse))
	mux.Get("/courses/lecture/:ID", dynamicMiddleware.Append(mw.Auth).ThenFunc(ShowLecture))

	// course admin
	mux.Get("/admin/courses/all", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(AdminAllCourses))
	mux.Post("/admin/courses/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(PostAdminCourse))
	mux.Get("/admin/courses/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(AdminCourse))
	mux.Get("/admin/courses/lecture/get-content/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(GetLectureContentJSON))
	mux.Post("/admin/courses/lecture/ajax/savelecture", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(SaveLecture))
	mux.Get("/admin/courses/lecture/:courseID/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(AdminLecture))
	mux.Post("/admin/courses/lecture/:courseID/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(PostAdminLecture))

	// public folder
	fileServer := http.FileServer(http.Dir("./client/clienthandlers/public/"))
	mux.Get("/client/static/", http.StripPrefix("/client/static", fileServer))

	return mux, nil
}
