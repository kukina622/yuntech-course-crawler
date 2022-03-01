package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"yuntech-course-crawler/crawler"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}
	jar, _ := cookiejar.New(nil)
	job := cron.New()
	job.AddFunc("*/3 * * * *", func() {
		logrus.Info("start crawling")
		task(jar)
	})
	job.Start()
	select {}
}

func task(jar *cookiejar.Jar) {
	username := os.Getenv("username")
	password := os.Getenv("password")
	courseList := strings.Split(os.Getenv("course"), ",")
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
		if courseSearchCrawler.Course.MaxPeole > courseSearchCrawler.Course.NowNumberOfPeople {
			loginResult := yunTechSSOCrawler.Login()
			if !loginResult {
				logrus.Fatal("warn username or password")
			}
			addResult := courseRegisterCrawler.AddCourse(serialNo)
			if !addResult {
				logrus.Fatal("error! please check your curriculum")
			}
			logrus.Info(fmt.Sprintf("%s: add successful", serialNo))
		}
	}
}

func loadCourse() []string {
	allCourseList := strings.Split(os.Getenv("course"), ",")
	return allCourseList
}
