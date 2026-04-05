package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"
	_ "modernc.org/sqlite"
)

type DB struct { db *sql.DB }

type Courses struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Subject string `json:"subject"`
	GradeLevel string `json:"grade_level"`
	Description string `json:"description"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type Lessons struct {
	ID string `json:"id"`
	CourseId string `json:"course_id"`
	Title string `json:"title"`
	LessonNumber int64 `json:"lesson_number"`
	Date string `json:"date"`
	DurationMinutes int64 `json:"duration_minutes"`
	Objectives string `json:"objectives"`
	Materials string `json:"materials"`
	Content string `json:"content"`
	Homework string `json:"homework"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil { return nil, err }
	db, err := sql.Open("sqlite", filepath.Join(d, "curriculum.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil { return nil, err }
	db.SetMaxOpenConns(1)
	db.Exec(`CREATE TABLE IF NOT EXISTS courses(id TEXT PRIMARY KEY, title TEXT NOT NULL, subject TEXT DEFAULT '', grade_level TEXT DEFAULT '', description TEXT DEFAULT '', status TEXT DEFAULT '', created_at TEXT DEFAULT(datetime('now')))`)
	db.Exec(`CREATE TABLE IF NOT EXISTS lessons(id TEXT PRIMARY KEY, course_id TEXT NOT NULL, title TEXT NOT NULL, lesson_number INTEGER DEFAULT 0, date TEXT DEFAULT '', duration_minutes INTEGER DEFAULT 0, objectives TEXT DEFAULT '', materials TEXT DEFAULT '', content TEXT DEFAULT '', homework TEXT DEFAULT '', status TEXT DEFAULT '', created_at TEXT DEFAULT(datetime('now')))`)
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }
func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string { return time.Now().UTC().Format(time.RFC3339) }

func (d *DB) CreateCourses(e *Courses) error {
	e.ID = genID(); e.CreatedAt = now()
	_, err := d.db.Exec(`INSERT INTO courses(id, title, subject, grade_level, description, status, created_at) VALUES(?, ?, ?, ?, ?, ?, ?)`, e.ID, e.Title, e.Subject, e.GradeLevel, e.Description, e.Status, e.CreatedAt)
	return err
}

func (d *DB) GetCourses(id string) *Courses {
	var e Courses
	if d.db.QueryRow(`SELECT id, title, subject, grade_level, description, status, created_at FROM courses WHERE id=?`, id).Scan(&e.ID, &e.Title, &e.Subject, &e.GradeLevel, &e.Description, &e.Status, &e.CreatedAt) != nil { return nil }
	return &e
}

func (d *DB) ListCourses() []Courses {
	rows, _ := d.db.Query(`SELECT id, title, subject, grade_level, description, status, created_at FROM courses ORDER BY created_at DESC`)
	if rows == nil { return nil }; defer rows.Close()
	var o []Courses
	for rows.Next() { var e Courses; rows.Scan(&e.ID, &e.Title, &e.Subject, &e.GradeLevel, &e.Description, &e.Status, &e.CreatedAt); o = append(o, e) }
	return o
}

func (d *DB) UpdateCourses(e *Courses) error {
	_, err := d.db.Exec(`UPDATE courses SET title=?, subject=?, grade_level=?, description=?, status=? WHERE id=?`, e.Title, e.Subject, e.GradeLevel, e.Description, e.Status, e.ID)
	return err
}

func (d *DB) DeleteCourses(id string) error {
	_, err := d.db.Exec(`DELETE FROM courses WHERE id=?`, id)
	return err
}

func (d *DB) CountCourses() int {
	var n int; d.db.QueryRow(`SELECT COUNT(*) FROM courses`).Scan(&n); return n
}

func (d *DB) SearchCourses(q string, filters map[string]string) []Courses {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (title LIKE ? OR subject LIKE ? OR grade_level LIKE ? OR description LIKE ?)"
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
	}
	if v, ok := filters["status"]; ok && v != "" { where += " AND status=?"; args = append(args, v) }
	rows, _ := d.db.Query(`SELECT id, title, subject, grade_level, description, status, created_at FROM courses WHERE `+where+` ORDER BY created_at DESC`, args...)
	if rows == nil { return nil }; defer rows.Close()
	var o []Courses
	for rows.Next() { var e Courses; rows.Scan(&e.ID, &e.Title, &e.Subject, &e.GradeLevel, &e.Description, &e.Status, &e.CreatedAt); o = append(o, e) }
	return o
}

