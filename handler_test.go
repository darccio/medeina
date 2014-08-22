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
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "This is the home page")
	})
	mux.Handle("/list", http.HandlerFunc(list))
	mr := NewMedeina()
	mr.OnHandler("api/v1/events", mux)
	return mr
}

func list(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a list")
}

func TestEmbeddedNegroni(t *testing.T) {
	testRequests(t, "GET", "/api/v1/events/", http.StatusOK)
	testRequests(t, "GET", "/api/v1/events/list", http.StatusOK)
	testRequests(t, "GET", "/api/v1/events/listx", http.StatusNotFound)
}

func testRequests(t *testing.T, method, path string, expectedStatus int) {
	r, _ := http.NewRequest(method, path, nil)
	w := new(httptest.ResponseRecorder)
	u := r.URL
	r.RequestURI = u.RequestURI()
	standard.ServeHTTP(w, r)
	if w.Code != expectedStatus {
		t.Errorf("Expected route %s %s found: Code=%d", method, u, w.Code)
		panic(t)
	}
}
