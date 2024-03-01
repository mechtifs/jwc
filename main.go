package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"time"

	"jwc/requests"

	"github.com/BurntSushi/toml"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var session *requests.Session
var resCh chan *Result
var gbkDecoder = simplifiedchinese.GBK.NewDecoder()
var courseRe = regexp.MustCompile(`<span id="teachIdChoose.+?" style="display:none;">([0-9A-F]{16})<\/span>`)
var config Config

func showSummary(courses []string, total, retries int) {
	fmt.Printf(" *** Course Selection Summary as of %s ***\n", time.Now().Format("Mon Jan _2 15:04:05 MST 2006"))
	fmt.Println("================================================================================")
	fmt.Printf("[%s %d/%d(%.2f%%) PAR:%d RE:%d]\n", config.Session.Credential, total-len(courses), total, float32(total-len(courses))/float32(total), runtime.NumGoroutine()-2, retries)
	fmt.Println("PENDING:", courses)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println()
}

func parseConfig() {
	if len(os.Args) == 1 {
		log.Fatal("[FATAL] Please specify the path of config file.")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("[FATAL] Failed to open file:", err)
	}
	defer file.Close()
	toml.NewDecoder(file).Decode(&config)
	if config.Client.Verbose {
		log.Println("[DEBUG]", config)
	}
}

func getIDs() ([]string, []string) {
	var courseIDs, teachIDs []string
	for _, target := range config.Session.Targets {
		var teachID string
		for {
			teachID = queryCourse(target)
			if teachID != "" {
				break
			}
			time.Sleep(time.Duration(config.Client.Delay) * time.Millisecond)
		}
		if teachID == "0" {
			continue
		}
		courseIDs = append(courseIDs, target)
		teachIDs = append(teachIDs, teachID)
	}
	return courseIDs, teachIDs
}

func main() {
	parseConfig()
	log.Printf("[NOTICE] Choosing %d course(s)\n", len(config.Session.Targets))
	session = requests.NewSession(&requests.SessionOptions{
		Header: map[string][]string{
			"User-Agent": {config.Session.UserAgent},
			"Cookie":     {"JSESSIONID=" + config.Session.Credential},
		},
	})

	courseIDs, teachIDs := getIDs()

	total := len(courseIDs)
	if total == 0 {
		log.Println("[NOTICE] Nothing to do. Quitting...")
		return
	}

	for range config.Client.Parallel {
		for i, courseID := range courseIDs {
			teachID := teachIDs[i]
			go addCourse(courseID, teachID)
		}
	}

	retries := 0
	if config.Summary.Enabled && config.Summary.Interval > 0 {
		go func() {
			for {
				time.Sleep(time.Duration(config.Summary.Interval) * time.Second)
				showSummary(courseIDs, total, retries)
			}
		}()
	}

	resCh = make(chan *Result, config.Client.Parallel*total)
	for res := range resCh {
		if res.State == Success || res.State == Conflict {
			if !config.Client.Keep {
				for i, courseID := range courseIDs {
					if courseID == res.CourseID {
						courseIDs = append(courseIDs[:i], courseIDs[i+1:]...)
						break
					}
				}
			}
			if res.State == Success {
				log.Println("[NOTICE]", res.CourseID, "is already chosen.")
			} else {
				log.Println("[WARNING]", res.CourseID, "is in conflict with another course.")
				total--
			}
			if len(courseIDs) == 0 {
				log.Println("[NOTICE] Nothing to do. Quitting...")
				break
			}
		} else if res.State == Overflowed {
			log.Println("[FATAL] Total credits overflowed. Quitting...")
			break
		} else {
			retries++
			time.Sleep(time.Duration(config.Client.Delay) * time.Millisecond)
			go addCourse(res.CourseID, res.TeachID)
		}
	}
}
