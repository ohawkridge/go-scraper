package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gocolly/colly"
	"github.com/jedib0t/go-pretty/v6/table"
)

const url = "https://www.teachinherts.com/find-a-job.htm"

var jobs []JobPosting

type JobPosting struct {
	title, school, location, hours, salary, description, detailsUrl, closingDate string
}

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

func stringToTime(date string) string {
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

	date = fmt.Sprintf("%s-%s-%s %d:00:00", day, month, parts[4], hour)

	// ^^^ This date string will parse as a Go datetime in the future
	// Go datetime formatting mnemonic
	// https://medium.com/@simplyianm/how-go-solves-date-and-time-formatting-8a932117c41c
	// n.b. 15 for 24-hour format
	// timeObj, err := time.Parse("01-02-2006 15:04:05", date)
	// if err != nil {
	//     panic(err)
	// }

	return date
}

func printJobs(jobs []JobPosting) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Title", "School", "Closes"})
	for _, job := range jobs {
		// N.B. title, school, location, hours, salary, description, detailsUrl, closingDate
		t.AppendRows([]table.Row{
			{job.title, job.school, job.closingDate},
		})
	}
	t.AppendFooter(table.Row{"COUNT", len(jobs)})
	t.SetStyle(table.StyleLight)
	t.Render()
}

func main() {
	// initialise a collector object
	collector := colly.NewCollector()

	collector.OnRequest(func(r *colly.Request) {
		// print the url of that request
		fmt.Println("Visiting", r.URL)
	})
	collector.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})
	collector.OnHTML("div.listing.joblisting > ul > li", func(e *colly.HTMLElement) {
		// initialise a new job struct every time we visit a page
		job := JobPosting{}

		// find the job parts and assign them to the struct
		job.title = e.ChildText("h3")
		job.school = e.ChildText("h4")
		job.location = e.ChildText("div > p:nth-child(4) span:nth-child(1)")
		job.hours = e.ChildText("div > p:nth-child(4) span:nth-child(2)")
		job.salary = e.ChildText("div > p:nth-child(5)")

		// convert dates like "16th Apr 2024 12:00 PM" to a format
		// that can be parsed as Go datetime objects in the future
		job.closingDate = stringToTime(e.ChildText("div > p.date"))

		job.description = e.ChildText("div > p:nth-child(7)")
		job.detailsUrl = "https://www.teachinherts.com" + e.ChildAttr("a", "href")
		jobs = append(jobs, job)
	})
	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("An error occurred:", e)
	})
	collector.Visit(url)
	printJobs(jobs)
	writeTxt()
}
