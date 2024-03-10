package main

import (
	"log"
	"time"

	"jwc/requests"
)

var session *JWSession
var config Config

func updateCourses(courseIDs []string) []*Course {
	var courses []*Course
	for _, courseID := range courseIDs {
		course := &Course{CourseID: courseID}
		for {
			course.updateTeachIDAvailability()
			if course.TeachID != "" {
				break
			}
			time.Sleep(time.Duration(config.Client.Delay) * time.Millisecond)
		}
		if course.TeachID == "0" {
			continue
		}
		courses = append(courses, course)
	}
	return courses
}

func modeAdd() {
	log.Printf("[NOTICE] Choosing %d course(s)\n", len(config.Session.Targets))
	courses := updateCourses(config.Session.Targets)

	total := len(courses)
	if total == 0 {
		log.Println("[NOTICE] Nothing to do. Quitting...")
		return
	}

	resCh := make(chan *Course, config.Client.Parallel*total)
	for range config.Client.Parallel {
		for _, course := range courses {
			go session.addCourse(course, resCh)
		}
	}

	retries := 0
	if config.Summary.Enabled && config.Summary.Interval > 0 {
		go func() {
			for {
				time.Sleep(time.Duration(config.Summary.Interval) * time.Second)
				showSummary(courses, total, retries)
			}
		}()
	}

	for res := range resCh {
		if res.State == Chosen || res.State == Conflict {
			if !config.Client.Keep {
				for i, courseInfo := range courses {
					if courseInfo.CourseID == res.CourseID {
						courses = append(courses[:i], courses[i+1:]...)
						break
					}
				}
			}
			if res.State == Chosen {
				log.Println("[NOTICE]", res.CourseID, "is already chosen.")
			} else {
				log.Println("[WARNING]", res.CourseID, "is in conflict with other courses.")
				total--
			}
			if len(courses) == 0 {
				log.Println("[NOTICE] Nothing to do. Quitting...")
				break
			}
		} else if res.State == Overflowed {
			log.Println("[FATAL] Total credits overflowed. Quitting...")
			break
		} else {
			retries++
			time.Sleep(time.Duration(config.Client.Delay) * time.Millisecond)
			go session.addCourse(res, resCh)
		}
	}
}

func modeChange() {
	srcCourse := &Course{CourseID: config.Session.Targets[0]}
	srcCourse.updateListID()
	if srcCourse.ListID == "" {
		log.Println("[NOTICE] Course not found. Quitting...")
		return
	}
	dstCourse := &Course{CourseID: config.Session.Targets[1]}
	for {
		dstCourse.updateTeachIDAvailability()
		if dstCourse.TeachID == "0" {
			log.Println("[NOTICE] Course not found. Quitting...")
			return
		}
		if !dstCourse.IsAvailable {
			log.Println("[WARNING] Course is currently full. Retrying in", config.Client.Delay, "milliseconds...")
			time.Sleep(time.Duration(config.Client.Delay) * time.Millisecond)
			continue
		}
		break
	}
	resCh := make(chan *Course, config.Client.Parallel*10)
	go func() {
		for range 10 {
			for range config.Client.Parallel {
				go session.quitCourse(srcCourse)
				go session.addCourse(dstCourse, resCh)
			}
			time.Sleep(time.Duration(config.Client.Delay) * time.Millisecond)
		}
	}()
	for res := range resCh {
		if res.State == Chosen {
			log.Println("[NOTICE]", res.CourseID, "is successfully changed. Quitting...")
			break
		} else if res.State == Conflict {
			log.Println("[WARNING]", res.CourseID, "is in conflict with other courses.")
		} else if res.State == Overflowed {
			log.Println("[FATAL] Total credits overflowed. Quitting...")
			break
		}
	}
}

func main() {
	parseConfig()
	session = NewJWSession(&requests.SessionOptions{
		Header: map[string][]string{
			"User-Agent": {config.Session.UserAgent},
			"Cookie":     {"JSESSIONID=" + config.Session.Credential},
		},
	})

	switch config.Client.Mode {
	case 0:
		modeAdd()
	case 1:
		modeChange()
	default:
		log.Println("[FATAL] Invalid mode. Quitting...")
	}
}
