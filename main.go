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
	scrape := flag.Bool("scrape", false, "Start scraping url")
	flag.Parse()

	if *scrape {
		fmt.Println("Scraping started...")
		// scrapeUrl(url)
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
