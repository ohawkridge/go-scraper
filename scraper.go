package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
)

const url string = "https://www.teachinherts.com/find-a-job.htm"
const dsn string = "admin:AcStrWTVJgodnDCuQXyb@tcp(jobs.clwa0i6yu538.eu-west-2.rds.amazonaws.com:3306)/jobs"

var (
	db   *sql.DB
	err  error
	jobs []JobPosting
)

type JobPosting struct {
	Title, School, Location, Hours, Salary, Description, DetailsUrl, ClosingDate string
}

// Converts a three-letter month name into a two-digit string.
// Parameters:
//   - month: the month name.
//
// Returns:
//   - a two-digit string like "01".
func monthToNumber(month string) string {
	monthMap := map[string]string{
		"jan": "01",
		"feb": "02",
		"mar": "03",
		"apr": "04",
		"may": "05",
		"jun": "06",
		"jul": "07",
		"aug": "08",
		"sep": "09",
		"oct": "10",
		"nov": "11",
		"dec": "12",
	}
	return monthMap[strings.ToLower(month)]
}

func hertsStringToTime(date string) string {
	parts := strings.Fields(date)

	// if day is one digit (eg, 1st), add a zero to the front
	day := parts[2][0:2]
	dayRune, _ := utf8.DecodeRuneInString(day[1:2])
	if !unicode.IsDigit(dayRune) {
		day = "0" + day[0:1]
	}

	// convert month to number
	month := monthToNumber(parts[3])

	// get the hour part and cast to integer
	hourDigits := strings.Split(parts[5], ":")
	hour, err := strconv.Atoi(hourDigits[0])
	if err != nil {
		panic(err)
	}

	// convert to 24-hour format
	if date[len(date)-2:] == "PM" {
		if hour < 12 {
			hour += 12
		}
	}

	date = fmt.Sprintf("%s-%s-%s %d:00:00", parts[4], month, day, hour)
	// ^^^ This date string should parse as GO datetime in the future
	// Go datetime formatting mnemonic
	// https://medium.com/@simplyianm/how-go-solves-date-and-time-formatting-8a932117c41c
	// timeObj, err := time.Parse("01-02-2006 15:04:05", date)

	return date
}

func main() {
	// Open database connection
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}

	// Initialise a collector object
	collector := colly.NewCollector()

	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*teachinherts.*",
		Parallelism: 2,
		Delay:       5 * time.Second,
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL, "✔️")
	})

	collector.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL, "✔️")
	})

	// Handle pagination by looking for next page link
	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		nextPage := e.Attr("href")
		if nextPage != "" {
			e.Request.Visit(nextPage)
		}
	})

	collector.OnHTML("div.listing.joblisting > ul > li", func(e *colly.HTMLElement) {
		// initialise a new job struct every time we visit a page
		job := JobPosting{}

		// find the job parts and assign them to the struct
		job.Title = e.ChildText("h3")
		job.School = e.ChildText("h4")
		job.Location = e.ChildText("div > p:nth-child(4) span:nth-child(1)")
		job.Hours = e.ChildText("div > p:nth-child(4) span:nth-child(2)")
		job.Salary = e.ChildText("div > p:nth-child(5)")

		// convert dates like "16th Apr 2024 12:00 PM" to a format
		// that can be parsed as Go datetime objects in the future
		job.ClosingDate = hertsStringToTime(e.ChildText("div > p.date"))

		job.Description = e.ChildText("div > p:nth-child(7)")
		job.DetailsUrl = "https://www.teachinherts.com" + e.ChildAttr("a", "href")
		jobs = append(jobs, job)
	})
	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("An error occurred:", e)
	})
	collector.Visit(url)

	// Write jobs to MySQL DB
	writeJobsToDb()

	// Write all-jobs template
	testTemplates(jobs)
}

func writeJobsToDb() {
	verifyDb()
	stmt, err := db.Prepare("INSERT INTO job (title, school, location, hours, salary, description, url, closing_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return
	}
	defer stmt.Close()

	for _, job := range jobs {
		_, err := stmt.Exec(job.Title, job.School, job.Location, job.Hours, job.Salary, job.Description, job.DetailsUrl, job.ClosingDate)
		if err != nil {
			fmt.Println("Error executing statement:", err)
			return
		}
	}
	fmt.Printf("Inserted %d record(s). ✔️\n", len(jobs))
	closeDb()
}

// Verifies connection to the database.
func verifyDb() {
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected to MySQL @ AWS ✔️")
}

func closeDb() {
	db.Close()
}
