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
	mux.Get("/admin/sections/assignments/download-graded-for-student/:ID", dynamicMiddleware.Append(mw.Auth).ThenFunc(DownloadGradeAssignmentForStudent))

	// courses & lectures
	mux.Get("/courses/all", dynamicMiddleware.Append(mw.Auth).ThenFunc(AllCourses))
	mux.Get("/courses/overview/:ID", dynamicMiddleware.Append(mw.Auth).ThenFunc(ShowCourse))
	mux.Get("/courses/lecture/:SectionID/:ID", dynamicMiddleware.Append(mw.Auth).ThenFunc(ShowLecture))

	// record access
	mux.Post("/courses/lecture/log/record-leaving", dynamicMiddleware.Append(mw.Auth).ThenFunc(StudentLeftLecture))

	// assignments (admin)
	mux.Get("/admin/assignments/assignments", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(Assignments))
	mux.Get("/admin/assignments/assignment/:ID", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(Assignment))
	mux.Post("/admin/assignments/assignment/:ID", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(GradeAssignment))
	mux.Get("/admin/assignments/assignments/download-graded/:UserID/:ID", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(DownloadGradeAssignment))

	// traffic
	mux.Get("/admin/sections/traffic", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(CourseTraffic))
	mux.Get("/admin/courses/course/ajax/traffic-data", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(CourseTrafficData))
	mux.Get("/admin/courses/course/ajax/traffic-data-for-student", dynamicMiddleware.Append(mw.Auth).ThenFunc(CourseTrafficDataForStudent))
	mux.Get("/admin/courses/course/ajax/traffic-data-for-student-admin", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(CourseTrafficDataForStudentAdmin))

	// section admin
	mux.Get("/admin/sections/all", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(AdminAllSections))
	mux.Get("/admin/sections/:ID", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(AdminSection))
	mux.Post("/admin/sections/:ID", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(PostAdminSection))
	mux.Get("/admin/sections/delete/:ID", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(DeleteSection))
	mux.Get("/admin/sections/remove-student/:SectionID/:ID", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(UnenrolStudent))
	mux.Get("/admin/sections/students/:ID", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(SectionStudents))
	mux.Post("/admin/sections/students/:ID", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(PostSectionStudents))

	// course admin
	mux.Get("/admin/courses/all", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(AdminAllCourses))
	mux.Post("/admin/courses/:ID", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(PostAdminCourse))
	mux.Get("/admin/courses/:ID", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(AdminCourse))
	mux.Get("/admin/courses/course/get-content/:ID", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(GetCourseContentJSON))
	mux.Post("/admin/courses/course/ajax/savecourse", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(SaveCourse))

	mux.Post("/admin/courses/ajax/save-lecture-sort-order", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(SaveLectureSortOrder))
	mux.Get("/admin/courses/lecture/get-content/:ID", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(GetLectureContentJSON))
	mux.Post("/admin/courses/lecture/ajax/savelecture", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(SaveLecture))
	mux.Get("/admin/courses/lecture/:courseID/:ID", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(AdminLecture))
	mux.Post("/admin/courses/lecture/:courseID/:ID", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(PostAdminLecture))

	mux.Get("/admin/members/all", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(MembersAll))
	mux.Get("/admin/members/all-deleted", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(handlers.Repo.DeletedMembersAll()))
	mux.Get("/admin/members/:id", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(MemberEdit))

	mux.Get("/admin/sections/:ID/accesses", dynamicMiddleware.Append(mw.Auth).Append(CoursesRole).ThenFunc(CourseSectionAccessHistory))

	// public folder
	fileServer := http.FileServer(http.Dir("./client/clienthandlers/public/"))
	mux.Get("/client/static/", http.StripPrefix("/client/static", fileServer))

	return mux, nil
}
