package main

import (
	"flag"
	"fmt"
)

const url string = "https://www.teachinherts.com/find-a-job.htm"

func main() {
	// Connect to database
	openDb()

	// Flags to control program actions
	var command string
	flag.StringVar(&command, "command", "", "Command to execute")
	flag.Parse()

	switch {
	case command == "scrape":
		scrapeUrl(url)
	case command == "delete":
		deleteAllJobs()
	case command == "locations":
		locations, err := getLocations()
		if err != nil {
			fmt.Println("Error getting job locations.\n", err)
		}
		// Output browse by location html file
		locationsToFile(locations)

		// For each location, generate html files
		for _, location := range locations {
			jobs, err := getLocationJobs(location)
			if err != nil {
				fmt.Printf("Error getting jobs for %s.\n", location.Location)
			}
			jobsToFile(jobs, location.Url)
		}
	case command == "full":
		// Create a detail.html for *every* job
		jobs := getJobs(-1)
		jobsToFiles(jobs)
	case command == "schools":
		// Find schools named in jobs
		schools, err := getSchools()
		if err != nil {
			fmt.Println("Error getting schools.\n", err)
		}

		// Output browse by school html file
		schoolsToFile(schools)
	case command == "base_salaries":
		processAllRecords()

	}

	// If still open, close the database connection
	open, err := verifyDb()
	if err != nil {
		fmt.Println("Error, can't Ping database.")
	}
	if open {
		closeDb()
	}
}
