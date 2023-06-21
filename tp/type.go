package tp

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
	Cookiesession1   http.Cookie `json:"cookiesession1,omitempty" bson:"cookiesession1,omitempty"`
	Potp_session     http.Cookie `json:"potp_session,omitempty" bson:"potp_session,omitempty"`
	Csrf_cookie_name http.Cookie `json:"csrf_cookie_name,omitempty" bson:"csrf_cookie_name,omitempty"`
}

type UserSession struct {
	Cookiesession1   string `json:"cookiesession1,omitempty" url:"cookiesession1,omitempty"`
	Potp_session     string `json:"potp_session,omitempty" url:"potp_session,omitempty"`
	Csrf_cookie_name string `json:"csrf_cookie_name,omitempty" url:"csrf_cookie_name,omitempty"`
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
