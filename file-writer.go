package main

import (
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
