package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jmoiron/sqlx"
	"github.com/popestr/rcpope-net/lambda/lib/courses"
	pgdb "github.com/popestr/rcpope-net/lambda/lib/db"
)

var db *sqlx.DB

func HandleRequest(ctx context.Context) (*courses.Response, error) {
	var response courses.Response
	response.QueryTime = time.Now().Format(time.RFC3339)

	courseList, err := courses.FetchCourses(db)
	if err != nil {
		return nil, fmt.Errorf("error fetching courses: %w", err)
	}
	response.Courses = courseList

	classifications, err := courses.FetchAbbreviations(db, "class")
	if err != nil {
		return nil, fmt.Errorf("error fetching classifications: %w", err)
	}
	response.Classifications = classifications

	languages, err := courses.FetchAbbreviations(db, "lang")
	if err != nil {
		return nil, fmt.Errorf("error fetching languages: %w", err)
	}
	response.Languages = languages

	return &response, nil
}

func main() {
	var err error
	db, err = pgdb.DB()
	if err != nil {
		panic(err)
	}

	lambda.Start(HandleRequest)
}
