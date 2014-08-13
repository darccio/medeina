package medeina

import (
	"bytes"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/oleiade/lane"
	"net/http"
)

type Router struct {
	*httprouter.Router
}

type Medeina struct {
	router *Router
	methods *lane.Stack
	path *lane.Deque
}

type Handle func()

type Method string

const (
	GET = "GET"
	POST = "POST"
	PUT = "PUT"
	PATCH = "PATCH"
	DELETE = "DELETE"
)

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
	// Make sure the Router conforms with the http.Handle interface
	// as in julienschmidt/httprouter.
	_ http.Handler = (*Medeina)(nil)
)

// Returns a new initialized Router with default httprouter's one.
func NewMedeina() *Medeina {
	return &Medeina {
		router: &Router {
			httprouter.New(),
		},
		methods: lane.NewStack(),
		path: lane.NewDeque(),
	}
}

func (m *Medeina) handle(method Method, handles []Handle) {
	m.methods.Push(method)
	for _, handle := range handles {
		handle()
	}
	m.methods.Pop()
}

func (m *Medeina) GET(handles ...Handle) {
	m.handle("GET", handles)
}

func (m *Medeina) POST(handles ...Handle) {
	m.handle("POST", handles)
}

func (m *Medeina) PUT(handles ...Handle) {
	m.handle("PUT", handles)
}

func (m *Medeina) PATCH(handles ...Handle) {
	m.handle("PATCH", handles)
}

func (m *Medeina) DELETE(handles ...Handle) {
	m.handle("DELETE", handles)
}

func (m *Medeina) On(path string, handles ...Handle) {
	m.path.Append(path)
	for _, handle := range handles {
		handle()
	}
	m.path.Pop()
}

func (m *Medeina) OnFunc(path string, handle func (*Medeina)) {
	m.path.Append(path)
	handle(m)
	m.path.Pop()
}

func (m *Medeina) Is(path string, handle httprouter.Handle, methods ...Method) {
	m.path.Append(path)
	fullPath := joinDeque(m.path)
	m.path.Pop()
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

func (m *Medeina) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}
