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

func (ctx *Context) NotFound(){
	ctx.ResponseWriter.WriteHeader(404)
	ctx.ResponseWriter.Write([]byte("Not Found"))
}

func (ctx *Context) BadRequest(){
	ctx.ResponseWriter.WriteHeader(400)
	ctx.ResponseWriter.Write([]byte("Bad Request"))
}

func Get(route string, handler interface{}){
	server.Get(route, handler)
}
func Post(route string, handler interface{}){
	server.Post(route, handler)
}
func Put(route string, handler interface{}){
	server.Put(route, handler)
}
func Delete(route string, handler interface{}){
	server.Delete(route, handler)
}

func Run(bind string){
	server.Run(bind)
}

var server = NewServer()













