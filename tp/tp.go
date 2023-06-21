package tp

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aiteung/atdb"
	"go.mongodb.org/mongo-driver/bson"
)

func Requests(method string, config Config, isPostBelanja bool, body io.Reader, cookies *Cookies) (respBody []byte) {
	client := http.Client{
		Timeout: time.Duration(config.RTO) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, config.URL, body)
	if err != nil {
		fmt.Print("http.NewRequest Got error ", err.Error())
	}
	req.AddCookie(&cookies.Cookiesession1)
	req.AddCookie(&cookies.Potp_session)
	req.AddCookie(&cookies.Csrf_cookie_name)
	if isPostBelanja {
		req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")

	} else {
		req.Header.Set("Accept", "*/*")
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", config.UserAgent)
	if isPostBelanja {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	} else {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	req.Header.Add("Host", config.Host)
	req.Header.Add("Origin", config.Origin)
	req.Header.Add("Referer", config.Referer)
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"111\", \"Not(A:Brand\";v=\"8\", \"Chromium\";v=\"111\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	if isPostBelanja {
		req.Header.Add("X-Requested-With", "XMLHttpRequest")
	}

	resp, err := client.Do(req)
	if err != nil { // loop disini buruk karena respon pun jadi error kalo di ulang
		fmt.Print("RTO: ", config.RTO, " , ", time.Now().Format("15:04:05"), err.Error())
	} else {
		//fmt.Println("Cookies response : ", resp.Cookies())
		for _, c := range resp.Cookies() {
			if c.Name == "potp_session" {
				cookies.Potp_session.Name = c.Name
				cookies.Potp_session.Value = c.Value
			} else if c.Name == "csrf_cookie_name" {
				cookies.Csrf_cookie_name.Name = c.Name
				cookies.Csrf_cookie_name.Value = c.Value
			} else if c.Name == "cookiesession1" {
				cookies.Cookiesession1.Name = c.Name
				cookies.Cookiesession1.Value = c.Value
			}
		}
		respStatusCode := resp.StatusCode
		fmt.Print("Status Code Response : ", respStatusCode)
		if respStatusCode == 200 {
			respBody, err = io.ReadAll(resp.Body)
			if err != nil {
				fmt.Print("could not read response body: ", err)
			}
			//fmt.Println("Body response : ", string(respBody))
		}
		defer resp.Body.Close()

	}
	return
}

func GetLoginInfo(cookies *Cookies, config Config) (user string) {
	body := Requests("GET", config, false, nil, cookies)
	for string(body) == "" {
		body = Requests("GET", config, false, nil, cookies)
	}
	//fmt.Println("GetLoginInfo Cookies: ", cookies)
	a := strings.Split(string(body), `#4F4F4F;`)
	if len(a) > 2 {
		b := strings.Split(a[2], "</h5>")
		c := strings.Split(b[0], ">")
		user = strings.TrimSpace(c[1])
	}

	return
}

func SetCookies(user UserSession) (usercookies Cookies) {
	var cookiesession1 http.Cookie
	var potp_session http.Cookie
	var csrf_cookie_name http.Cookie
	cookiesession1 = http.Cookie{
		Name:     "cookiesession1",
		Value:    user.Cookiesession1,
		Path:     "/",
		Expires:  time.Now().Local().Add(time.Hour * time.Duration(8760)),
		HttpOnly: true,
		Secure:   false,
	}
	potp_session = http.Cookie{
		Name:     "potp_session",
		Value:    user.Potp_session,
		Path:     "/",
		Expires:  time.Now().Local().Add(time.Hour * time.Duration(2)),
		HttpOnly: true,
		Secure:   true,
	}
	csrf_cookie_name = http.Cookie{
		Name:     "csrf_cookie_name",
		Value:    user.Csrf_cookie_name,
		Path:     "/",
		Expires:  time.Now().Local().Add(time.Hour * time.Duration(2)),
		HttpOnly: true,
		Secure:   true,
	}
	usercookies = Cookies{
		Cookiesession1:   cookiesession1,
		Potp_session:     potp_session,
		Csrf_cookie_name: csrf_cookie_name,
	}

	return
}

func UpdateCookies(email string, config Config) string {
	filter := bson.M{"username": email}
	usercred := atdb.GetOneDoc[UserCred](config.MongoConn, "user", filter)
	POTPCookies := SetCookies(usercred.Session)
	user := GetLoginInfo(&POTPCookies, config)
	fmt.Println("User Full Name: ", user)
	usercred.Session.Cookiesession1 = POTPCookies.Cookiesession1.Value
	usercred.Session.Potp_session = POTPCookies.Potp_session.Value
	usercred.Session.Csrf_cookie_name = POTPCookies.Csrf_cookie_name.Value
	usercred.Nama = user
	atdb.ReplaceOneDoc(config.MongoConn, "user", filter, usercred)
	return user
}
