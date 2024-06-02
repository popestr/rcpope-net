package courses

import "github.com/jmoiron/sqlx"

const courseQuery = `
SELECT semester, course_code, course_name, course_topic, classification, code_available, languages, summary, 
STRING_AGG(CASE WHEN type = 'class' THEN icon_html ELSE NULL end, ' ') AS classification_icons, 
STRING_AGG(CASE WHEN type = 'lang' THEN icon_html ELSE NULL END, ' ') AS lang_icons, 
STRING_AGG(CONCAT('course--', abbreviation), ' ') AS course_classes
FROM courses 
INNER JOIN abbreviations ON (courses.classification LIKE CONCAT('%', abbreviations.abbreviation, '%') OR courses.languages LIKE CONCAT('%', abbreviations.abbreviation, '%'))
GROUP BY course_code, semester, course_name, course_topic, classification, code_available, languages, summary, semester_id
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
	err := db.Select(&abbreviations, "SELECT abbreviation, icon_html, longname FROM abbreviations WHERE type = $1", abbrType)
	if err != nil {
		return nil, err
	}

	return abbreviations, nil
}