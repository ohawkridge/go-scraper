package main

import (
	"fmt"
	"html/template"
	"os"
)

func testTemplates(jobs []JobPosting) {
	// Create new template and parse for errors
	tmpl, err := template.New("all-jobs.tmpl").ParseFiles("all-jobs.tmpl")
	if err != nil {
		panic(err)
	}
	var f *os.File
	f, err = os.Create("all-jobs.html")
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

func writeDetails(jobs []JobPosting) {
	for i, job := range jobs {
		// TESTING
		if i == 3 {
			break
		}
		tmpl, err := template.New("detail.tmpl").ParseFiles("detail.tmpl")
		if err != nil {
			panic(err)
		}
		var f *os.File
		fileStr := fmt.Sprintf("jobs/job-%d.html", i)
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

func writeLocations(locations []Location) {
	// Create a new template and check for errors
	tmpl, err := template.New("locations.tmpl").ParseFiles("locations.tmpl")
	if err != nil {
		panic(err)
	}
	var f *os.File
	f, err = os.Create("locations.html")
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
