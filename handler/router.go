package handler

import (
	"net/http"
	"os"
	"strings"
)

type RouterFunc func(req Req, res Res) 

func Cwd(path string) string {
	d, _ := os.Getwd()
	return d + path
}

// ---------------
// ROUTER 
// ---------------

type router struct {
	routes 		[]route
	middlewares []RouterFunc
}

func RouterInit(hndlr handler) router {
	println("INITIALIZING ROUTER")
	routes := make([]route, 0)
	mids := make([]RouterFunc, 0)
	r := router{routes, mids}
	hndlr.routers = append(hndlr.routers, r)
	return r
}

func (r* router) All(path string, fns ...RouterFunc) {
	println("ADDING ROUTE: ALL | " + path)
	m := make(map[string][]RouterFunc)
	m["ALL"] = fns
	r.routes = append(r.routes, route{path, m, nil})
}

func (r* router) UseMethod(method string, path string, fns ...RouterFunc) {
	m := make(map[string][]RouterFunc)
	m[strings.ToUpper(method)] = fns
	r.routes = append(r.routes, route{path, m, nil})
}

func (r* router) Route(path string) route {
	for _, rt := range r.routes {
		if rt.path == path {
			return rt
		}
	}
	return route{}
}

func (r* router) Load() {
	println("LOADING ROUTER ROUTES")
	if len(r.routes) == 0 {
		println(" >> WARNING: NO ROUTES IN ROUTER")
	}
	for _, rt := range r.routes {
		http.HandleFunc(rt.path, func(w http.ResponseWriter, req *http.Request) {
			for m, fns := range rt.funcs {
				mt := strings.ToUpper(m)
				if mt == "ALL" {
					println("LOADING PATH " + rt.path + "METHOD ALL")
					for _, fn := range fns {
						fn(req, w)
					}
				} else if mt == req.Method {
					println("LOADING PATH " + rt.path + "METHOD " + mt)
					for _, fn := range fns {
						fn(req, w)
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
		})
	}
}

// ---------------
// INDIVIDUAL METHODS
// ---------------

// The GET method requests a representation of the specified resource. Requests using GET should only retrieve data.
func (r* router) Get(path string, fn RouterFunc) {
	r.UseMethod("GET", path, fn)
}

// The HEAD method asks for a response identical to a GET request, but without the response body.
func (r* router) Head(path string, fn RouterFunc) {
	r.UseMethod("HEAD", path, )
}

// The POST method submits an entity to the specified resource, often causing a change in state or side effects on the server.
func (r* router) Post(path string, fn RouterFunc) {
	r.UseMethod("POST", path, fn)
}

// The PUT method replaces all current representations of the target resource with the request payload.
func (r* router) Put(path string, fn RouterFunc) {
	r.UseMethod("PUT", path, fn)
}

// The DELETE method deletes the specified resource.
func (r* router) Delete(path string, fn RouterFunc) {
	r.UseMethod("DELETE", path, fn)
}

// The CONNECT method establishes a tunnel to the server identified by the target resource.
func (r* router) Connect(path string, fn RouterFunc) {
	r.UseMethod("CONNECT", path, fn)
}

// The OPTIONS method describes the communication options for the target resource.
func (r* router) Options(path string, fn RouterFunc) {
	r.UseMethod("OPTIONS", path, fn)
}

// The TRACE method performs a message loop-back test along the path to the target resource.
func (r* router) Trace(path string, fn RouterFunc) {
	r.UseMethod("TRACE", path, fn)
}

// The PATCH method applies partial modifications to a resource.
func (r* router) Patch(path string, fn RouterFunc) {
	r.UseMethod("PATCH", path, fn)
}