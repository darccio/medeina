// Copyright (c) 2014 Dario Castañé. Licensed under the MIT License.
package medeina

import (
	"bytes"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/oleiade/lane"
	"net/http"
	"net/url"
	"strings"
)

// Internal router struct. It can be useful to keep Medeina
// router-agnostic.
type router struct {
	*httprouter.Router
}

// Medeina is a goddess willing to help you with your trees... of routes.
// Allow it to be part of your chain of HTTP Handlers and she will handle
// all those messy branches that your once-used-to-be-simple router got.
type Medeina struct {
	router  *router
	methods *lane.Stack
	path    *lane.Deque
}

// Medeina closures definition.
type Handle func()

// HTTP Methods available as constants.
// We could use strings but it was cleaner to force
// specefic values in an enum-like fashion.
type Method string

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PATCH  = "PATCH"
	DELETE = "DELETE"
)

var Methods = []Method{ GET, POST, PUT, PATCH, DELETE }

// Joins a deque using slashes. This is not a
// generic function.
func joinDeque(s *lane.Deque) string {
	var (
		buffer bytes.Buffer
		bDeque *lane.Deque
	)
	bDeque = lane.NewDeque()
	for e := s.Shift(); e != nil; e = s.Shift() {
		subpath := fmt.Sprintf("%v", e)
		if !(subpath == "" && s.Empty()) {
			buffer.WriteString("/")
			buffer.WriteString(subpath)
		}
		bDeque.Append(e)
	}
	for e := bDeque.Shift(); e != nil; e = bDeque.Shift() {
		s.Append(e)
	}
	return buffer.String()
}

var (
	// Make sure this conforms with the http.Handle interface
	// as in julienschmidt/httprouter.
	_ http.Handler = (*Medeina)(nil)
)

// Returns a new initialized Medeina tree routing with default httprouter's one.
func NewMedeina() *Medeina {
	return &Medeina{
		router: &router{
			httprouter.New(),
		},
		methods: lane.NewStack(),
		path:    lane.NewDeque(),
	}
}

// Core logic of handling routes in a tree.
func (m *Medeina) handle(method Method, handle Handle) {
	m.methods.Push(method)
	handle()
	m.methods.Pop()
}

// Switches context to use GET method as default in the closure.
// You can override for a route while using Is, setting which
// methods you want.
func (m *Medeina) GET(handles Handle) {
	m.handle("GET", handles)
}

// Switches context to use POST method as default in the closure.
func (m *Medeina) POST(handles Handle) {
	m.handle("POST", handles)
}

// Switches context to use PUT method as default in the closure.
func (m *Medeina) PUT(handles Handle) {
	m.handle("PUT", handles)
}

// Switches context to use PATCH method as default in the closure.
func (m *Medeina) PATCH(handles Handle) {
	m.handle("PATCH", handles)
}

// Switches context to use DELETE method as default in the closure.
func (m *Medeina) DELETE(handles Handle) {
	m.handle("DELETE", handles)
}

// Adds a new subpath to the current context. Everything under the
// closure will use all the previously set path as root for their
// URLs.
func (m *Medeina) On(path string, handle Handle) {
	m.path.Append(path)
	handle()
	m.path.Pop()
}

// As On but using a function which accepts a routing tree as parameter.
// This will be useful to split routes definition in several functions.
func (m *Medeina) OnFunc(path string, handle func(*Medeina)) {
	m.path.Append(path)
	handle(m)
	m.path.Pop()
}

// As On but using a function which accepts a standard http.Handler,
// delegating further route handling to the handler. It adds a HttpRouter
// catch-all matcher called 'medeina_subpath'.
// This will be useful to split routes definition in several functions.
func (m *Medeina) OnHandler(path string, handle http.Handler) {
	m.path.Append(path)
	m.Handler("*medeina_subpath", handle, Methods...)
	m.path.Pop()
}

// Sets a canonical path. A canonical path means no further entries are in the path.
func (m *Medeina) Is(path string, handle httprouter.Handle, methods ...Method) {
	m.path.Append(path)
	fullPath := joinDeque(m.path)
	m.path.Pop()
	// If any method is provided, it overrides the default one.
	if len(methods) > 0 {
		for _, method := range methods {
			sm := string(method)
			m.router.Handle(sm, fullPath, handle)
		}
	} else {
		method := m.methods.Head()
		if method == nil {
			panic(fmt.Errorf("you cannot set an endpoint outside a HTTP method scope or without passing methods by parameter"))
		}
		m.router.Handle(string(method.(Method)), fullPath, handle)
	}
}

// As Is but delegateing on a standard http.Handler.
// There is no equivalent functions for specific HTTP methods, so you must use
// this in order to add standard http.Handlers.
func (m *Medeina) Handler(path string, handle http.Handler, methods ...Method) {
	m.path.Append(path)
	fullPath := joinDeque(m.path)
	m.path.Pop()
	// If any method is provided, it overrides the default one.
	if len(methods) > 0 {
		for _, method := range methods {
			sm := string(method)
			m.router.Handler(sm, fullPath, handle)
		}
	} else {
		method := m.methods.Head()
		if method == nil {
			panic(fmt.Errorf("you cannot set an endpoint outside a HTTP method scope or without passing methods by parameter"))
		}
		m.router.Handler(string(method.(Method)), fullPath, handle)
	}
}

// Utility function to use with http.Handler compatible routers. Modifies
// the request's URL in order to make subrouters relative to the prefix.
// If you use a router as subrouter without this they need to match the full
// path.
func HandlerPathPrefix(prefix string, handle http.Handler) http.Handler {
	if !strings.HasPrefix(prefix, "/") {
		prefix = fmt.Sprintf("/%s", prefix)
	}
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		old := r.URL
		r.URL, _ = url.ParseRequestURI(strings.Replace(old.Path, prefix, "", 1))
		handle.ServeHTTP(w, r)
		r.URL = old
	})
}

// Makes the routing tree implement the http.Handler interface.
func (m *Medeina) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}
