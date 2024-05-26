package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Struct to represent job locations
type Location struct {
	Location string
	Url      string
}

func openDb() {
	// Open database connection
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
}

func verifyDb() {
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected to MySQL AWS 🔌")
}

func getLocations() {
	verifyDb()

	var locations []Location
	// Find distinct locations excluding empty
	query := "SELECT DISTINCT location FROM job WHERE NOT location='';"

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}

	for rows.Next() {
		var location Location
		err = rows.Scan(&location.Location)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		location.Url = locationToSlug(location.Location)
		locations = append(locations, location)

		// For each location, find matching jobs and display them on their own page
		getLocationJobs(location)
	}

	// fmt.Printf("Locations: %#v\n", locations)
	writeLocations(locations)
}

// Converts a location like "Welwyn / Hatfield District"
// into a URL path like "welwyn-hatfield-district"
// Parameters:
//   - location: the location string.
//
// Returns:
//   - a hypen separated string.
func locationToSlug(location string) string {
	location = strings.ToLower(location)

	// Regular expresssion to find all non-alphabet characters
	re := regexp.MustCompile(`[^a-z]+`)
	return re.ReplaceAllString(location, "-")
}

func getLocationJobs(location Location) {
	verifyDb()

	// Find jobs for location
	query := fmt.Sprintf("SELECT * FROM job WHERE location='%s';", location.Location)

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}

	var jobs []JobPosting
	for rows.Next() {
		var job JobPosting
		// Scan the row into variables
		err = rows.Scan(&job.ID, &job.Title, &job.School, &job.Location, &job.Hours, &job.Salary, &job.Description, &job.DetailsUrl, &job.ClosingDate)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		// Convert 2024-05-24 → 23 days
		job.ClosingDate = relativeDate(job.ClosingDate)
		jobs = append(jobs, job)
	}

	tmpl, err := template.New("job-cards.tmpl").ParseFiles("job-cards.tmpl")
	if err != nil {
		panic(err)
	}
	var f *os.File
	fileStr := fmt.Sprintf("locations/%s.html", location.Url)
	f, err = os.Create(fileStr)
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

func relativeDate(dateStr string) string {
	// Parse the date string into a time.Time object
	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return ""
	}

	// Calculate the difference between the parsed date and the current date
	now := time.Now()
	duration := parsedDate.Sub(now)

	// Convert the duration to days
	days := int(duration.Hours() / 24)

	// Determine the relative time message
	var relativeTime string
	if days > 0 {
		relativeTime = fmt.Sprintf("Closes in %d days", days)
	} else if days < 0 {
		relativeTime = fmt.Sprintf("Closed %d days ago", -days)
	} else {
		relativeTime = "Closes today"
	}

	return relativeTime
}
