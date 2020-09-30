package clienthandlers

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"github.com/tsawler/goblender/pkg/handlers"
	mw "github.com/tsawler/goblender/pkg/middleware"
	"net/http"
)

// ClientRoutes is used to handle custom routes for specific clients
func ClientRoutes(mux *pat.PatternServeMux, standardMiddleWare, dynamicMiddleware alice.Chain) (*pat.PatternServeMux, error) {

	// assignments
	mux.Get("/admin/users/profile", dynamicMiddleware.Append(mw.Auth).ThenFunc(StudentProfile))
	mux.Get("/admin/assignments-for-student", dynamicMiddleware.Append(mw.Auth).ThenFunc(StudentAssignments))
	mux.Get("/courses/assignments/submit-an-assignment", dynamicMiddleware.Append(mw.Auth).ThenFunc(SubmitAssignment))
	mux.Post("/courses/assignments/submit-an-assignment", dynamicMiddleware.Append(mw.Auth).ThenFunc(PostSubmitAssignment))

	// courses & lectures
	mux.Get("/courses/all", dynamicMiddleware.Append(mw.Auth).ThenFunc(AllCourses))
	mux.Get("/courses/overview/:ID", dynamicMiddleware.Append(mw.Auth).ThenFunc(ShowCourse))
	mux.Get("/courses/lecture/:ID", dynamicMiddleware.Append(mw.Auth).ThenFunc(ShowLecture))

	// record access
	mux.Post("/courses/lecture/log/record-leaving", dynamicMiddleware.Append(mw.Auth).ThenFunc(StudentLeftLecture))

	// assignments (admin)
	mux.Get("/admin/courses/assignments", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(Assignments))
	mux.Get("/admin/courses/assignment/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(Assignment))
	mux.Post("/admin/courses/assignment/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(GradeAssignment))

	// traffic
	mux.Get("/admin/courses/traffic", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(CourseTraffic))
	mux.Get("/admin/courses/course/ajax/traffic-data", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(CourseTrafficData))
	mux.Get("/admin/courses/course/ajax/traffic-data-for-student", dynamicMiddleware.Append(mw.Auth).ThenFunc(CourseTrafficDataForStudent))
	mux.Get("/admin/courses/course/ajax/traffic-data-for-student-admin", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(CourseTrafficDataForStudentAdmin))

	// section admin
	mux.Get("/admin/sections/all", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(AdminAllSections))
	mux.Get("/admin/sections/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(AdminSection))
	mux.Post("/admin/sections/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(PostAdminSection))
	mux.Get("/admin/sections/delete/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(DeleteSection))
	mux.Get("/admin/sections/remove-student/:SectionID/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(UnenrolStudent))
	mux.Get("/admin/sections/students/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(SectionStudents))
	mux.Post("/admin/sections/students/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(PostSectionStudents))

	// course admin
	mux.Get("/admin/courses/all", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(AdminAllCourses))
	mux.Post("/admin/courses/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(PostAdminCourse))
	mux.Get("/admin/courses/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(AdminCourse))
	mux.Get("/admin/courses/course/get-content/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(GetCourseContentJSON))
	mux.Post("/admin/courses/course/ajax/savecourse", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(SaveCourse))

	mux.Post("/admin/courses/ajax/save-lecture-sort-order", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(SaveLectureSortOrder))
	mux.Get("/admin/courses/lecture/get-content/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(GetLectureContentJSON))
	mux.Post("/admin/courses/lecture/ajax/savelecture", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(SaveLecture))
	mux.Get("/admin/courses/lecture/:courseID/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(AdminLecture))
	mux.Post("/admin/courses/lecture/:courseID/:ID", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(PostAdminLecture))

	mux.Get("/admin/members/all", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(MembersAll))
	mux.Get("/admin/members/all-deleted", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(handlers.Repo.DeletedMembersAll(app)))
	mux.Get("/admin/members/:id", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(MemberEdit))

	mux.Get("/admin/courses/:ID/accesses", dynamicMiddleware.Append(mw.Auth).Append(mw.PagesRole).ThenFunc(CourseAccessHistory))

	// public folder
	fileServer := http.FileServer(http.Dir("./client/clienthandlers/public/"))
	mux.Get("/client/static/", http.StripPrefix("/client/static", fileServer))

	return mux, nil
}
