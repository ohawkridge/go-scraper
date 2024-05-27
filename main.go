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
