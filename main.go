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

var total int
var session *requests.Session
var resCh chan *Result
var gbkDecoder = simplifiedchinese.GBK.NewDecoder()
var courseRe = regexp.MustCompile(`<span id="teachIdChoose.+?" style="display:none;">([0-9A-F]{16})<\/span>`)
var config Config

func showSummary(courses []string, retries int) {
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

func main() {
	parseConfig()
	courseIDs := config.Session.Targets
	total = len(courseIDs)
	log.Printf("[NOTICE] Choosing %d course(s)\n", total)
	session = requests.NewSession(&requests.SessionOptions{
		Header: map[string][]string{
			"User-Agent": {config.Session.UserAgent},
			"Cookie":     {"JSESSIONID=" + config.Session.Credential},
		},
	})
	resCh = make(chan *Result, config.Client.Parallel*total)
	for i := 0; i < config.Client.Parallel; i++ {
		for _, courseID := range courseIDs {
			for {
				teachID := queryCourse(courseID)
				if teachID != "" {
					go addCourse(courseID, teachID)
					break
				}
				if teachID == "0" {
					break
				}
				time.Sleep(time.Duration(config.Client.Delay) * time.Millisecond)
			}
		}
	}
	retries := 0
	if config.Summary.Enabled {
		go func() {
			for {
				time.Sleep(time.Duration(config.Summary.Interval) * time.Second)
				showSummary(courseIDs, retries)
			}
		}()
	}
	for res := range resCh {
		if res.Success {
			if !config.Client.Keep {
				for i, courseID := range courseIDs {
					if courseID == res.CourseID {
						courseIDs = append(courseIDs[:i], courseIDs[i+1:]...)
						break
					}
				}
			}
			log.Println("[NOTICE]", res, "may be already chosen.")
			if len(courseIDs) == 0 {
				log.Println("[NOTICE] Nothing to do. Quitting...")
				break
			}
			continue
		}
		retries++
		time.Sleep(time.Duration(config.Client.Delay) * time.Millisecond)
		go addCourse(res.CourseID, res.TeachID)
	}
}
