package gh

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	URL       string
	UserAgent string
	RTO       int
	Host      string
	Origin    string
	Referer   string
	MongoConn *mongo.Database
}

type Cookies struct {
	User_session                http.Cookie `json:"user_session,omitempty" bson:"user_session,omitempty"`
	Host_user_session_same_site http.Cookie `json:"__Host-user_session_same_site,omitempty" bson:"__Host-user_session_same_site,omitempty"`
	Gh_sess                     http.Cookie `json:"_gh_sess,omitempty" bson:"_gh_sess,omitempty"`
	Tz                          http.Cookie `json:"tz,omitempty" bson:"tz,omitempty"`
	Dotcom_user                 http.Cookie `json:"dotcom_user,omitempty" bson:"dotcom_user,omitempty"`
	Logged_in                   http.Cookie `json:"logged_in,omitempty" bson:"logged_in,omitempty"`
	Has_recent_activity         http.Cookie `json:"has_recent_activity,omitempty" bson:"has_recent_activity,omitempty"`
	Device_id                   http.Cookie `json:"_device_id,omitempty" bson:"_device_id,omitempty"`
	Preferred_color_mode        http.Cookie `json:"preferred_color_mode,omitempty" bson:"preferred_color_mode,omitempty"`
	Color_mode                  http.Cookie `json:"color_mode,omitempty" bson:"color_mode,omitempty"`
	Octo                        http.Cookie `json:"_octo,omitempty" bson:"_octo,omitempty"`
}

type UserSession struct {
	User_session                string `json:"user_session,omitempty" bson:"user_session,omitempty"`
	Host_user_session_same_site string `json:"__Host-user_session_same_site,omitempty" bson:"__Host-user_session_same_site,omitempty"`
	Gh_sess                     string `json:"_gh_sess,omitempty" bson:"_gh_sess,omitempty"`
	Tz                          string `json:"tz,omitempty" bson:"tz,omitempty"`
	Dotcom_user                 string `json:"dotcom_user,omitempty" bson:"dotcom_user,omitempty"`
	Logged_in                   string `json:"logged_in,omitempty" bson:"logged_in,omitempty"`
	Has_recent_activity         string `json:"has_recent_activity,omitempty" bson:"has_recent_activity,omitempty"`
	Device_id                   string `json:"_device_id,omitempty" bson:"_device_id,omitempty"`
	Preferred_color_mode        string `json:"preferred_color_mode,omitempty" bson:"preferred_color_mode,omitempty"`
	Color_mode                  string `json:"color_mode,omitempty" bson:"color_mode,omitempty"`
	Octo                        string `json:"_octo,omitempty" bson:"_octo,omitempty"`
}

type UserCred struct {
	Username string      `json:"username,omitempty" url:"username,omitempty"`
	Password string      `json:"password,omitempty" url:"password,omitempty"`
	Nama     string      `json:"nama,omitempty" url:"nama,omitempty"`
	Session  UserSession `json:"session,omitempty" url:"session,omitempty"`
}

type UserCookies struct {
	Username string      `json:"username,omitempty" url:"username,omitempty"`
	Password string      `json:"password,omitempty" url:"password,omitempty"`
	Cookies  interface{} `json:"cookies,omitempty" url:"cookies,omitempty"`
}
