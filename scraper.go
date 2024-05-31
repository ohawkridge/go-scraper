package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
	"github.com/jedib0t/go-pretty/v6/table"
)

type JobPosting struct {
	ID, Title, School, Location, Hours, Salary, Description, DetailsUrl, ClosingDate string
}

// Converts a three-letter month name into a two-digit string.
// Parameters:
//   - month: the month name.
//
// Returns:
//   - a two-digit string like "01".
func monthToNumber(month string) string {
	month = strings.ToLower(month)
	monthMap := map[string]string{
		"jan":       "01",
		"feb":       "02",
		"mar":       "03",
		"apr":       "04",
		"may":       "05",
		"jun":       "06",
		"jul":       "07",
		"aug":       "08",
		"sep":       "09",
		"oct":       "10",
		"nov":       "11",
		"dec":       "12",
		"january":   "01",
		"february":  "02",
		"march":     "03",
		"april":     "04",
		"june":      "06",
		"july":      "07",
		"august":    "08",
		"september": "09",
		"october":   "10",
		"november":  "11",
		"december":  "12",
	}
	return monthMap[strings.ToLower(month)]
}

// Like 5 June 2024 at 9am
func govStringToTime(date string) string {
	parts := strings.Fields(date)

	// Make day 2 digits
	day := parts[0]
	if len(day) == 1 {
		day = "0" + parts[0]
	}

	month := monthToNumber(parts[1])
	year := parts[2]

	// TODO Parse times like 9am/11:59pm/12pm(midday)

	date = fmt.Sprintf("%s-%s-%s %02d:00:00", year, month, day, 0)
	return date
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

// Use GOColly package to scrape all jobs from url
func scrapeUrl(url string) {
	var jobs []JobPosting

	// Initialise a collector object
	collector := colly.NewCollector()

	// Apply rate limits
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*teachinherts.*",
		Parallelism: 2,
		Delay:       3,               // Start with a 3 second delay
		RandomDelay: 5 * time.Second, // Additional random delay between 0 and 5 seconds
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

	// Write jobs into the database
	insertJobs(jobs)
}

func scrapeUrl2(url2 string) {
	// var jobs []JobPosting
	// Output nice table in console
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Title", "Closing", "URL"})

	// Initialise a collector object
	collector := colly.NewCollector()

	// Apply rate limits
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*service.gov.uk.*",
		Parallelism: 2,
		Delay:       3,               // Start with a 3 second delay
		RandomDelay: 5 * time.Second, // Additional random delay between 0 and 5 seconds
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL, "✔️")
	})

	collector.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL, "✔️")
	})

	// Handle pagination by looking for next page link
	//
	// collector.OnHTML("a.govuk-pagination__link", func(e *colly.HTMLElement) {
	// 	nextPage := e.Attr("href")
	// 	if nextPage != "" {
	// 		e.Request.Visit(nextPage)
	// 	}
	// })

	collector.OnHTML("div.search-results__item", func(e *colly.HTMLElement) {
		// initialise a new job struct every time we visit a page
		job := JobPosting{}

		// find the job parts and assign them to the struct
		job.Title = removeWords(e.ChildText("h2.govuk-heading-m"))

		// Separate school and location (in same <p> tag)
		address := e.ChildText("p.govuk-body.address")
		index := strings.IndexRune(address, ',')
		job.School = address[:index]
		job.Location = address[index+2:]

		// job.Hours = e.ChildText("dl.govuk-summary-list:nth-child(3) div.govuk-summary-list__row dd")
		job.Hours = e.ChildText("div.govuk-summary-list__row:nth-of-type(3) dd")
		job.Salary = e.ChildText("div.govuk-summary-list__row:first-of-type dd")

		// convert dates like "11 June 2024 at 12pm (midday)"
		date := e.ChildText("div.govuk-summary-list__row:nth-of-type(4) dd")
		job.ClosingDate = govStringToTime(date)

		url := "https://teaching-vacancies.service.gov.uk" + e.ChildAttr("h2 a.govuk-link.view-vacancy-link", "href")
		job.DetailsUrl = url
		job.Description = fmt.Sprintf("See <a href='%s'>full job details</a>", url)

		// TESTING
		if len(job.Title) > 30 {
			job.Title = job.Title[:29] + "..."
		}
		if len(job.School) > 26 {
			job.School = job.School[:25] + "..."
		}
		if len(job.Salary) > 26 {
			job.Salary = job.Salary[:25] + "..."
		}
		t.AppendRow(table.Row{job.Title, job.ClosingDate, job.DetailsUrl})

		// jobs = append(jobs, job)
	})
	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("An error occurred:", e)
	})
	collector.Visit(url2)
	// t.AppendFooter(table.Row{"COUNT", len(jobs2)})
	t.AppendFooter(table.Row{"COUNT", t.Length()})
	t.SetStyle(table.StyleLight)
	t.Render()
	// Write jobs into the database
	// insertJobs(jobs)
}

func removeWords(s string) string {
	re := regexp.MustCompile(`(?i)(quick apply)`)
	out := re.ReplaceAllString(s, "")
	return out
}
