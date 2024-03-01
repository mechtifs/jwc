package main

import (
	"fmt"
	"jwc/requests"
	"log"
	"strings"
	"time"
)

func queryCourse(courseID string) string {
	var resp *requests.Response
	var err error
	for {
		resp, err = session.Post(config.Session.BaseURI+"/vatuu/CourseStudentAction",
			&requests.RequestOptions{
				Header: map[string][]string{
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
				Body: []byte("setAction=studentCourseSysSchedule&selectAction=TeachID&key1=" + courseID),
			},
		)
		if err != nil {
			log.Println("[WARNING]", err)
			continue
		}
		break
	}
	if config.Client.Verbose {
		log.Println("[DEBUG]", string(resp.Content))
	}
	if strings.Contains(string(resp.Content), "当前选课系统不可用") {
		log.Println("[WARNING] System is not yet available. Retrying in", config.Client.Delay, "milliseconds...")
		return ""
	}
	matches := courseRe.FindStringSubmatch(string(resp.Content))
	if len(matches) == 0 {
		log.Printf("[WARNING] Course %s not found. Skipping...\n", courseID)
		return "0"
	}
	return matches[1]
}

func addCourse(courseID, teachID string) {
	var resp *requests.Response
	var err error
	for {
		resp, err = session.Get(config.Session.BaseURI+fmt.Sprintf("/vatuu/CourseStudentAction?setAction=addStudentCourseApply&teachId=%s&isBook=0&tt=%d", teachID, time.Now().Unix()), nil)
		if err != nil {
			log.Println("[WARNING]", err)
			continue
		}
		break
	}
	text, _ := gbkDecoder.String(string(resp.Content))
	if config.Client.Verbose {
		log.Println("[DEBUG]", text)
	}
	var state StateType
	if strings.Contains(text, "成功") || strings.Contains(text, courseID+"冲突") || strings.Contains(text, "已选该课程") {
		state = Success
	} else if strings.Contains(text, "冲突") {
		state = Conflict
	} else if strings.Contains(text, "上限") {
		state = Overflowed
	}
	resCh <- &Result{CourseID: courseID, TeachID: teachID, State: state}
}
