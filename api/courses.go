package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var courseQuery = `
	SELECT semester, course_code, course_name, course_topic, classification, code_available, languages, summary, 
	STRING_AGG(CASE WHEN type = 'class' THEN icon_html ELSE NULL end, ' ') AS classification_icons, 
	STRING_AGG(CASE WHEN type = 'lang' THEN icon_html ELSE NULL END, ' ') AS lang_icons, 
	STRING_AGG(CONCAT('course--', abbreviation), ' ') AS course_classes
	FROM courses 
	INNER JOIN abbreviations ON (courses.classification LIKE CONCAT('%', abbreviations.abbreviation, '%') OR courses.languages LIKE CONCAT('%', abbreviations.abbreviation, '%'))
	GROUP BY course_code, semester, course_name, course_topic, classification, code_available, languages, summary, semester_id
	ORDER BY semester_id DESC
`

type Course struct {
	Semester            string `json:"semester"`
	CourseCode          string `json:"course_code"`
	CourseName          string `json:"course_name"`
	CourseTopic         string `json:"course_topic"`
	Classification      string `json:"classification"`
	CodeAvailable       bool   `json:"code_available"`
	Languages           string `json:"languages"`
	Summary             string `json:"summary"`
	ClassificationIcons string `json:"classification_icons"`
	LangIcons           string `json:"lang_icons"`
	CourseClasses       string `json:"course_classes"`
}

type Abbreviation struct {
	Abbreviation string `json:"abbreviation"`
	IconHTML     string `json:"icon_html"`
	Longname     string `json:"longname"`
}

type Response struct {
	QueryTime       string         `json:"query_time"`
	FreshQuery      bool           `json:"fresh_query"`
	Courses         []Course       `json:"courses"`
	Classifications []Abbreviation `json:"classifications"`
	Languages       []Abbreviation `json:"languages"`
}

func Courses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := GetDB()
	if err != nil {
		log.Println("Error connecting to database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	runQuery := r.URL.Query().Get("secret") == refreshSecret
	var response Response

	// Check if we can use the cached data
	if !runQuery && isCacheValid() {
		cacheData, err := os.ReadFile(cacheFilename)
		if err != nil {
			log.Println("Error reading cache file:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(cacheData, &response); err != nil {
			log.Println("Error unmarshaling cache data:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		response.QueryTime = time.Now().Format(time.RFC3339)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Fetch fresh data from the database
	response.FreshQuery = runQuery
	response.QueryTime = time.Now().Format(time.RFC3339)

	courses, err := fetchCourses(db)
	if err != nil {
		log.Println("Error fetching courses:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response.Courses = courses

	classifications, err := fetchAbbreviations(db, "class")
	if err != nil {
		log.Println("Error fetching classifications:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response.Classifications = classifications

	languages, err := fetchAbbreviations(db, "lang")
	if err != nil {
		log.Println("Error fetching languages:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response.Languages = languages

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Println("Error marshaling response data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write to cache if a fresh query was made
	if runQuery {
		if err := os.WriteFile(cacheFilename, jsonData, 0644); err != nil {
			log.Println("Error writing to cache file:", err)
		}
	}

	w.Write(jsonData)
}

func isCacheValid() bool {
	info, err := os.Stat(cacheFilename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

type CourseSql struct {
	Semester sql.NullString `db:"semester"`
	CourseCode sql.NullString `db:"course_code"`
	CourseName sql.NullString `db:"course_name"`
	CourseTopic sql.NullString `db:"course_topic"`
	Classification sql.NullString `db:"classification"`
	CodeAvailable sql.NullInt32 `db:"code_available"`
	Languages sql.NullString `db:"languages"`
	Summary sql.NullString `db:"summary"`
	ClassificationIcons sql.NullString `db:"classification_icons"`
	LangIcons sql.NullString `db:"lang_icons"`
	CourseClasses sql.NullString `db:"course_classes"`
}

func (c *CourseSql) Course() Course {
	return Course{
		Semester: c.Semester.String,
		CourseCode: c.CourseCode.String,
		CourseName: c.CourseName.String,
		CourseTopic: c.CourseTopic.String,
		Classification: c.Classification.String,
		CodeAvailable: c.CodeAvailable.Int32 == 1,
		Languages: c.Languages.String,
		Summary: c.Summary.String,
		ClassificationIcons: c.ClassificationIcons.String,
		LangIcons: c.LangIcons.String,
		CourseClasses: c.CourseClasses.String,
	}
}

func fetchCourses(db *sqlx.DB) ([]Course, error) {
	rows, err := db.Query(courseQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course CourseSql
		if err := rows.Scan(&course.Semester, &course.CourseCode, &course.CourseName, &course.CourseTopic, &course.Classification, &course.CodeAvailable,
			&course.Languages, &course.Summary, &course.ClassificationIcons, &course.LangIcons, &course.CourseClasses); err != nil {
			return nil, err
		}

		courses = append(courses, course.Course())
	}
	return courses, nil
}

func fetchAbbreviations(db *sqlx.DB, abbrType string) ([]Abbreviation, error) {
	rows, err := db.Query("SELECT abbreviation, icon_html, longname FROM abbreviations WHERE type = $1", abbrType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var abbreviations []Abbreviation
	for rows.Next() {
		var abbreviation Abbreviation
		if err := rows.Scan(&abbreviation.Abbreviation, &abbreviation.IconHTML, &abbreviation.Longname); err != nil {
			return nil, err
		}
		abbreviations = append(abbreviations, abbreviation)
	}
	return abbreviations, nil
}

const (
	port = "5432"
	cacheFilename = "cache.json"
)

var (
	host = MustGet("POSTGRES_HOST")
	user = MustGet("POSTGRES_USER")
	password = MustGet("POSTGRES_PASSWORD")
	dbname = MustGet("POSTGRES_DATABASE")
	refreshSecret = MustGet("REFRESH_SECRET")
)

func MustGet(secretName string) string {
	val, ok := os.LookupEnv(secretName)
	if !ok {
		panic("missing required environment variable: " + secretName)
	}
	return val
}

func GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)
}

func GetDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", GetConnectionString())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
