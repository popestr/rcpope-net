package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/popestr/rcpope-net/pkg/courses"
	"github.com/popestr/rcpope-net/pkg/db"
)

func Courses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := db.GetDB()
	if err != nil {
		log.Println("Error connecting to database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var response courses.Response
	response.QueryTime = time.Now().Format(time.RFC3339)

	courseList, err := courses.FetchCourses(db)
	if err != nil {
		log.Println("Error fetching courses:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response.Courses = courseList

	classifications, err := courses.FetchAbbreviations(db, "class")
	if err != nil {
		log.Println("Error fetching classifications:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response.Classifications = classifications

	languages, err := courses.FetchAbbreviations(db, "lang")
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

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}