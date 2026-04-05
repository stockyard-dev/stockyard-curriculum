package server

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/stockyard-dev/stockyard-curriculum/internal/store"
)

type Server struct {
	db     *store.DB
	mux    *http.ServeMux
	limits Limits
}

func New(db *store.DB, limits Limits) *Server {
	s := &Server{db: db, mux: http.NewServeMux(), limits: limits}
	s.mux.HandleFunc("GET /api/courses", s.listCourses)
	s.mux.HandleFunc("POST /api/courses", s.createCourses)
	s.mux.HandleFunc("GET /api/courses/export.csv", s.exportCourses)
	s.mux.HandleFunc("GET /api/courses/{id}", s.getCourses)
	s.mux.HandleFunc("PUT /api/courses/{id}", s.updateCourses)
	s.mux.HandleFunc("DELETE /api/courses/{id}", s.delCourses)
	s.mux.HandleFunc("GET /api/lessons", s.listLessons)
	s.mux.HandleFunc("POST /api/lessons", s.createLessons)
	s.mux.HandleFunc("GET /api/lessons/export.csv", s.exportLessons)
	s.mux.HandleFunc("GET /api/lessons/{id}", s.getLessons)
	s.mux.HandleFunc("PUT /api/lessons/{id}", s.updateLessons)
	s.mux.HandleFunc("DELETE /api/lessons/{id}", s.delLessons)
	s.mux.HandleFunc("GET /api/stats", s.stats)
	s.mux.HandleFunc("GET /api/health", s.health)
	s.mux.HandleFunc("GET /health", s.health)
	s.mux.HandleFunc("GET /ui", s.dashboard)
	s.mux.HandleFunc("GET /ui/", s.dashboard)
	s.mux.HandleFunc("GET /", s.root)
	s.mux.HandleFunc("GET /api/tier", func(w http.ResponseWriter, r *http.Request) {
		wj(w, 200, map[string]any{"tier": s.limits.Tier, "upgrade_url": "https://stockyard.dev/curriculum/"})})
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.mux.ServeHTTP(w, r) }
func wj(w http.ResponseWriter, c int, v any) { w.Header().Set("Content-Type", "application/json"); w.WriteHeader(c); json.NewEncoder(w).Encode(v) }
func we(w http.ResponseWriter, c int, m string) { wj(w, c, map[string]string{"error": m}) }
func (s *Server) root(w http.ResponseWriter, r *http.Request) { if r.URL.Path != "/" { http.NotFound(w, r); return }; http.Redirect(w, r, "/ui", 302) }
func oe[T any](s []T) []T { if s == nil { return []T{} }; return s }
func init() { log.SetFlags(log.LstdFlags | log.Lshortfile) }

func (s *Server) listCourses(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	filters := map[string]string{}
	if v := r.URL.Query().Get("status"); v != "" { filters["status"] = v }
	if q != "" || len(filters) > 0 { wj(w, 200, map[string]any{"courses": oe(s.db.SearchCourses(q, filters))}); return }
	wj(w, 200, map[string]any{"courses": oe(s.db.ListCourses())})
}

func (s *Server) createCourses(w http.ResponseWriter, r *http.Request) {
	if s.limits.MaxItems > 0 { if s.db.CountCourses() >= s.limits.MaxItems { we(w, 402, "Free tier limit reached. Upgrade at https://stockyard.dev/curriculum/"); return } }
	var e store.Courses
	json.NewDecoder(r.Body).Decode(&e)
	if e.Title == "" { we(w, 400, "title required"); return }
	s.db.CreateCourses(&e)
	wj(w, 201, s.db.GetCourses(e.ID))
}

func (s *Server) getCourses(w http.ResponseWriter, r *http.Request) {
	e := s.db.GetCourses(r.PathValue("id"))
	if e == nil { we(w, 404, "not found"); return }
	wj(w, 200, e)
}

func (s *Server) updateCourses(w http.ResponseWriter, r *http.Request) {
	existing := s.db.GetCourses(r.PathValue("id"))
	if existing == nil { we(w, 404, "not found"); return }
	var patch store.Courses
	json.NewDecoder(r.Body).Decode(&patch)
	patch.ID = existing.ID; patch.CreatedAt = existing.CreatedAt
	if patch.Title == "" { patch.Title = existing.Title }
	if patch.Subject == "" { patch.Subject = existing.Subject }
	if patch.GradeLevel == "" { patch.GradeLevel = existing.GradeLevel }
	if patch.Description == "" { patch.Description = existing.Description }
	if patch.Status == "" { patch.Status = existing.Status }
	s.db.UpdateCourses(&patch)
	wj(w, 200, s.db.GetCourses(patch.ID))
}

