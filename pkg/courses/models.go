package courses

import (
	"database/sql"
	"strings"
)

type Course struct {
	Semester            string `json:"semester"`
	CourseCode          string `json:"course_code"`
	CourseName          string `json:"course_name"`
	CourseTopic         string `json:"course_topic"`
	CodeAvailable       bool   `json:"code_available"`
	Summary             string `json:"summary"`
	Classifications      []string `json:"classifications"`
	Languages 			[]string `json:"languages"`
}

type Abbreviation struct {
	Abbreviation string `json:"abbreviation" db:"abbreviation"`
	Longname     string `json:"longname" db:"longname"`
}

type Response struct {
	QueryTime       string         `json:"query_time"`
	Courses         []Course       `json:"courses"`
	Classifications []Abbreviation `json:"classifications"`
	Languages       []Abbreviation `json:"languages"`
}

type CourseSql struct {
	Semester            sql.NullString `db:"semester"`
	CourseCode          sql.NullString `db:"course_code"`
	CourseName          sql.NullString `db:"course_name"`
	CourseTopic         sql.NullString `db:"course_topic"`
	CodeAvailable       sql.NullInt32  `db:"code_available"`
	Summary             sql.NullString `db:"summary"`
	Classifications 	sql.NullString `db:"classification"`
	Languages           sql.NullString `db:"languages"`
}

func (c *CourseSql) Course() Course {
	var classifications, languages []string

	if c.Classifications.Valid {
		classifications = strings.Split(c.Classifications.String, " ")
	}

	if c.Languages.Valid {
		languages = strings.Split(c.Languages.String, " ")
	}

	return Course{
		Semester:            c.Semester.String,
		CourseCode:          c.CourseCode.String,
		CourseName:          c.CourseName.String,
		CourseTopic:         c.CourseTopic.String,
		CodeAvailable:       c.CodeAvailable.Int32 == 1,
		Summary:             c.Summary.String,
		Classifications:     classifications,
		Languages:           languages,
	}
}