package pkg

import (
	"net/http"
)

//cookie 默认值
var cookie = &http.Cookie{
	Path:     "/",
	Domain:   "127.0.0.1",
	MaxAge:   7 * 86400,
	HttpOnly: true,
}

//写cookie
func WriteCookie(rw http.ResponseWriter, name string, val string) {
	http.SetCookie(rw, &http.Cookie{Name: name, Value: val, Path: cookie.Path,
		Domain: cookie.Domain, MaxAge: cookie.MaxAge, HttpOnly: cookie.HttpOnly})
}

//删除cookie
func DelCookie(rw http.ResponseWriter, name string) {
	http.SetCookie(rw, &http.Cookie{Name: name, MaxAge: -1})
}

//读cookie
func ReadCookie(req *http.Request, name string) string {
	mycookie, err := req.Cookie(name)
	if err != nil {
		return ""
	}
	return mycookie.Value
}

func SetCookieDefVal(domain string, path string, maxAge int, httpOnly bool) {
	cookie.Domain = domain
	cookie.Path = path
	cookie.MaxAge = maxAge * 86400
	cookie.HttpOnly = httpOnly
}