func (d *DB) CreateLessons(e *Lessons) error {
	e.ID = genID(); e.CreatedAt = now()
	_, err := d.db.Exec(`INSERT INTO lessons(id, course_id, title, lesson_number, date, duration_minutes, objectives, materials, content, homework, status, created_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, e.ID, e.CourseId, e.Title, e.LessonNumber, e.Date, e.DurationMinutes, e.Objectives, e.Materials, e.Content, e.Homework, e.Status, e.CreatedAt)
	return err
}

func (d *DB) GetLessons(id string) *Lessons {
	var e Lessons
	if d.db.QueryRow(`SELECT id, course_id, title, lesson_number, date, duration_minutes, objectives, materials, content, homework, status, created_at FROM lessons WHERE id=?`, id).Scan(&e.ID, &e.CourseId, &e.Title, &e.LessonNumber, &e.Date, &e.DurationMinutes, &e.Objectives, &e.Materials, &e.Content, &e.Homework, &e.Status, &e.CreatedAt) != nil { return nil }
	return &e
}

func (d *DB) ListLessons() []Lessons {
	rows, _ := d.db.Query(`SELECT id, course_id, title, lesson_number, date, duration_minutes, objectives, materials, content, homework, status, created_at FROM lessons ORDER BY created_at DESC`)
	if rows == nil { return nil }; defer rows.Close()
	var o []Lessons
	for rows.Next() { var e Lessons; rows.Scan(&e.ID, &e.CourseId, &e.Title, &e.LessonNumber, &e.Date, &e.DurationMinutes, &e.Objectives, &e.Materials, &e.Content, &e.Homework, &e.Status, &e.CreatedAt); o = append(o, e) }
	return o
}

func (d *DB) UpdateLessons(e *Lessons) error {
	_, err := d.db.Exec(`UPDATE lessons SET course_id=?, title=?, lesson_number=?, date=?, duration_minutes=?, objectives=?, materials=?, content=?, homework=?, status=? WHERE id=?`, e.CourseId, e.Title, e.LessonNumber, e.Date, e.DurationMinutes, e.Objectives, e.Materials, e.Content, e.Homework, e.Status, e.ID)
	return err
}

func (d *DB) DeleteLessons(id string) error {
	_, err := d.db.Exec(`DELETE FROM lessons WHERE id=?`, id)
	return err
}

func (d *DB) CountLessons() int {
	var n int; d.db.QueryRow(`SELECT COUNT(*) FROM lessons`).Scan(&n); return n
}

func (d *DB) SearchLessons(q string, filters map[string]string) []Lessons {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (course_id LIKE ? OR title LIKE ? OR objectives LIKE ? OR materials LIKE ? OR content LIKE ? OR homework LIKE ?)"
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
	}
	if v, ok := filters["status"]; ok && v != "" { where += " AND status=?"; args = append(args, v) }
	rows, _ := d.db.Query(`SELECT id, course_id, title, lesson_number, date, duration_minutes, objectives, materials, content, homework, status, created_at FROM lessons WHERE `+where+` ORDER BY created_at DESC`, args...)
	if rows == nil { return nil }; defer rows.Close()
	var o []Lessons
	for rows.Next() { var e Lessons; rows.Scan(&e.ID, &e.CourseId, &e.Title, &e.LessonNumber, &e.Date, &e.DurationMinutes, &e.Objectives, &e.Materials, &e.Content, &e.Homework, &e.Status, &e.CreatedAt); o = append(o, e) }
	return o
}
