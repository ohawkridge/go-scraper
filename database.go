package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
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

type School struct {
	Name    string
	Url     string
	NumJobs int
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
		job.ClosingDate = getRelativeDate(job.ClosingDate)
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
		// Turn location into url-safe-slug
		location.Url = locationToSlug(location.Location)
		locations = append(locations, location)
	}

	return locations, err
}

// Get schools based on all jobs in the database.
//
// Returns:
//   - array of School objects.
func getSchools() ([]School, error) {
	verifyDb()

	var schools []School
	// Find distinct schools named in jobs
	query := "SELECT DISTINCT school FROM job WHERE NOT school='' ORDER BY school ASC;"

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}

	for rows.Next() {
		var school School
		err = rows.Scan(&school.Name)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		// Turn school's name into url-safe-slug
		school.Url = locationToSlug(school.Name)

		// Get the number of jobs at this school
		n, err := countSchoolJobs(school.Name)
		if err == nil {
			school.NumJobs = n
		} else {
			fmt.Printf("Error counting jobs at %s\n%s\n", school.Name, err)
		}

		schools = append(schools, school)
	}

	return schools, err
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
		job.ClosingDate = getRelativeDate(job.ClosingDate)
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

// Delete *all* jobs ahead of re-scraping
func deleteAllJobs() {
	verifyDb()

	// TRUNCATE is faster because it removes all rows by deallocating the data pages used by the table
	query := "TRUNCATE TABLE job;"
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}

	fmt.Println("All jobs deleted ✔️")
}

// Count the number of jobs at a school
func countSchoolJobs(school string) (int, error) {
	// verifyDb() should be connected

	// Construct COUNT query
	query := "SELECT COUNT(*) AS num_jobs FROM job WHERE school = ?"

	var count int
	row := db.QueryRow(query, school)
	err = row.Scan(&count)
	if err != nil {
		return -1, err
	}
	// fmt.Printf("Found %d jobs at %s\n", count, school)
	return count, nil
}

func numsOnly(s string) string {
	// Remove all non-digit characters
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(s, "")
}

// Check if salary is hourly
func isHourly(s string) bool {
	// Construct a regular expression pattern to match any of the substrings
	re := regexp.MustCompile(`(?i)(hourly|per hour|hour|p/h|hrs|closing)`)
	return re.MatchString(s)
}

func contains5Digits(s string) bool {
	// Improved to check for 5 digits *after* £
	if strings.Contains(s, "£") {
		s = s[strings.IndexRune(s, '£'):]
	}
	// Compile the regular expression to match at least 5 digits
	re := regexp.MustCompile(`(\D*\d){5}`)
	return re.MatchString(s)
}

func extractSalary(s string) (int, error) {
	// Find the index of the first '£'
	index := strings.IndexRune(s, '£')
	if index == -1 {
		return 0, err
	} else {
		// Extract the next 6 characters (£19,999) as an int
		fmt.Println(s)
		s = s[index : index+8]
		fmt.Println(s)
		salary, err := strconv.Atoi(numsOnly(s))
		return salary, err
	}
}

// Admin function to process all records
func processAllRecords() {
	verifyDb()

	// Get all records
	rows, err := db.Query("SELECT id, salary FROM job;")
	if err != nil {
		fmt.Println("Error executing query:", err)
	}
	for rows.Next() {
		var salary string
		var id int
		err := rows.Scan(&id, &salary)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			panic(err.Error())
		}
		if !isHourly(salary) && contains5Digits(salary) {
			extracted, err := extractSalary(salary)
			if err != nil {
				fmt.Println("Error extracting salary", err)
			}

			// Execute the UPDATE query
			query := `UPDATE job SET base_salary = ? WHERE id = ?`
			result, err := db.Exec(query, extracted, id)
			if err != nil {
				fmt.Println("Error executing query", err)
			} else {
				n, err := result.RowsAffected()
				if err != nil {
					log.Fatalf("Error fetching number of rows affected: %v", err)
				}
				fmt.Printf("Updated job %d (%d row affected)\n", id, n)
			}
		}
	}
}
