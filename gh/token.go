package gh

import (
	"fmt"
	"strings"

	"github.com/whatsauth/watoken"
)

const WebKitForm = `------WebKitFormBoundary#16RANDOMSTRING#
Content-Disposition: form-data; name="authenticity_token"

#DATACSRF#
------WebKitFormBoundary#16RANDOMSTRING#
Content-Disposition: form-data; name="value"

#VALUE#
------WebKitFormBoundary#16RANDOMSTRING#--`

const WebKitFormContentType = `multipart/form-data; boundary=----WebKitFormBoundary#16RANDOMSTRING#`

const CheckNameURL = "https://github.com/settings/personal-access-tokens/check-name"

func GeneratePayloadCheckName(authenticity_token string, value string) (payload, contenttype string) {
	csrf := strings.ReplaceAll(WebKitForm, "#DATACSRF#", authenticity_token)
	val := strings.ReplaceAll(csrf, "#VALUE#", value)
	random := watoken.RandomString(16)
	payload = strings.ReplaceAll(val, "#16RANDOMSTRING#", random)
	contenttype = strings.ReplaceAll(WebKitFormContentType, "#16RANDOMSTRING#", random)
	return
}

func GetCsrfAutoCheckInfo(cookies *Cookies, isInit bool, config Config) (user string) {
	body := Requests("GET", config, isInit, false, nil, cookies)
	for string(body) == "" {
		body = Requests("GET", config, isInit, false, nil, cookies)
	}
	dtml := strings.Split(string(body), `<auto-check src="/settings/personal-access-tokens/check-name" required>`)
	if len(dtml) > 1 {
		trimkiri := strings.Split(dtml[1], `<input type="hidden" value="`)
		trimkanan := strings.Split(trimkiri[1], `" data-csrf="true" />`)
		user = strings.TrimSpace(trimkanan[0])
	}

	return
}

func SendCheckName(tokenname string, username string, config Config) string {
	usercred := GetUserFromDB(username, config)
	GHCookies := SetCookies(usercred.Session, false)
	token := GetCsrfAutoCheckInfo(&GHCookies, false, config)
	fmt.Println("\nauthenticity_token: ", token)
	var payload string
	payload, config.Content_type = GeneratePayloadCheckName(token, tokenname)
	config.URL = CheckNameURL
	fmt.Println(payload)
	fmt.Println(len(payload))
	fmt.Println(config)
	body := Requests("POST", config, false, true, strings.NewReader(payload), &GHCookies)
	UpdateCookieInDB(username, "", usercred, GHCookies, config)
	return string(body)
}
