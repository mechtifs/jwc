package main

import (
	"fmt"
	"runtime"
	"time"
)

func showSummary(courses []*Course, total, retries int) {
	fmt.Printf(" *** Course Selection Summary as of %s ***\n", time.Now().Format("Mon Jan _2 15:04:05 MST 2006"))
	fmt.Println("================================================================================")
	fmt.Printf("[%s %d/%d(%.2f%%) PAR:%d RE:%d]\n", config.Session.Credential, total-len(courses), total, float32(total-len(courses))/float32(total), runtime.NumGoroutine()-2, retries)
	var coursesIDs []string
	for _, course := range courses {
		coursesIDs = append(coursesIDs, course.CourseID)
	}
	fmt.Println("PENDING:", coursesIDs)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println()
}
