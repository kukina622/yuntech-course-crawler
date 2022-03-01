package crawler

import (
	"github.com/anaskhan96/soup"
	"io/ioutil"
	"log"
	"net/url"
)

type CourseRegisterCrawler struct {
	YunTechSSOCrawler
}

const ADD_COURSE_URL = "https://webapp.yuntech.edu.tw/AAXCCS/CourseSelectionRegister.aspx"

func (crawler *CourseRegisterCrawler) AddCourse(SerialNo string) bool {
	body := crawler.search(SerialNo)
	body = crawler.register(body)
	body = crawler.nextStep(body)
	result := crawler.save(body)
	return result
}

func (crawler *CourseRegisterCrawler) search(SerialNo string) string {
	fres, err := crawler.Client.Get(ADD_COURSE_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer fres.Body.Close()
	fbody, _ := ioutil.ReadAll(fres.Body)
	fdoc := soup.HTMLParse(string(fbody))
	payload := url.Values{}
	payload.Add("__EVENTTARGET", "ctl00$ContentPlaceHolder1$QueryButton")
	payload.Add("__EVENTARGUMENT", "")
	payload.Add("__LASTFOCUS", "")
	payload.Add("__VIEWSTATE", fdoc.Find("input", "name", "__VIEWSTATE").Attrs()["value"])
	payload.Add("__VIEWSTATEGENERATOR", fdoc.Find("input", "name", "__VIEWSTATEGENERATOR").Attrs()["value"])
	payload.Add("__VIEWSTATEENCRYPTED", "")
	payload.Add("__EVENTVALIDATION", fdoc.Find("input", "name", "__EVENTVALIDATION").Attrs()["value"])
	payload.Add("ctl00$PlaceHolderMultipleLanguage$SelectLang", "")
	payload.Add("ctl00$ContentPlaceHolder1$CollegeDDL", "")
	payload.Add("ctl00$ContentPlaceHolder1$DayNightDDL", "1")
	payload.Add("ctl00$ContentPlaceHolder1$EduSysDDL", "B")
	payload.Add("ctl00$ContentPlaceHolder1$MajOpDDL", "")
	payload.Add("ctl00$ContentPlaceHolder1$CurrentSubjTextBox", SerialNo)
	payload.Add("ctl00$ContentPlaceHolder1$CourseNameTextBox", "")
	payload.Add("ctl00$ContentPlaceHolder1$TeacherNameTextBox", "")
	payload.Add("ctl00$ContentPlaceHolder1$CourProgramDDL", "")
	sres, err := crawler.Client.PostForm(ADD_COURSE_URL, payload)
	if err != nil {
		log.Fatal(err)
	}
	defer sres.Body.Close()
	sbody, _ := ioutil.ReadAll(sres.Body)
	return string(sbody)
}

func (crawler *CourseRegisterCrawler) register(body string) string {
	doc := soup.HTMLParse(body)
	checkBox := doc.Find("table", "id", "ContentPlaceHolder1_QueryCourseGridView").Find("input", "type", "checkbox").Attrs()["name"]
	payload := url.Values{}
	payload.Add("__EVENTTARGET", "ctl00$ContentPlaceHolder1$RegisterButton")
	payload.Add("__EVENTARGUMENT", "")
	payload.Add("__LASTFOCUS", "")
	payload.Add("__VIEWSTATE", doc.Find("input", "name", "__VIEWSTATE").Attrs()["value"])
	payload.Add("__VIEWSTATEGENERATOR", doc.Find("input", "name", "__VIEWSTATEGENERATOR").Attrs()["value"])
	payload.Add("__VIEWSTATEENCRYPTED", "")
	payload.Add("__EVENTVALIDATION", doc.Find("input", "name", "__EVENTVALIDATION").Attrs()["value"])
	payload.Add("ctl00$PlaceHolderMultipleLanguage$SelectLang", "")
	payload.Add("ctl00$ContentPlaceHolder1$CollegeDDL", "")
	payload.Add("ctl00$ContentPlaceHolder1$DayNightDDL", "1")
	payload.Add("ctl00$ContentPlaceHolder1$EduSysDDL", "B")
	payload.Add("ctl00$ContentPlaceHolder1$MajOpDDL", "")
	payload.Add("ctl00$ContentPlaceHolder1$CurrentSubjTextBox", "")
	payload.Add("ctl00$ContentPlaceHolder1$CourseNameTextBox", "")
	payload.Add("ctl00$ContentPlaceHolder1$TeacherNameTextBox", "")
	payload.Add("ctl00$ContentPlaceHolder1$CourProgramDDL", "")
	payload.Add(checkBox, "on")
	resp, err := crawler.Client.PostForm(ADD_COURSE_URL, payload)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	_body, _ := ioutil.ReadAll(resp.Body)
	return string(_body)
}

func (crawler *CourseRegisterCrawler) nextStep(body string) string {
	doc := soup.HTMLParse(body)
	payload := url.Values{}
	payload.Add("__EVENTTARGET", "ctl00$ContentPlaceHolder1$NextStepButton")
	payload.Add("__EVENTARGUMENT", "")
	payload.Add("__LASTFOCUS", "")
	payload.Add("__VIEWSTATE", doc.Find("input", "name", "__VIEWSTATE").Attrs()["value"])
	payload.Add("__VIEWSTATEGENERATOR", doc.Find("input", "name", "__VIEWSTATEGENERATOR").Attrs()["value"])
	payload.Add("__VIEWSTATEENCRYPTED", "")
	payload.Add("__EVENTVALIDATION", doc.Find("input", "name", "__EVENTVALIDATION").Attrs()["value"])
	payload.Add("ctl00$PlaceHolderMultipleLanguage$SelectLang", "")
	payload.Add("ctl00$ContentPlaceHolder1$CollegeDDL", "")
	payload.Add("ctl00$ContentPlaceHolder1$DayNightDDL", "1")
	payload.Add("ctl00$ContentPlaceHolder1$EduSysDDL", "B")
	payload.Add("ctl00$ContentPlaceHolder1$MajOpDDL", "")
	payload.Add("ctl00$ContentPlaceHolder1$CurrentSubjTextBox", "")
	payload.Add("ctl00$ContentPlaceHolder1$CourseNameTextBox", "")
	payload.Add("ctl00$ContentPlaceHolder1$TeacherNameTextBox", "")
	payload.Add("ctl00$ContentPlaceHolder1$CourProgramDDL", "")
	resp, err := crawler.Client.PostForm(ADD_COURSE_URL, payload)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	_body, _ := ioutil.ReadAll(resp.Body)
	return string(_body)
}

func (crawler *CourseRegisterCrawler) save(body string) bool {
	doc := soup.HTMLParse(body)
	payload := url.Values{}
	payload.Add("__EVENTTARGET", "ctl00$ContentPlaceHolder1$SaveButton")
	payload.Add("__EVENTARGUMENT", "")
	payload.Add("__VIEWSTATE", doc.Find("input", "name", "__VIEWSTATE").Attrs()["value"])
	payload.Add("__VIEWSTATEGENERATOR", doc.Find("input", "name", "__VIEWSTATEGENERATOR").Attrs()["value"])
	payload.Add("__VIEWSTATEENCRYPTED", "")
	payload.Add("__EVENTVALIDATION", doc.Find("input", "name", "__EVENTVALIDATION").Attrs()["value"])
	payload.Add("ctl00$PlaceHolderMultipleLanguage$SelectLang", "")
	resp, err := crawler.Client.PostForm(ADD_COURSE_URL, payload)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}
