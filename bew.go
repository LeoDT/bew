package bew

import (
	"net/http"
	"time"
)

type Context struct {
	Request *http.Request
	Server  *Server
	Params BaseDict
	Json JsonDict
	Files FileDict
	http.ResponseWriter
}

type Cookie struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	Expires  time.Time
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

func (c *Cookie) ToHttpCookie() (cookie *http.Cookie) {
	cookie = &http.Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Path:     c.Path,
		Domain:   c.Domain,
		Expires:  c.Expires,
		MaxAge:   c.MaxAge,
		Secure:   c.Secure,
		HttpOnly: c.HttpOnly,
	}

	return
}

func (ctx *Context) Init() {
	r := ctx.Request

	if r.Header.Get("Content-Type") == "application/json" {
		json := JsonDict{}

		json.Parse(r)
		ctx.Json = json
	} else {
		ctx.Params = r.URL.Query()
		r.ParseForm()
		for k, v := range r.Form {
			ctx.Params.Set(k, v[0])
		}
	}
}

func (ctx *Context) Abort(status int, body string) {
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write([]byte(body))
}

func (ctx *Context) Redirect(status int, url string) {
	ctx.ResponseWriter.Header().Set("Location", url)
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write([]byte("Response to: " + url))
}

func (ctx *Context) NotFound() {
	ctx.ResponseWriter.WriteHeader(404)
	ctx.ResponseWriter.Write([]byte("Not Found"))
}

func (ctx *Context) BadRequest() {
	ctx.ResponseWriter.WriteHeader(400)
	ctx.ResponseWriter.Write([]byte("Bad Request"))
}

func (ctx *Context) Cookie(name string) string {
	cookie, err := ctx.Request.Cookie(name)

	if err != nil {
		return ""
	}

	return cookie.Value
}

func (ctx *Context) SetCookie(cookie *Cookie) {
	http.SetCookie(ctx.ResponseWriter, cookie.ToHttpCookie())
}

func Get(route string, handler interface{}) {
	server.Get(route, handler)
}
func Post(route string, handler interface{}) {
	server.Post(route, handler)
}
func Put(route string, handler interface{}) {
	server.Put(route, handler)
}
func Delete(route string, handler interface{}) {
	server.Delete(route, handler)
}

func Run(bind string) {
	server.Run(bind)
}

var server = NewServer()
