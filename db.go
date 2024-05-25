package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func connectDb() {
	dsn := "admin:AcStrWTVJgodnDCuQXyb@tcp(jobs.clwa0i6yu538.eu-west-2.rds.amazonaws.com:3306)/jobs"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Verify the connection to the database
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected to MySQL @ AWS ✔️")
}

func insertJob(job JobPosting) {
	// dsn := "admin:AcStrWTVJgodnDCuQXyb@tcp(jobs.clwa0i6yu538.eu-west-2.rds.amazonaws.com:3306)/jobs"
	// db, err := sql.Open("mysql", dsn)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()

	// // Verify the connection to the database
	// err = db.Ping()
	// if err != nil {
	// 	fmt.Println("Error connecting to the database:", err)
	// 	return
	// }
	// fmt.Println("Connected to MySQL @ AWS.")

	// Insert a record into the database
	// N.B. title, school, location, hours, salary, description, detailsUrl, closingDate
	query := "INSERT INTO job (title, school, location, hours, salary, description, url, closing_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := db.Exec(query, job.title, job.school, job.location, job.hours, job.salary, job.description, job.detailsUrl, job.closingDate)
	if err != nil {
		panic(err.Error())
	}

	// Get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error fetching rows affected:", err)
		return
	}
	fmt.Printf("%d row(s) affected.", rowsAffected)
}

func deleteAllJobs() {
	// dsn := "admin:AcStrWTVJgodnDCuQXyb@tcp(jobs.clwa0i6yu538.eu-west-2.rds.amazonaws.com:3306)/jobs"
	// db, err := sql.Open("mysql", dsn)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()

	// // Verify the connection to the database
	// err = db.Ping()
	// if err != nil {
	// 	fmt.Println("Error connecting to the database:", err)
	// 	return
	// }
	// fmt.Println("Connected to MySQL @ AWS.")
	query := "DELETE FROM job;"
	result, err := db.Exec(query)
	if err != nil {
		panic(err.Error())
	}

	// Get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error fetching rows affected:", err)
		return
	}
	fmt.Printf("%d row(s) affected ✔️\n", rowsAffected)
}
