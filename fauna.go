package main

import (
// "fmt"
// "time"

// "github.com/fauna/fauna-go"
)

// const key = "fnAFeptU9dAAy5W1QRC6xBd4ucthzDPU1YTNZ_a2"

// func saveJobs(jobs []JobPosting) {
// 	client := fauna.NewClient(key, fauna.Timeouts{QueryTimeout: 20 * time.Second})

// 	for _, job := range jobs {
// 		createJob, _ := fauna.FQL(`Job.create({ title: ${title}})`, map[string]any{"title": job.title})
// 		res, err := client.Query(createJob)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(res.Data.(*fauna.Document).Data["name"])
// 	}
// }
