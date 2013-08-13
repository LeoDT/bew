// server and route
package bew

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
)

type route struct {
	r       string
	regex   *regexp.Regexp
	method  string
	handler reflect.Value
}

type Server struct {
	l      net.Listener
	routes []route
}

var rules = map[string]*regexp.Regexp{
	"int":     regexp.MustCompile("<[^:>]+?:int>"),
	"path":    regexp.MustCompile("<[^:>]+?:path>"),
	"default": regexp.MustCompile("<[^:>]+?>"),
}
var rule_replacer = map[string]string{
	"int":     `(\d+)`,
	"path":    `(.+?)`,
	"default": `([^/]+?)`,
}

func (r *route) compileRouteRegex() {
	route_pattern_string := r.r

	for k, v := range rules {
		route_pattern_string = v.ReplaceAllString(route_pattern_string, rule_replacer[k])
	}

	route_pattern_regex := regexp.MustCompile("^" + route_pattern_string + "/?$")
	r.regex = route_pattern_regex
}

func (r *route) ParseParams(params []string) (parsed []interface{}) {
	type_regex := regexp.MustCompile("<[^:]+?:([^:]+?)>|<[^:]+?>")

	type_list := type_regex.FindAllStringSubmatch(r.r, -1)

	parsed = make([]interface{}, len(params))
	for i, t := range type_list {
		switch t[len(t)-1] {
		case "int":
			item, _ := strconv.Atoi(params[i])
			parsed[i] = item
		default:
			parsed[i] = params[i]
		}
	}

	return parsed
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
	if r.RequestURI == "*" {
		c.Header().Set("Connection", "close")
		c.WriteHeader(400)
		return
	}
	s.route(c, r)
}

// Route related methods
func (s *Server) addRoute(r string, method string, handler interface{}) {
	new_route := route{r: r, method: method, handler: reflect.ValueOf(handler)}
	new_route.compileRouteRegex()

	s.routes = append(s.routes, new_route)
}

func (s *Server) Get(r string, handler interface{}) {
	s.addRoute(r, "GET", handler)
}

func (s *Server) Post(r string, handler interface{}) {
	s.addRoute(r, "POST", handler)
}

func (s *Server) Put(r string, handler interface{}) {
	s.addRoute(r, "PUT", handler)
}

func (s *Server) Delete(r string, handler interface{}) {
	s.addRoute(r, "DELETE", handler)
}

func matchRoute(r route, path string) (match bool, result []interface{}) {
	match = r.regex.MatchString(path)

	if match {
		pattern := r.regex.FindAllStringSubmatch(path, -1)

		result = r.ParseParams(pattern[0][1:])
	}

	return
}

var contextType reflect.Type

func init() {
	contextType = reflect.TypeOf(Context{})
}

func requireContext(handler *reflect.Value) bool {
	handler_type := handler.Type()

	if handler_type.NumIn() == 0 {
		return false
	}

	arg0 := handler_type.In(0)
	if arg0.Kind() != reflect.Ptr {
		return false
	}

	if arg0.Elem() == contextType {
		return true
	}

	return false
}

func (s *Server) route(c http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path
	ctx := &Context{Request: r, ResponseWriter: c, Server: s}
	ctx.Init()

	for _, route := range s.routes {
		match, result := matchRoute(route, requestPath)
		if !match {
			continue
		} else {
			var args []reflect.Value

			if requireContext(&route.handler) {
				args = append(args, reflect.ValueOf(ctx))
			}

			for _, arg := range result {
				args = append(args, reflect.ValueOf(arg))
			}

			ret := route.handler.Call(args)

			if len(ret) < 1 {
				return
			}

			ret0 := ret[0]

			var content []byte
			if ret0.Kind() == reflect.String {
				content = []byte(ret0.String())
			} else if ret0.Kind() == reflect.Map {
				json_content := make(map[string]interface{})
				for _, k := range ret0.MapKeys() {
					json_content[k.String()] = ret0.MapIndex(k).Interface()
				}

				json_string, err := json.Marshal(json_content)

				if err != nil {
					ctx.Abort(500, "Internal Error")
					return
				}

				ctx.Header().Set("Content-Type", "application/json")
				content = json_string
			} else if ret0.Kind() == reflect.Struct {
				json_content := make(map[string]interface{})
				type_ret := ret0.Type()
				for i := 0; i < ret0.NumField(); i++ {
					f := ret0.Field(i)
					if f.CanInterface() {
						// Only jsonify the exported field
						json_content[type_ret.Field(i).Name] = f.Interface()
					}
				}

				json_string, err := json.Marshal(json_content)

				if err != nil {
					ctx.Abort(500, "Internal Error")
					return
				}

				ctx.Header().Set("Content-Type", "application/json")
				content = json_string
			}

			if len(content) < 1 {
				// ctx.Abort(500, "Internal Error")
				return
			}

			ctx.Header().Set("Content-Length", strconv.Itoa(len(content)))

			ctx.Write(content)
		}

		return
	}

	ctx.NotFound()
}
