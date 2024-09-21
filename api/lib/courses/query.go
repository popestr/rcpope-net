package courses

import "github.com/jmoiron/sqlx"

const courseQuery = `
SELECT semester, course_code, course_name, course_topic, classification, code_available, languages, summary, classification
FROM courses 
ORDER BY semester_id DESC
`

func FetchCourses(db *sqlx.DB) ([]Course, error) {
	var courses []CourseSql
	err := db.Select(&courses, courseQuery)
	if err != nil {
		return nil, err
	}

	var courseList []Course
	for _, course := range courses {
		courseList = append(courseList, course.Course())
	}

	return courseList, nil
}

func FetchAbbreviations(db *sqlx.DB, abbrType string) ([]Abbreviation, error) {
	var abbreviations []Abbreviation
	err := db.Select(&abbreviations, "SELECT abbreviation, longname FROM abbreviations WHERE type = $1", abbrType)
	if err != nil {
		return nil, err
	}

	return abbreviations, nil
}
