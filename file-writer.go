package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

// Outputs array of JobPostings as cards in a single html file.
//
// Parameters:
//   - jobs: an array of JobPosting objects.
func jobsToFile(jobs []JobPosting, filename string) {
	// Create new template and parse for errors
	tmpl, err := template.New("template-job-cards.tmpl").ParseFiles("template-job-cards.tmpl")
	if err != nil {
		panic(err)
	}
	var f *os.File
	f, err = os.Create(fmt.Sprintf("html/%s.html", filename))
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
		fileStr := fmt.Sprintf("html/j/job-%d.html", i)
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
	deleteFilesInDir("html/location")

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

func schoolsToFile(schools []School) {
	deleteFilesInDir("html/school")
	// Create a new template and check for errors
	tmpl, err := template.New("template-schools.tmpl").ParseFiles("template-schools.tmpl")
	if err != nil {
		panic(err)
	}
	var f *os.File
	f, err = os.Create("html/schools.html")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, schools)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}

// Delete all file in a directory.
// Called prior to writing new files.
//
// Parameters:
//   - dir: the path of the directory to delete.
func deleteFilesInDir(dir string) error {
	// Read all files in the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// Iterate through the files and delete each one
	for _, file := range files {
		if !file.IsDir() { // Ensure it is a file, not a directory
			err = os.Remove(filepath.Join(dir, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
