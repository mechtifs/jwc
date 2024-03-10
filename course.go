package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"jwc/requests"
)

type StateType int

const (
	Failed StateType = iota
	Chosen
	Conflict
	Overflowed
)

type Course struct {
	CourseID    string
	TeachID     string
	ListID      string
	IsAvailable bool
	State       StateType
}

var teachIDRe = regexp.MustCompile(`<span id="teachIdChoose.+?" style="display:none;">([0-9A-F]{16})<\/span>`)

func (c *Course) updateListID() {
	var resp *requests.Response
	var err error
	for {
		resp, err = session.Get(config.Session.BaseURI+"/vatuu/CourseStudentAction?setAction=studentCourseSysList&viewType=delCourse", nil)
		if err != nil {
			log.Println("[WARNING]", err)
			continue
		}
		break
	}
	if config.Client.Verbose {
		log.Println("[DEBUG]", string(resp.Content))
	}
	re := regexp.MustCompile(fmt.Sprintf(`'([0-9A-F]{32})','%s'`, c.CourseID))
	matches := re.FindStringSubmatch(string(resp.Content))
	if len(matches) == 0 {
		return
	}
	c.ListID = matches[1]
}

func (c *Course) updateTeachIDAvailability() {
	var resp *requests.Response
	var err error
	for {
		resp, err = session.Post(config.Session.BaseURI+"/vatuu/CourseStudentAction",
			&requests.RequestOptions{
				Header: map[string][]string{
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
				Body: []byte("setAction=studentCourseSysSchedule&selectAction=TeachID&key1=" + c.CourseID),
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
		return
	}
	matches := teachIDRe.FindStringSubmatch(string(resp.Content))
	if len(matches) == 0 {
		log.Printf("[WARNING] Course %s not found. Skipping...\n", c.CourseID)
		c.TeachID = "0"
	}
	c.TeachID = matches[1]
	c.IsAvailable = !strings.Contains(string(resp.Content), "已选满")
}
