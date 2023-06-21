# atcookies

How to use :

```go
package main

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
 URL:       "<https://github.com/settings/personal-access-tokens/new>",
 UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36",
 RTO:       10,
 Host:      "github.com",
 Origin:    "<https://github.com>",
 MongoConn: atdb.MongoConnect(DBMongoinfo),
}

func TestInitCookies(m*testing.T) {
 InitCookies("iteungbot", os.Getenv("ITEUNGGITHUBSESSION"), config)
}

func TestUpdateCookies(m *testing.T) {
 UpdateCookies("iteungbot", config)
}
```

## Development

manage cookies

```sh
git tag v0.0.1
git push origin --tags
go list -m github.com/aiteung/atcookies@v0.0.1
```
