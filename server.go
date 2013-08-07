// server and route 
package bew

import (
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strconv"
)

type route struct {
	r       string
	method  string
	handler reflect.Value
}

type Server struct {
	l      net.Listener
	routes []route
}

func NewServer() (s *Server) {
	s = &Server{}
	return
}

func (s *Server) Run(bind string) {
	mux := http.NewServeMux()

	mux.Handle("/", s)

	l, err := net.Listen("tcp", bind)
	if err != nil {
		fmt.Println("bind " + bind + " error")
	}

	s.l = l
	err = http.Serve(s.l, mux)
}

func (s *Server) ServeHTTP(c http.ResponseWriter, r *http.Request) {
	s.route(c, r)
}


// Route related methods
func (s *Server) addRoute(r string, method string, handler interface{}) {
	s.routes = append(s.routes, route{r: r, method: method, handler: reflect.ValueOf(handler)})
}

func (s *Server) Get(r string, handler interface{}) {
	s.addRoute(r, "GET", handler)
}

func (s *Server) Post(r string, handler interface{}) {
	s.addRoute(r, "POST", handler)
}

func (s *Server) route(c http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path

	for _, route := range s.routes {
		if route.r != requestPath {
			continue
		} else {
			var args []reflect.Value
			args = append(args, reflect.ValueOf(r))
			ret := route.handler.Call(args)
			content := []byte(ret[0].String())

			c.Header().Set("Content-Length", strconv.Itoa(len(content)))
			
			c.Write(content)
		}
	}
}













