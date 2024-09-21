package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/popestr/rcpope-net/lambda/lib/courses"
	"github.com/popestr/rcpope-net/lambda/lib/db"
)

func HandleRequest(ctx context.Context) (*courses.Response, error) {
	db, err := db.GetDB()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	defer db.Close()

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
	lambda.Start(HandleRequest)
}
