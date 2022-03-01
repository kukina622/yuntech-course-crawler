package crawler

import (
	// "fmt"
	"github.com/anaskhan96/soup"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	LOGIN_URL       = "https://webapp.yuntech.edu.tw/YunTechSSO/Account/Login"
	CHECK_LOGIN_URL = "https://webapp.yuntech.edu.tw/YunTechSSO/Account/IsLogined"
)

type YunTechSSOCrawler struct {
	Username string
	Password string
	Client   *http.Client
}

func (crawler *YunTechSSOCrawler) Login() bool {
	if crawler.checkLogin() {
		return true
	}
	token := crawler.getLoginToken()
	payload := url.Values{}
	payload.Add("__RequestVerificationToken", token)
	payload.Add("pLoginName", crawler.Username)
	payload.Add("pLoginPassword", crawler.Password)
	payload.Add("pRememberMe", "true")

	resp, err := crawler.Client.PostForm(LOGIN_URL, payload)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return crawler.checkLogin()
}

func (crawler *YunTechSSOCrawler) getLoginToken() string {
	resp, err := crawler.Client.Get(LOGIN_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	doc := soup.HTMLParse(string(body))
	token := doc.Find("input", "name", "__RequestVerificationToken").Attrs()["value"]
	return token
}

func (crawler *YunTechSSOCrawler) checkLogin() bool {
	resp, err := crawler.Client.Get(CHECK_LOGIN_URL)
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body) == "True"
}