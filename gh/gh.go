package gh

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aiteung/atdb"
	"go.mongodb.org/mongo-driver/bson"
)

func Requests(method string, config Config, isInit bool, isPost bool, body io.Reader, cookies *Cookies) (respBody []byte) {
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
	req.AddCookie(&cookies.User_session)
	req.AddCookie(&cookies.Host_user_session_same_site)
	if !isInit {
		req.AddCookie(&cookies.Gh_sess)
		req.AddCookie(&cookies.Tz)
		req.AddCookie(&cookies.Dotcom_user)
		req.AddCookie(&cookies.Logged_in)
		req.AddCookie(&cookies.Has_recent_activity)
		req.AddCookie(&cookies.Device_id)
		req.AddCookie(&cookies.Preferred_color_mode)
		req.AddCookie(&cookies.Color_mode)
		req.AddCookie(&cookies.Octo)
	}
	req.Header.Set("Accept", "*/*")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", config.UserAgent)
	if isPost {
		req.Header.Add("Content-Type", config.Content_type)
	}

	req.Header.Add("Host", config.Host)
	req.Header.Add("Origin", config.Origin)
	req.Header.Add("Referer", config.Referer)

	resp, err := client.Do(req)
	if err != nil { // loop disini buruk karena respon pun jadi error kalo di ulang
		fmt.Print("RTO: ", config.RTO, " , ", time.Now().Format("15:04:05"), err.Error())
	} else {
		//fmt.Println("Cookies response : ", resp.Cookies())
		for _, c := range resp.Cookies() {
			if c.Name == "user_session" {
				cookies.User_session.Name = c.Name
				cookies.User_session.Value = c.Value
			} else if c.Name == "__Host-user_session_same_site" {
				cookies.Host_user_session_same_site.Name = c.Name
				cookies.Host_user_session_same_site.Value = c.Value
				cookies.Host_user_session_same_site.Secure = c.Secure
			} else if c.Name == "_gh_sess" {
				cookies.Gh_sess.Name = c.Name
				cookies.Gh_sess.Value = c.Value
			} else if c.Name == "tz" {
				cookies.Tz.Name = c.Name
				cookies.Tz.Value = c.Value
			} else if c.Name == "dotcom_user" {
				cookies.Dotcom_user.Name = c.Name
				cookies.Dotcom_user.Value = c.Value
			} else if c.Name == "logged_in" {
				cookies.Logged_in.Name = c.Name
				cookies.Logged_in.Value = c.Value
			} else if c.Name == "has_recent_activity" {
				cookies.Has_recent_activity.Name = c.Name
				cookies.Has_recent_activity.Value = c.Value
			} else if c.Name == "_device_id" {
				cookies.Device_id.Name = c.Name
				cookies.Device_id.Value = c.Value
			} else if c.Name == "preferred_color_mode" {
				cookies.Preferred_color_mode.Name = c.Name
				cookies.Preferred_color_mode.Value = c.Value
			} else if c.Name == "color_mode" {
				cookies.Color_mode.Name = c.Name
				cookies.Color_mode.Value = c.Value
			} else if c.Name == "_octo" {
				cookies.Octo.Name = c.Name
				cookies.Octo.Value = c.Value
			}
		}
		respStatusCode := resp.StatusCode
		fmt.Print("Status Code Response : ", respStatusCode)
		if respStatusCode == 200 {
			respBody, err = io.ReadAll(resp.Body)
			if err != nil {
				fmt.Print("could not read response body: ", err)
			}
			//fmt.Println("Body response : ", respBody)
		}
		defer resp.Body.Close()

	}
	return
}

func GetLoginInfo(cookies *Cookies, isInit bool, config Config) (user string) {
	body := Requests("GET", config, isInit, false, nil, cookies)
	for string(body) == "" {
		body = Requests("GET", config, isInit, false, nil, cookies)
	}
	dtml := strings.Split(string(body), `<span data-view-component="true" class="Truncate-text">`)
	if len(dtml) > 2 {
		name := strings.Split(dtml[2], "</span>")
		user = strings.TrimSpace(name[0])
	}

	return
}

