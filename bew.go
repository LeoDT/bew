package bew

import (
	"net/http"
)

type Context struct{
	Request *http.Request
	Server *Server
	http.ResponseWriter
}

func (ctx *Context) Abort(status int, body string){
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write([]byte(body))
}

func (ctx *Context) Redirect(status int, url string){
	ctx.ResponseWriter.Header().Set("Location", url)
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write([]byte("Response to: " + url))
}

func Get(route string, handler interface{}){
	server.Get(route, handler)
}

func Run(bind string){
	server.Run(bind)
}

var server = NewServer()













