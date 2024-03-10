package main

import (
	"fmt"
	"jwc/requests"
	"log"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

type JWSession struct {
	requests.Session
}

func NewJWSession(opt *requests.SessionOptions) *JWSession {
	return &JWSession{
		*requests.NewSession(opt),
	}
}

func (s *JWSession) addCourse(course *Course, resCh chan *Course) {
	var resp *requests.Response
	var err error
	for {
		resp, err = s.Get(config.Session.BaseURI+fmt.Sprintf("/vatuu/CourseStudentAction?setAction=addStudentCourseApply&teachId=%s&isBook=0&tt=%d", course.TeachID, time.Now().Unix()), nil)
		if err != nil {
			log.Println("[WARNING]", err)
			continue
		}
		break
	}
	gbkDecoder := simplifiedchinese.GBK.NewDecoder()
	text, _ := gbkDecoder.String(string(resp.Content))
	if config.Client.Verbose {
		log.Println("[DEBUG]", text)
	}
	var state StateType
	if strings.Contains(text, "成功") || strings.Contains(text, course.CourseID+"冲突") || strings.Contains(text, "已选该课程") {
		state = Chosen
	} else if strings.Contains(text, "冲突") {
		state = Conflict
	} else if strings.Contains(text, "上限") {
		state = Overflowed
	}
	course.State = state
	resCh <- course
}

// No matter what parameter you pass to this api, it will always return the same response...
func (s *JWSession) quitCourse(info *Course) {
	var resp *requests.Response
	var err error
	for {
		resp, err = s.Get(config.Session.BaseURI+fmt.Sprintf("/vatuu/CourseStudentAction?setAction=delStudentCourseList&listId=%s&teachId=%s&tt=%d", info.ListID, info.TeachID, time.Now().Unix()), nil)
		if err != nil {
			log.Println("[WARNING]", err)
			continue
		}
		break
	}
	if config.Client.Verbose {
		gbkDecoder := simplifiedchinese.GBK.NewDecoder()
		text, _ := gbkDecoder.String(string(resp.Content))
		log.Println("[DEBUG]", text)
	}
}
