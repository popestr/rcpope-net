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

func FetchAbbreviations(db *sqlx.DB, abbrType string) ([]Abbreviation, error) {
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