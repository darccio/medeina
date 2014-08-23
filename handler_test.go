package medeina

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var standard http.Handler

func init() {
	standard = loadStandard()
}

func loadStandard() http.Handler {
	mux := http.NewServeMux()
	//mux.HandleFunc("/api/v1/events/", func(w http.ResponseWriter, r *http.Request) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//if r.URL.Path != "/api/v1/events/" {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, "This is the home page")
	})
	//mux.HandleFunc("/api/v1/events/list", list)
	mux.HandleFunc("/list", list)
	mr := NewMedeina()
	mr.OnHandler("api/v1/events", HandlerPathPrefix("/api/v1/events", mux))
	return mr
}

func list(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a list")
}

func TestEmbeddedNegroni(t *testing.T) {
	testRequests(t, "GET", "/api/v1/events/", http.StatusOK)
	testRequests(t, "GET", "/api/v1/events/xlistx", http.StatusNotFound)
	testRequests(t, "GET", "/api/v1/events/list", http.StatusOK)
}

func testRequests(t *testing.T, method, path string, expectedStatus int) {
	r, _ := http.NewRequest(method, path, nil)
	w := new(httptest.ResponseRecorder)
	u := r.URL
	r.RequestURI = u.RequestURI()
	standard.ServeHTTP(w, r)
	if w.Code != expectedStatus {
		t.Errorf("Expected %d for route %s %s found: Code=%d", expectedStatus, method, u, w.Code)
		panic(t)
	}
}
