package main

import (
	"fmt"
	"html/template"
	"os"
)

// Outputs array of JobPostings as cards in a single html file.
//
// Parameters:
//   - jobs: an array of JobPosting objects.
func jobsToFile(jobs []JobPosting) {
	// Create new template and parse for errors
	tmpl, err := template.New("template-job-cards.tmpl").ParseFiles("template-job-cards.tmpl")
	if err != nil {
		panic(err)
	}
	var f *os.File
	f, err = os.Create("html/jobs.html")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, jobs)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}

// Outputs array of JobPostings as individual html files.
// N.B. Templates don't work inside folders !!
//
// Parameters:
//   - jobs: an array of JobPosting objects.
func jobsToFiles(jobs []JobPosting) {
	for i, job := range jobs {
		tmpl, err := template.New("template-detail.tmpl").ParseFiles("template-detail.tmpl")
		if err != nil {
			panic(err)
		}
		var f *os.File
		fileStr := fmt.Sprintf("html/jobs/job-%d.html", i)
		f, err = os.Create(fileStr)
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(f, job)
		if err != nil {
			panic(err)
		}
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}
}

// Outputs array of Location objects to a single html file.
//
// Parameters:
//   - locations: an array of Location objects.
func locationsToFile(locations []Location) {
	// Create a new template and check for errors
	tmpl, err := template.New("template-locations.tmpl").ParseFiles("template-locations.tmpl")
	if err != nil {
		panic(err)
	}
	var f *os.File
	f, err = os.Create("html/locations.html")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, locations)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}