func SetCookies(user UserSession, isInit bool) (usercookies Cookies) {
	var user_session http.Cookie
	var host_user_session_same_site http.Cookie
	user_session = http.Cookie{
		Name:  "user_session",
		Value: user.User_session,
		Path:  "/",
		//Expires:  time.Now().Local().Add(time.Hour * time.Duration(8760)),
		//HttpOnly: true,
		//Secure:   true,
	}
	host_user_session_same_site = http.Cookie{
		Name:  "__Host-user_session_same_site",
		Value: user.Host_user_session_same_site,
		Path:  "/",
		//Expires:  time.Now().Local().Add(time.Hour * time.Duration(2)),
		HttpOnly: true,
		Secure:   true,
	}
	usercookies = Cookies{
		User_session:                user_session,
		Host_user_session_same_site: host_user_session_same_site,
	}
	if !isInit {
		var gh_sess http.Cookie
		var tz http.Cookie
		var dotcom_user http.Cookie
		var logged_in http.Cookie
		var has_recent_activity http.Cookie
		var device_id http.Cookie
		var preferred_color_mode http.Cookie
		var color_mode http.Cookie
		var octo http.Cookie
		gh_sess = http.Cookie{
			Name:  "_gh_sess",
			Value: user.Gh_sess,
			Path:  "/",
			//Expires:  time.Now().Local().Add(time.Hour * time.Duration(2)),
			HttpOnly: true,
			Secure:   true,
		}
		tz = http.Cookie{
			Name:  "tz",
			Value: user.Tz,
			Path:  "/",
			//Expires:  time.Now().Local().Add(time.Hour * time.Duration(2)),
			HttpOnly: true,
			Secure:   true,
		}
		dotcom_user = http.Cookie{
			Name:  "dotcom_user",
			Value: user.Dotcom_user,
			Path:  "/",
			//Expires:  time.Now().Local().Add(time.Hour * time.Duration(2)),
			HttpOnly: true,
			Secure:   true,
		}
		logged_in = http.Cookie{
			Name:  "logged_in",
			Value: user.Logged_in,
			Path:  "/",
			//Expires:  time.Now().Local().Add(time.Hour * time.Duration(2)),
			HttpOnly: true,
			Secure:   true,
		}
		has_recent_activity = http.Cookie{
			Name:  "has_recent_activity",
			Value: user.Has_recent_activity,
			Path:  "/",
			//Expires:  time.Now().Local().Add(time.Hour * time.Duration(2)),
			HttpOnly: true,
			Secure:   true,
		}
		device_id = http.Cookie{
			Name:  "_device_id",
			Value: user.Device_id,
			Path:  "/",
			//Expires:  time.Now().Local().Add(time.Hour * time.Duration(2)),
			HttpOnly: true,
			Secure:   true,
		}
		preferred_color_mode = http.Cookie{
			Name:  "preferred_color_mode",
			Value: user.Preferred_color_mode,
			Path:  "/",
			//Expires:  time.Now().Local().Add(time.Hour * time.Duration(2)),
			HttpOnly: true,
			Secure:   true,
		}
		color_mode = http.Cookie{
			Name:  "color_mode",
			Value: user.Color_mode,
			Path:  "/",
			//Expires:  time.Now().Local().Add(time.Hour * time.Duration(2)),
			HttpOnly: true,
			Secure:   true,
		}
		octo = http.Cookie{
			Name:  "_octo",
			Value: user.Octo,
			Path:  "/",
			//Expires:  time.Now().Local().Add(time.Hour * time.Duration(2)),
			HttpOnly: true,
			Secure:   true,
		}
		usercookies = Cookies{
			User_session:                user_session,
			Host_user_session_same_site: host_user_session_same_site,
			Gh_sess:                     gh_sess,
			Tz:                          tz,
			Dotcom_user:                 dotcom_user,
			Logged_in:                   logged_in,
			Has_recent_activity:         has_recent_activity,
			Device_id:                   device_id,
			Preferred_color_mode:        preferred_color_mode,
			Color_mode:                  color_mode,
			Octo:                        octo,
		}
	}

	return
}

func UpdateCookies(username string, config Config) string {
	usercred := GetUserFromDB(username, config)
	GHCookies := SetCookies(usercred.Session, false)
	userfullname := GetLoginInfo(&GHCookies, false, config)
	fmt.Println("\nUser Full Name: ", userfullname)
	UpdateCookieInDB(username, userfullname, usercred, GHCookies, config)
	return userfullname
}

func InitCookies(username string, user_session string, config Config) string {
	var usercred UserCred
	usercred.Session.User_session = user_session
	usercred.Session.Host_user_session_same_site = user_session
	GHCookies := SetCookies(usercred.Session, true)
	userfullname := GetLoginInfo(&GHCookies, true, config)
	fmt.Println("\nUser Full Name: ", userfullname)
	CookieToSessionDB(&usercred, GHCookies, userfullname)
	atdb.InsertOneDoc(config.MongoConn, "user", usercred)
	return userfullname
}

func UpdateCookieInDB(username string, userfullname string, usercred UserCred, GHCookies Cookies, config Config) {
	CookieToSessionDB(&usercred, GHCookies, userfullname)
	filter := bson.M{"username": username}
	atdb.ReplaceOneDoc(config.MongoConn, "user", filter, usercred)
}

func GetUserFromDB(username string, config Config) UserCred {
	filter := bson.M{"username": username}
	usercred := atdb.GetOneDoc[UserCred](config.MongoConn, "user", filter)
	return usercred
}

func CookieToSessionDB(usercred *UserCred, GHCookies Cookies, userfullname string) {
	usercred.Session.User_session = GHCookies.User_session.Value
	usercred.Session.Host_user_session_same_site = GHCookies.Host_user_session_same_site.Value
	usercred.Session.Gh_sess = GHCookies.Gh_sess.Value
	usercred.Session.Tz = GHCookies.Tz.Value
	usercred.Session.Dotcom_user = GHCookies.Dotcom_user.Value
	usercred.Session.Logged_in = GHCookies.Logged_in.Value
	usercred.Session.Has_recent_activity = GHCookies.Has_recent_activity.Value
	usercred.Session.Device_id = GHCookies.Device_id.Value
	usercred.Session.Preferred_color_mode = GHCookies.Preferred_color_mode.Value
	usercred.Session.Color_mode = GHCookies.Color_mode.Value
	usercred.Session.Octo = GHCookies.Octo.Value
	if userfullname != "" {
		usercred.Nama = userfullname
	}

}
