package crawler

import (
	"github.com/anaskhan96/soup"
	"log"
	"net/url"
	"regexp"
	"strconv"
)

const COURSE_URL string = "https://webapp.yuntech.edu.tw/webnewcas/course/querycour.aspx"

type course struct {
	Name              string
	MaxPeople         int
	NowNumberOfPeople int
}

type CourseSearchCrawler struct {
	Course course
}

func (crawler *CourseSearchCrawler) QueryCourse(serialNo string) {
	resp, err := soup.Get(COURSE_URL)
	if err != nil {
		log.Fatal(err)
	}
	doc := soup.HTMLParse(resp)
	payload := url.Values{}
	payload.Add("__VIEWSTATE", doc.Find("input", "id", "__VIEWSTATE").Attrs()["value"])
	payload.Add("ctl00_MainContent_ToolkitScriptManager1_HiddenField", doc.Find("input", "id", "ctl00_MainContent_ToolkitScriptManager1_HiddenField").Attrs()["value"])
	payload.Add("__VIEWSTATEGENERATOR", doc.Find("input", "id", "__VIEWSTATEGENERATOR").Attrs()["value"])
	payload.Add("__EVENTVALIDATION", doc.Find("input", "id", "__EVENTVALIDATION").Attrs()["value"])
	payload.Add("ctl00$MainContent$AcadSeme", doc.Find("select", "id", "ctl00_MainContent_AcadSeme").Find("option", "selected", "selected").Attrs()["value"])
	payload.Add("ctl00$MainContent$CurrentSubj", serialNo)
	// fixed parameter
	payload.Add("ctl00$MainContent$Submit", "執行查詢")
	payload.Add("__EVENTTARGET", "")
	payload.Add("__EVENTARGUMENT", "")
	payload.Add("__VIEWSTATEENCRYPTED", "")
	payload.Add("ctl00$MainContent$College", "")
	payload.Add("ctl00$MainContent$DeptCode", "")
	payload.Add("ctl00$MainContent$TextBoxWatermarkExtender3_ClientState", "")
	payload.Add("ctl00$MainContent$SubjName", "")
	payload.Add("ctl00$MainContent$TextBoxWatermarkExtender1_ClientState", "")
	payload.Add("ctl00$MainContent$Instructor", "")
	payload.Add("ctl00$MainContent$TextBoxWatermarkExtender2_ClientState", "")
	resp, err = soup.PostForm(COURSE_URL, payload)
	if err != nil {
		log.Fatal(err)
	}
	course := crawler.parseCoursePage(resp)
	crawler.Course = course
}

func (crawler *CourseSearchCrawler) parseCoursePage(rawPage string) course {
	course := course{}
	doc := soup.HTMLParse(rawPage)
	courseRow := doc.Find("tr", "class", "GridView_Row").FindAll("td")
	// course name
	course.Name = courseRow[2].Find("a").Text()
	// NowNumberOfPeople
	_nowNumberOfPeople, _ := strconv.Atoi(courseRow[9].Find("span").Text())
	course.NowNumberOfPeople = _nowNumberOfPeople
	// max people
	r, _ := regexp.Compile("[0-9]+")
	strMax := r.FindString(courseRow[10].Find("span").FullText())
	_maxPeople, err := strconv.Atoi(strMax)
	if err != nil {
		_maxPeople = 9999999999
	}
	course.MaxPeople = _maxPeople
	return course
}
