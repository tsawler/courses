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

	// public folder
	fileServer := http.FileServer(http.Dir("./client/clienthandlers/public/"))
	mux.Get("/client/static/", http.StripPrefix("/client/static", fileServer))

	return mux, nil
}
