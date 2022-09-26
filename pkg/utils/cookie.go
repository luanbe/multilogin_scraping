package util

import (
	"github.com/tebeka/selenium"
	"net/http"
	"time"
)

func ConvertSeleniumToHttpCookies(seleniumCookies []selenium.Cookie) []*http.Cookie {
	var httpCookies []*http.Cookie
	for _, cookie := range seleniumCookies {

		httpCookie := &http.Cookie{
			Name:    cookie.Name,
			Value:   cookie.Value,
			Path:    cookie.Path,
			Domain:  cookie.Domain,
			Secure:  cookie.Secure,
			Expires: time.Unix(int64(cookie.Expiry), 0),
		}
		httpCookies = append(httpCookies, httpCookie)
	}
	return httpCookies
}
