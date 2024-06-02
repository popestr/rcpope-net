package courses

import "database/sql"

type Course struct {
	Semester            string `json:"semester" db:"semester"`
	CourseCode          string `json:"course_code" db:"course_code"`
	CourseName          string `json:"course_name" db:"course_name"`
	CourseTopic         string `json:"course_topic" db:"course_topic"`
	Classification      string `json:"classification" db:"classification"`
	CodeAvailable       bool   `json:"code_available" db:"code_available"`
	Languages           string `json:"languages" db:"languages"`
	Summary             string `json:"summary" db:"summary"`
	ClassificationIcons string `json:"classification_icons" db:"classification_icons"`
	LangIcons           string `json:"lang_icons" db:"lang_icons"`
	CourseClasses       string `json:"course_classes" db:"course_classes"`
}

type Abbreviation struct {
	Abbreviation string `json:"abbreviation" db:"abbreviation"`
	IconHTML     string `json:"icon_html" db:"icon_html"`
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
	Classification      sql.NullString `db:"classification"`
	CodeAvailable       sql.NullInt32  `db:"code_available"`
	Languages           sql.NullString `db:"languages"`
	Summary             sql.NullString `db:"summary"`
	ClassificationIcons sql.NullString `db:"classification_icons"`
	LangIcons           sql.NullString `db:"lang_icons"`
	CourseClasses       sql.NullString `db:"course_classes"`
}

func (c *CourseSql) Course() Course {
	return Course{
		Semester:            c.Semester.String,
		CourseCode:          c.CourseCode.String,
		CourseName:          c.CourseName.String,
		CourseTopic:         c.CourseTopic.String,
		Classification:      c.Classification.String,
		CodeAvailable:       c.CodeAvailable.Int32 == 1,
		Languages:           c.Languages.String,
		Summary:             c.Summary.String,
		ClassificationIcons: c.ClassificationIcons.String,
		LangIcons:           c.LangIcons.String,
		CourseClasses:       c.CourseClasses.String,
	}
}