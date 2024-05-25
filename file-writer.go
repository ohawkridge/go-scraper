package main

import (
	"os"
	"text/template"
	"time"
)

func writeTxt() {
	ts := time.Now().UTC().String()
	err := os.WriteFile("tmp/test.txt", []byte("Last run:\n"+ts), 0755)
	if err != nil {
		panic(err)
	}
}

func testTemplates() {
	testJobs := []JobPosting{
		{
			"Teacher of Maths",
			"Chauncy School",
			"Ware",
			"Full time",
			"£40,000-£50,000 MPS",
			"Blah blah blah",
			"https://chauncyschool.com/job",
			"2024-12-31",
		},
		{
			"Teacher of Art",
			"King Harold Academy",
			"Ware",
			"Full time",
			"£40,000-£50,000 UPS",
			"Blah blah blah",
			"https://kha.co.uk/",
			"2024-12-31",
		},
	}

	var tmplFile = "templates/list.tmpl"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, testJobs)
	if err != nil {
		panic(err)
	}
}