func (s *Server) delCourses(w http.ResponseWriter, r *http.Request) {
	s.db.DeleteCourses(r.PathValue("id"))
	wj(w, 200, map[string]string{"deleted": "ok"})
}

func (s *Server) exportCourses(w http.ResponseWriter, r *http.Request) {
	items := s.db.ListCourses()
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=courses.csv")
	cw := csv.NewWriter(w)
	cw.Write([]string{"id", "title", "subject", "grade_level", "description", "status", "created_at"})
	for _, e := range items {
		cw.Write([]string{e.ID, fmt.Sprintf("%v", e.Title), fmt.Sprintf("%v", e.Subject), fmt.Sprintf("%v", e.GradeLevel), fmt.Sprintf("%v", e.Description), fmt.Sprintf("%v", e.Status), e.CreatedAt})
	}
	cw.Flush()
}

func (s *Server) listLessons(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	filters := map[string]string{}
	if v := r.URL.Query().Get("status"); v != "" { filters["status"] = v }
	if q != "" || len(filters) > 0 { wj(w, 200, map[string]any{"lessons": oe(s.db.SearchLessons(q, filters))}); return }
	wj(w, 200, map[string]any{"lessons": oe(s.db.ListLessons())})
}

func (s *Server) createLessons(w http.ResponseWriter, r *http.Request) {
	var e store.Lessons
	json.NewDecoder(r.Body).Decode(&e)
	if e.CourseId == "" { we(w, 400, "course_id required"); return }
	if e.Title == "" { we(w, 400, "title required"); return }
	s.db.CreateLessons(&e)
	wj(w, 201, s.db.GetLessons(e.ID))
}

func (s *Server) getLessons(w http.ResponseWriter, r *http.Request) {
	e := s.db.GetLessons(r.PathValue("id"))
	if e == nil { we(w, 404, "not found"); return }
	wj(w, 200, e)
}

func (s *Server) updateLessons(w http.ResponseWriter, r *http.Request) {
	existing := s.db.GetLessons(r.PathValue("id"))
	if existing == nil { we(w, 404, "not found"); return }
	var patch store.Lessons
	json.NewDecoder(r.Body).Decode(&patch)
	patch.ID = existing.ID; patch.CreatedAt = existing.CreatedAt
	if patch.CourseId == "" { patch.CourseId = existing.CourseId }
	if patch.Title == "" { patch.Title = existing.Title }
	if patch.Date == "" { patch.Date = existing.Date }
	if patch.Objectives == "" { patch.Objectives = existing.Objectives }
	if patch.Materials == "" { patch.Materials = existing.Materials }
	if patch.Content == "" { patch.Content = existing.Content }
	if patch.Homework == "" { patch.Homework = existing.Homework }
	if patch.Status == "" { patch.Status = existing.Status }
	s.db.UpdateLessons(&patch)
	wj(w, 200, s.db.GetLessons(patch.ID))
}

func (s *Server) delLessons(w http.ResponseWriter, r *http.Request) {
	s.db.DeleteLessons(r.PathValue("id"))
	wj(w, 200, map[string]string{"deleted": "ok"})
}

func (s *Server) exportLessons(w http.ResponseWriter, r *http.Request) {
	items := s.db.ListLessons()
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=lessons.csv")
	cw := csv.NewWriter(w)
	cw.Write([]string{"id", "course_id", "title", "lesson_number", "date", "duration_minutes", "objectives", "materials", "content", "homework", "status", "created_at"})
	for _, e := range items {
		cw.Write([]string{e.ID, fmt.Sprintf("%v", e.CourseId), fmt.Sprintf("%v", e.Title), fmt.Sprintf("%v", e.LessonNumber), fmt.Sprintf("%v", e.Date), fmt.Sprintf("%v", e.DurationMinutes), fmt.Sprintf("%v", e.Objectives), fmt.Sprintf("%v", e.Materials), fmt.Sprintf("%v", e.Content), fmt.Sprintf("%v", e.Homework), fmt.Sprintf("%v", e.Status), e.CreatedAt})
	}
	cw.Flush()
}

func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	m := map[string]any{}
	m["courses_total"] = s.db.CountCourses()
	m["lessons_total"] = s.db.CountLessons()
	wj(w, 200, m)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	m := map[string]any{"status": "ok", "service": "curriculum"}
	m["courses"] = s.db.CountCourses()
	m["lessons"] = s.db.CountLessons()
	wj(w, 200, m)
}
