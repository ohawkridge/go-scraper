package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const dsn string = "admin:AcStrWTVJgodnDCuQXyb@tcp(jobs.clwa0i6yu538.eu-west-2.rds.amazonaws.com:3306)/jobs"

// Struct to represent job locations
type Location struct {
	Location string
	Url      string
}

var (
	db  *sql.DB
	err error
)

// Open database connection
func openDb() error {
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	return nil
}

// Verify database connection
func verifyDb() (bool, error) {
	err = db.Ping()
	if err != nil {
		return false, err
	}
	fmt.Println("Connected to MySQL AWS 🔌")

	return true, err
}

// Close database connection
func closeDb() {
	db.Close()
}

// Get jobs from the database.
// Parameters:
//   - limit: # of rows to get (-1 = all rows).
//
// Returns:
//   - array of JobPosting objects.
func getJobs(limit int) []JobPosting {
	verifyDb()

	// Define the query, limiting rows if needed
	query := "SELECT * FROM job ORDER BY id ASC"
	if limit != -1 {
		query = fmt.Sprintf("SELECT * FROM job ORDER BY id ASC LIMIT %d", limit)
	}

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		panic(err.Error())
	}

	var jobs []JobPosting
	for rows.Next() {
		var job JobPosting
		// Scan the row into JobPosting object
		err = rows.Scan(&job.ID, &job.Title, &job.School, &job.Location, &job.Hours, &job.Salary, &job.Description, &job.DetailsUrl, &job.ClosingDate)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			panic(err.Error())
		}
		jobs = append(jobs, job)
	}

	return jobs
}

// Get locations based on all jobs in the database.
//
// Returns:
//   - array of Location objects.
func getLocations() ([]Location, error) {
	verifyDb()

	var locations []Location
	// Find distinct locations (excluding empty)
	query := "SELECT DISTINCT location FROM job WHERE NOT location='';"

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}

	for rows.Next() {
		var location Location
		err = rows.Scan(&location.Location)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		location.Url = locationToSlug(location.Location)
		locations = append(locations, location)

		// For each location, find matching jobs
		getLocationJobs(location)
	}

	// fmt.Printf("Locations: %#v\n", locations)
	return locations, err
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

func getLocationJobs(location Location) ([]JobPosting, error) {
	verifyDb()

	// Find jobs for location
	query := fmt.Sprintf("SELECT * FROM job WHERE location='%s';", location.Location)

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}

	var jobs []JobPosting
	for rows.Next() {
		var job JobPosting
		// Scan the row into a JobPosting object
		err = rows.Scan(&job.ID, &job.Title, &job.School, &job.Location, &job.Hours, &job.Salary, &job.Description, &job.DetailsUrl, &job.ClosingDate)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, err
}

// Converts a MySQL date like 2024-05-24 into
// a relative date like '26 days time'.
//
// Parameters:
//   - dateStr: the date as a string like 'YYYY-MM-DD'.
//
// Returns:
//   - a string like '28 days'.
func getRelativeDate(dateStr string) string {
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
		relativeTime = fmt.Sprintf("%d days time", days)
	} else if days < 0 {
		relativeTime = fmt.Sprintf("%d days ago", -days)
	} else {
		relativeTime = "today"
	}

	return relativeTime
}

// Inserts jobs into the database.
// Parameters:
//   - jobs: an array of JobPosting objects.
func insertJobs(jobs []JobPosting) {
	verifyDb()

	// Prepare INSERT query
	stmt, err := db.Prepare("INSERT INTO job (title, school, location, hours, salary, description, url, closing_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return
	}
	defer stmt.Close()

	// Insert each job
	// TODO Bulk insert?
	for _, job := range jobs {
		_, err := stmt.Exec(job.Title, job.School, job.Location, job.Hours, job.Salary, job.Description, job.DetailsUrl, job.ClosingDate)
		if err != nil {
			fmt.Println("Error executing statement:", err)
			return
		}
	}
	fmt.Printf("Inserted %d record(s). ✔️\n", len(jobs))
}
