package gh

import (
	"os"
	"testing"

	"github.com/aiteung/atdb"
)

var DBMongoinfo = atdb.DBInfo{
	DBString: os.Getenv("MONGOSTRINGAWANGGA"),
	DBName:   "github",
}

var config = Config{
	URL:       "https://github.com/settings/personal-access-tokens/new",
	UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36",
	RTO:       10,
	Host:      "github.com",
	Origin:    "https://github.com",
	MongoConn: atdb.MongoConnect(DBMongoinfo),
}

func TestUpdateCookies(m *testing.T) {
	UpdateCookies("iteungbot", config)
}
