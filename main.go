package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"yuntech-course-crawler/crawler"
	"yuntech-course-crawler/utils"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}
	logger := cron.VerbosePrintfLogger(log.New(os.Stdout, "", log.LstdFlags))
	jar, _ := cookiejar.New(nil)
	job := cron.New(cron.WithChain(cron.Recover(logger)))
	job.AddFunc("*/1 * * * *", func() {
		result := task(jar)
		if result {
			job.Stop()
			logrus.Info("job stop")
			os.Exit(0)
		}
	})
	job.Start()
	select {}
}

func task(jar *cookiejar.Jar) bool {
	username := os.Getenv("username")
	password := os.Getenv("password")
	courseList := loadCourse()
	if len(courseList) == 0 {
		return true
	}
	logrus.Info("start crawling")
	courseSearchCrawler := crawler.CourseSearchCrawler{}
	yunTechSSOCrawler := crawler.YunTechSSOCrawler{
		Username: username,
		Password: password,
		Client:   &http.Client{Jar: jar},
	}
	courseRegisterCrawler := crawler.CourseRegisterCrawler{
		YunTechSSOCrawler: yunTechSSOCrawler,
	}
	for _, serialNo := range courseList {
		courseSearchCrawler.QueryCourse(serialNo)
		if courseSearchCrawler.Course.MaxPeople > courseSearchCrawler.Course.NowNumberOfPeople {
			logrus.Info(fmt.Sprintf("%s: has vacancies", serialNo))
			loginResult := yunTechSSOCrawler.Login()
			if !loginResult {
				panic("warn username or password")
			}
			addResult := courseRegisterCrawler.AddCourse(serialNo)
			if !addResult {
				panic("error! please check your curriculum")
			}
			logrus.Info(fmt.Sprintf("%s: added successfully", serialNo))
			addComplete(serialNo)
		}else{
			logrus.Info(fmt.Sprintf("%s: is full", serialNo))
		}
	}
	return false
}

func loadCourse() []string {
	allCourseList := strings.Split(os.Getenv("course"), ",")
	_, err := os.Stat("complete.tmp")
	// file existed
	if err == nil {
		completeCourseBytes, _ := ioutil.ReadFile("complete.tmp")
		completeCourse := strings.Split(string(completeCourseBytes), ",")
		for _, j := range completeCourse[:len(completeCourse)-1] {
			index := utils.FindIndex(allCourseList, j)
			if index != -1 {
				allCourseList = utils.RemoveIndex(allCourseList, index)
			}
		}
	}
	return allCourseList
}

func addComplete(serialNo string) {
	f, err := os.OpenFile("complete.tmp", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(fmt.Sprintf("%s,", serialNo)); err != nil {
		panic(err)
	}
}
