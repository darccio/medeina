// Copyright 2013 Julien Schmidt, All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.
// Adapted from julienschmidt/go-http-routing-benchmark/github_test.go
// Copyright (c) 2014 Dario Castañé. Licensed under the MIT License.

package medeina

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type route struct {
	method string
	path   string
}

// http://developer.github.com/v3/
var githubAPI = []route{
	// OAuth Authorizations
	{"GET", "/authorizations"},
	{"GET", "/authorizations/:id"},
	{"POST", "/authorizations"},
	//{"PUT", "/authorizations/clients/:client_id"},
	//{"PATCH", "/authorizations/:id"},
	{"DELETE", "/authorizations/:id"},
	{"GET", "/applications/:client_id/tokens/:access_token"},
	{"DELETE", "/applications/:client_id/tokens"},
	{"DELETE", "/applications/:client_id/tokens/:access_token"},

	// Activity
	{"GET", "/events"},
	{"GET", "/repos/:owner/:repo/events"},
	{"GET", "/networks/:owner/:repo/events"},
	{"GET", "/orgs/:org/events"},
	{"GET", "/users/:user/received_events"},
	{"GET", "/users/:user/received_events/public"},
	{"GET", "/users/:user/events"},
	{"GET", "/users/:user/events/public"},
	{"GET", "/users/:user/events/orgs/:org"},
	{"GET", "/feeds"},
	{"GET", "/notifications"},
	{"GET", "/repos/:owner/:repo/notifications"},
	{"PUT", "/notifications"},
	{"PUT", "/repos/:owner/:repo/notifications"},
	{"GET", "/notifications/threads/:id"},
	//{"PATCH", "/notifications/threads/:id"},
	{"GET", "/notifications/threads/:id/subscription"},
	{"PUT", "/notifications/threads/:id/subscription"},
	{"DELETE", "/notifications/threads/:id/subscription"},
	{"GET", "/repos/:owner/:repo/stargazers"},
	{"GET", "/users/:user/starred"},
	{"GET", "/user/starred"},
	{"GET", "/user/starred/:owner/:repo"},
	{"PUT", "/user/starred/:owner/:repo"},
	{"DELETE", "/user/starred/:owner/:repo"},
	{"GET", "/repos/:owner/:repo/subscribers"},
	{"GET", "/users/:user/subscriptions"},
	{"GET", "/user/subscriptions"},
	{"GET", "/repos/:owner/:repo/subscription"},
	{"PUT", "/repos/:owner/:repo/subscription"},
	{"DELETE", "/repos/:owner/:repo/subscription"},
	{"GET", "/user/subscriptions/:owner/:repo"},
	{"PUT", "/user/subscriptions/:owner/:repo"},
	{"DELETE", "/user/subscriptions/:owner/:repo"},

	// Gists
	{"GET", "/users/:user/gists"},
	{"GET", "/gists"},
	//{"GET", "/gists/public"},
	//{"GET", "/gists/starred"},
	{"GET", "/gists/:id"},
	{"POST", "/gists"},
	//{"PATCH", "/gists/:id"},
	{"PUT", "/gists/:id/star"},
	{"DELETE", "/gists/:id/star"},
	{"GET", "/gists/:id/star"},
	{"POST", "/gists/:id/forks"},
	{"DELETE", "/gists/:id"},

	// Git Data
	{"GET", "/repos/:owner/:repo/git/blobs/:sha"},
	{"POST", "/repos/:owner/:repo/git/blobs"},
	{"GET", "/repos/:owner/:repo/git/commits/:sha"},
	{"POST", "/repos/:owner/:repo/git/commits"},
	//{"GET", "/repos/:owner/:repo/git/refs/*ref"},
	{"GET", "/repos/:owner/:repo/git/refs"},
	{"POST", "/repos/:owner/:repo/git/refs"},
	//{"PATCH", "/repos/:owner/:repo/git/refs/*ref"},
	//{"DELETE", "/repos/:owner/:repo/git/refs/*ref"},
	{"GET", "/repos/:owner/:repo/git/tags/:sha"},
	{"POST", "/repos/:owner/:repo/git/tags"},
	{"GET", "/repos/:owner/:repo/git/trees/:sha"},
	{"POST", "/repos/:owner/:repo/git/trees"},

	// Issues
	{"GET", "/issues"},
	{"GET", "/user/issues"},
	{"GET", "/orgs/:org/issues"},
	{"GET", "/repos/:owner/:repo/issues"},
	{"GET", "/repos/:owner/:repo/issues/:number"},
	{"POST", "/repos/:owner/:repo/issues"},
	//{"PATCH", "/repos/:owner/:repo/issues/:number"},
	{"GET", "/repos/:owner/:repo/assignees"},
	{"GET", "/repos/:owner/:repo/assignees/:assignee"},
	{"GET", "/repos/:owner/:repo/issues/:number/comments"},
	//{"GET", "/repos/:owner/:repo/issues/comments"},
	//{"GET", "/repos/:owner/:repo/issues/comments/:id"},
	{"POST", "/repos/:owner/:repo/issues/:number/comments"},
	//{"PATCH", "/repos/:owner/:repo/issues/comments/:id"},
	//{"DELETE", "/repos/:owner/:repo/issues/comments/:id"},
	{"GET", "/repos/:owner/:repo/issues/:number/events"},
	//{"GET", "/repos/:owner/:repo/issues/events"},
	//{"GET", "/repos/:owner/:repo/issues/events/:id"},
	{"GET", "/repos/:owner/:repo/labels"},
	{"GET", "/repos/:owner/:repo/labels/:name"},
	{"POST", "/repos/:owner/:repo/labels"},
	//{"PATCH", "/repos/:owner/:repo/labels/:name"},
	{"DELETE", "/repos/:owner/:repo/labels/:name"},
	{"GET", "/repos/:owner/:repo/issues/:number/labels"},
	{"POST", "/repos/:owner/:repo/issues/:number/labels"},
	{"DELETE", "/repos/:owner/:repo/issues/:number/labels/:name"},
	{"PUT", "/repos/:owner/:repo/issues/:number/labels"},
	{"DELETE", "/repos/:owner/:repo/issues/:number/labels"},
	{"GET", "/repos/:owner/:repo/milestones/:number/labels"},
	{"GET", "/repos/:owner/:repo/milestones"},
	{"GET", "/repos/:owner/:repo/milestones/:number"},
	{"POST", "/repos/:owner/:repo/milestones"},
	//{"PATCH", "/repos/:owner/:repo/milestones/:number"},
	{"DELETE", "/repos/:owner/:repo/milestones/:number"},

	// Miscellaneous
	{"GET", "/emojis"},
	{"GET", "/gitignore/templates"},
	{"GET", "/gitignore/templates/:name"},
	{"POST", "/markdown"},
	{"POST", "/markdown/raw"},
	{"GET", "/meta"},
	{"GET", "/rate_limit"},

	// Organizations
	{"GET", "/users/:user/orgs"},
	{"GET", "/user/orgs"},
	{"GET", "/orgs/:org"},
	//{"PATCH", "/orgs/:org"},
	{"GET", "/orgs/:org/members"},
	{"GET", "/orgs/:org/members/:user"},
	{"DELETE", "/orgs/:org/members/:user"},
	{"GET", "/orgs/:org/public_members"},
	{"GET", "/orgs/:org/public_members/:user"},
	{"PUT", "/orgs/:org/public_members/:user"},
	{"DELETE", "/orgs/:org/public_members/:user"},
	{"GET", "/orgs/:org/teams"},
	{"GET", "/teams/:id"},
	{"POST", "/orgs/:org/teams"},
	//{"PATCH", "/teams/:id"},
	{"DELETE", "/teams/:id"},
	{"GET", "/teams/:id/members"},
	{"GET", "/teams/:id/members/:user"},
	{"PUT", "/teams/:id/members/:user"},
	{"DELETE", "/teams/:id/members/:user"},
	{"GET", "/teams/:id/repos"},
	{"GET", "/teams/:id/repos/:owner/:repo"},
	{"PUT", "/teams/:id/repos/:owner/:repo"},
	{"DELETE", "/teams/:id/repos/:owner/:repo"},
	{"GET", "/user/teams"},

	// Pull Requests
	{"GET", "/repos/:owner/:repo/pulls"},
	{"GET", "/repos/:owner/:repo/pulls/:number"},
	{"POST", "/repos/:owner/:repo/pulls"},
	//{"PATCH", "/repos/:owner/:repo/pulls/:number"},
	{"GET", "/repos/:owner/:repo/pulls/:number/commits"},
	{"GET", "/repos/:owner/:repo/pulls/:number/files"},
	{"GET", "/repos/:owner/:repo/pulls/:number/merge"},
	{"PUT", "/repos/:owner/:repo/pulls/:number/merge"},
	{"GET", "/repos/:owner/:repo/pulls/:number/comments"},
	//{"GET", "/repos/:owner/:repo/pulls/comments"},
	//{"GET", "/repos/:owner/:repo/pulls/comments/:number"},
	{"PUT", "/repos/:owner/:repo/pulls/:number/comments"},
	//{"PATCH", "/repos/:owner/:repo/pulls/comments/:number"},
	//{"DELETE", "/repos/:owner/:repo/pulls/comments/:number"},

	// Repositories
	{"GET", "/user/repos"},
	{"GET", "/users/:user/repos"},
	{"GET", "/orgs/:org/repos"},
	{"GET", "/repositories"},
	{"POST", "/user/repos"},
	{"POST", "/orgs/:org/repos"},
	{"GET", "/repos/:owner/:repo"},
	//{"PATCH", "/repos/:owner/:repo"},
	{"GET", "/repos/:owner/:repo/contributors"},
	{"GET", "/repos/:owner/:repo/languages"},
	{"GET", "/repos/:owner/:repo/teams"},
	{"GET", "/repos/:owner/:repo/tags"},
	{"GET", "/repos/:owner/:repo/branches"},
	{"GET", "/repos/:owner/:repo/branches/:branch"},
	{"DELETE", "/repos/:owner/:repo"},
	{"GET", "/repos/:owner/:repo/collaborators"},
	{"GET", "/repos/:owner/:repo/collaborators/:user"},
	{"PUT", "/repos/:owner/:repo/collaborators/:user"},
	{"DELETE", "/repos/:owner/:repo/collaborators/:user"},
	{"GET", "/repos/:owner/:repo/comments"},
	{"GET", "/repos/:owner/:repo/commits/:sha/comments"},
	{"POST", "/repos/:owner/:repo/commits/:sha/comments"},
	{"GET", "/repos/:owner/:repo/comments/:id"},
	//{"PATCH", "/repos/:owner/:repo/comments/:id"},
	{"DELETE", "/repos/:owner/:repo/comments/:id"},
	{"GET", "/repos/:owner/:repo/commits"},
	{"GET", "/repos/:owner/:repo/commits/:sha"},
	{"GET", "/repos/:owner/:repo/readme"},
	//{"GET", "/repos/:owner/:repo/contents/*path"},
	//{"PUT", "/repos/:owner/:repo/contents/*path"},
	//{"DELETE", "/repos/:owner/:repo/contents/*path"},
	//{"GET", "/repos/:owner/:repo/:archive_format/:ref"},
	{"GET", "/repos/:owner/:repo/keys"},
	{"GET", "/repos/:owner/:repo/keys/:id"},
	{"POST", "/repos/:owner/:repo/keys"},
	//{"PATCH", "/repos/:owner/:repo/keys/:id"},
	{"DELETE", "/repos/:owner/:repo/keys/:id"},
	{"GET", "/repos/:owner/:repo/downloads"},
	{"GET", "/repos/:owner/:repo/downloads/:id"},
	{"DELETE", "/repos/:owner/:repo/downloads/:id"},
	{"GET", "/repos/:owner/:repo/forks"},
	{"POST", "/repos/:owner/:repo/forks"},
	{"GET", "/repos/:owner/:repo/hooks"},
	{"GET", "/repos/:owner/:repo/hooks/:id"},
	{"POST", "/repos/:owner/:repo/hooks"},
	//{"PATCH", "/repos/:owner/:repo/hooks/:id"},
	{"POST", "/repos/:owner/:repo/hooks/:id/tests"},
	{"DELETE", "/repos/:owner/:repo/hooks/:id"},
	{"POST", "/repos/:owner/:repo/merges"},
	{"GET", "/repos/:owner/:repo/releases"},
	{"GET", "/repos/:owner/:repo/releases/:id"},
	{"POST", "/repos/:owner/:repo/releases"},
	//{"PATCH", "/repos/:owner/:repo/releases/:id"},
	{"DELETE", "/repos/:owner/:repo/releases/:id"},
	{"GET", "/repos/:owner/:repo/releases/:id/assets"},
	{"GET", "/repos/:owner/:repo/stats/contributors"},
	{"GET", "/repos/:owner/:repo/stats/commit_activity"},
	{"GET", "/repos/:owner/:repo/stats/code_frequency"},
	{"GET", "/repos/:owner/:repo/stats/participation"},
	{"GET", "/repos/:owner/:repo/stats/punch_card"},
	{"GET", "/repos/:owner/:repo/statuses/:ref"},
	{"POST", "/repos/:owner/:repo/statuses/:ref"},

	// Search
	{"GET", "/search/repositories"},
	{"GET", "/search/code"},
	{"GET", "/search/issues"},
	{"GET", "/search/users"},
	{"GET", "/legacy/issues/search/:owner/:repository/:state/:keyword"},
	{"GET", "/legacy/repos/search/:keyword"},
	{"GET", "/legacy/user/search/:keyword"},
	{"GET", "/legacy/user/email/:email"},

	// Users
	{"GET", "/users/:user"},
	{"GET", "/user"},
	//{"PATCH", "/user"},
	{"GET", "/users"},
	{"GET", "/user/emails"},
	{"POST", "/user/emails"},
	{"DELETE", "/user/emails"},
	{"GET", "/users/:user/followers"},
	{"GET", "/user/followers"},
	{"GET", "/users/:user/following"},
	{"GET", "/user/following"},
	{"GET", "/user/following/:user"},
	{"GET", "/users/:user/following/:target_user"},
	{"PUT", "/user/following/:user"},
	{"DELETE", "/user/following/:user"},
	{"GET", "/users/:user/keys"},
	{"GET", "/user/keys"},
	{"GET", "/user/keys/:id"},
	{"POST", "/user/keys"},
	//{"PATCH", "/user/keys/:id"},
	{"DELETE", "/user/keys/:id"},
}

var medeina http.Handler

func init() {
	medeina = loadMedeina()
}

// Hack to refactor some common patterns.
func setEndpoints(mr *Medeina, handler httprouter.Handle, endpoints []string) {
	for _, endpoint := range endpoints {
		mr.Is(endpoint, handler)
	}
}

// Hack to refactor some common patterns.
func setBasicEndpoints(mr *Medeina, name string, post, del bool) {
	mr.On(name, func() {
		mr.GET(func() {
			mr.Is("", testHandlerParams)
			mr.Is(":id", testHandlerParams)
		})
		if post {
			mr.Is("", testHandlerParams, POST)
		}
		if del {
			mr.Is(":id", testHandlerParams, DELETE)
		}
	})
}

// Builds a Medeina router creating an exact copy of githubAPI.
// It's a mess **on purpose**. I wanted to use different styles of working with Medeina, in order to find any pitfalls.
func loadMedeina() http.Handler {
	mr := NewMedeina()
	mr.On("authorizations", func() {
		mr.Is("", testHandler, GET, POST)
		mr.Is(":id", testHandler, GET, DELETE)
	})
	mr.On("user", func() {
		mr.On("starred", func() {
			mr.Is("", testHandler, GET)
			mr.Is(":owner/:repo", testHandler, GET, PUT, DELETE)
		})
		mr.On("subscriptions", func() {
			mr.Is("", testHandler, GET)
			mr.Is(":owner/:repo", testHandler, GET, PUT, DELETE)
		})
		mr.On("repos", func() {
			mr.Is("", testHandler, GET, POST)
		})
		mr.Is("issues", testHandler, GET)
		mr.Is("orgs", testHandler, GET)
		mr.Is("teams", testHandler, GET)
		mr.Is("emails", testHandler, GET, POST, DELETE)
		mr.Is("followers", testHandler, GET)
		mr.On("following", func() {
			mr.Is("", testHandler, GET)
			mr.Is(":user", testHandler, GET, PUT, DELETE)
		})
		mr.On("keys", func() {
			mr.Is("", testHandler, GET, POST)
			mr.Is(":id", testHandler, GET, DELETE)
		})
	})
	mr.On("repos/:owner/:repo", func() {
		mr.GET(func() {
			endpoints := []string{
				"", "events", "notifications", "stargazers", "subscribers",
				"contributors", "languages", "teams", "tags", "readme", "forks",
				"assignees", "assignees/:assignee", "milestones",
				"milestones/:number/labels", "milestones/:number",
				"branches", "branches/:branch", "collaborators", "collaborators/:user",
				"releases/:id/assets", "statuses/:ref",
			}
			setEndpoints(mr, testHandlerParams, endpoints)
		})
		mr.Is("subscription", testHandlerParams, GET, PUT, DELETE)
		mr.On("git", func() {
			endpoints := []string{
				"blobs", "commits", "refs", "tags", "trees",
			}
			for _, endpoint := range endpoints {
				mr.POST(func() {
					mr.Is(endpoint, testHandlerParams)
				})
				mr.GET(func() {
					if endpoint != "refs" {
						mr.Is(fmt.Sprintf("%s/:sha", endpoint), testHandlerParams)
					}
				})
			}
		})
		mr.On("issues", func() {
			mr.GET(func() {
				mr.Is("", testHandlerParams)
				mr.On(":number", func() {
					endpoints := []string{
						"", "comments", "events", "labels",
					}
					setEndpoints(mr, testHandlerParams, endpoints)
				})
			})
		})
		mr.Is("pulls", testHandlerParams, GET, POST)
		mr.On("pulls/:number", func() {
			mr.GET(func() {
				endpoints := []string{
					"commits", "files", "merge", "comments",
				}
				setEndpoints(mr, testHandlerParams, endpoints)
			})
			mr.PUT(func() {
				mr.Is("merge", testHandlerParams)
				mr.Is("comments", testHandlerParams)
			})
		})
		mr.On("commits", func() {
			f := func() {
				mr.Is(":sha/comments", testHandlerParams)
			}
			mr.GET(f)
			mr.POST(f)
			mr.Is(":sha", testHandlerParams, GET)
		})
		mr.On("stats", func() {
			mr.GET(func() {
				endpoints := []string{
					"contributors", "commit_activity", "code_frequency",
					"participation", "punch_card",
				}
				setEndpoints(mr, testHandlerParams, endpoints)
			})
		})
		mr.POST(func() {
			mr.Is("issues/:number/comments", testHandlerParams)
			mr.Is("labels", testHandlerParams)
			mr.Is("issues/:number/labels", testHandlerParams)
			mr.Is("milestones", testHandlerParams)
			mr.Is("forks", testHandlerParams)
			mr.Is("merges", testHandlerParams)
			mr.Is("hooks/:id/tests", testHandlerParams)
			mr.Is("statuses/:ref", testHandlerParams)
		})
		mr.DELETE(func() {
			mr.Is("labels/:name", testHandlerParams)
			mr.Is("issues/:number/labels/:name", testHandlerParams)
			mr.Is("issues/:number/labels", testHandlerParams)
			mr.Is("milestones/:number", testHandlerParams)
			mr.Is("collaborators/:user", testHandlerParams)
			mr.Is("", testHandlerParams)
		})
		mr.PUT(func() {
			mr.Is("issues/:number/labels", testHandlerParams)
			mr.Is("collaborators/:user", testHandlerParams)
		})
		setBasicEndpoints(mr, "keys", true, true)
		setBasicEndpoints(mr, "downloads", false, true)
		setBasicEndpoints(mr, "comments", false, true)
		setBasicEndpoints(mr, "hooks", true, true)
		setBasicEndpoints(mr, "releases", true, true)
	})
	mr.On("applications/:client_id/tokens", func() {
		mr.Is("", testHandler, DELETE)
		mr.Is(":access_token", testHandler, GET, DELETE)
	})
	mr.On("orgs/:org", func() {
		mr.GET(func() {
			endpoints := []string{
				"", "repos", "events", "issues", "teams",
			}
			setEndpoints(mr, testHandler, endpoints)
			endpoints = []string{
				"members", "public_members",
			}
			for _, endpoint := range endpoints {
				mr.On(endpoint, func() {
					mr.Is("", testHandler)
					mr.Is(":user", testHandler)
				})
			}
		})
		mr.POST(func() {
			endpoints := []string{
				"repos", "teams",
			}
			setEndpoints(mr, testHandler, endpoints)
		})
		mr.DELETE(func() {
			mr.Is("members/:user", testHandler)
			mr.Is("public_members/:user", testHandler)
		})
		mr.PUT(func() {
			mr.Is("public_members/:user", testHandler)
		})
	})
	mr.On("users/:user", func() {
		mr.GET(func() {
			endpoints := []string{
				"", "received_events", "received_events/public", "events",
				"events/public", "events/orgs/:org", "starred", "subscriptions",
				"orgs", "repos", "followers", "following", "keys",
				"following/:target_user",
			}
			setEndpoints(mr, testHandler, endpoints)
		})
	})
	mr.On("gists", func() {
		mr.Is("", testHandler, GET, POST)
		mr.On(":id", func() {
			mr.Is("", testHandler, GET, DELETE)
			mr.Is("star", testHandler, GET, PUT, DELETE)
			mr.Is("forks", testHandler, POST)
		})
	})
	mr.On("team/:id", func() {
		f := func(base string, path string) {
			mr.Is(base, testHandler, GET)
			mr.On(fmt.Sprintf("%s/%s", base, path), func() {
				mr.Is("", testHandler, GET, PUT, DELETE)
			})
		}
		f("repos", ":owner/repo")
		f("members", ":user")
	})
	mr.GET(func() {
		endpoints := []string{
			"emojis", "gitignore/templates", "gitignore/templates/:name", "meta",
			"rate_limit", "events", "networks/:owner/:repo/events", "feeds",
			"search/repositories", "search/code", "search/issues", "search/users",
			"legacy/issues/search/:owner/:repository/:state/:keyword",
			"legacy/repos/search/:keyword", "legacy/user/search/:keyword",
			"legacy/user/email/:email", "issues",
		}
		setEndpoints(mr, testHandler, endpoints)
		mr.On("notifications", func() {
			mr.Is("", testHandler)
			mr.Is("threads/:id", testHandler)
			mr.Is("threads/:id/subscription", testHandler)
		})
	})
	mr.POST(func() {
		mr.On("markdown", func() {
			mr.Is("", testHandler)
			mr.Is("raw", testHandler)
		})
	})
	mr.PUT(func() {
		mr.On("notifications", func() {
			mr.Is("threads/:id/subscription", testHandler)
		})
	})
	mr.DELETE(func() {
		mr.Is("notifications/threads/:id/subscription", testHandler)
	})
	return mr
}

func testHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Medeina or Medeinė (derived from medis (tree) and medė (forest)) [...] is a one of the main deities in the Lithuanian mythology.")
}

func testHandlerParams(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if strings.HasPrefix(r.URL.String(), "/repos/imdario/medeina") {
		if p.ByName("owner") != "imdario" && p.ByName("repo") != "medeina" {
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprint(w, "Medeina is single, unwilling to get married, though voluptuous and beautiful huntress.")
	} else {
		fmt.Fprint(w, "She is depicted as a young woman and a she-wolf (cf. vilkmergė) with an escort of wolves. Her duty is not to help the hunters, but to protect the forest.")
	}
}

// Fakes an expected valid request.
func testRequest(t *testing.T, router http.Handler, r *http.Request) {
	w := new(httptest.ResponseRecorder)
	u := r.URL
	r.RequestURI = u.RequestURI()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Handling route %s failed: Code=%d", u, w.Code)
	}
}

func TestStatic(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user/repos", nil)
	testRequest(t, medeina, req)
}

func TestParam(t *testing.T) {
	req, _ := http.NewRequest("GET", "/repos/imdario/medeina/stargazers", nil)
	testRequest(t, medeina, req)
}

func TestAll(t *testing.T) {
	w := new(httptest.ResponseRecorder)
	r, _ := http.NewRequest("GET", "/", nil)
	u := r.URL
	rq := u.RawQuery
	for _, route := range githubAPI {
		r.Method = route.method
		r.RequestURI = route.path
		u.Path = route.path
		u.RawQuery = rq
		medeina.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Errorf("Handling route %s failed: Code=%d", u, w.Code)
		}
	}
}

// Fakes an expected not found request.
func requestNotFound(t *testing.T, method, path string) {
	r, _ := http.NewRequest(method, path, nil)
	w := new(httptest.ResponseRecorder)
	u := r.URL
	r.RequestURI = u.RequestURI()
	medeina.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Errorf("Not expected route %s %s found: Code=%d", method, u, w.Code)
		panic(t)
	}
}

func TestNotFound(t *testing.T) {
	requestNotFound(t, "GET", "/repos/imdario/medeina/stargzrs")
	requestNotFound(t, "GET", "/deadmau5")
	requestNotFound(t, "DELETE", "/pulls")
}

// Builds a map with all the missing methods from githubAPI.
// This assures no extra routes are created. Combined with TestAll,
// we make sure everything was created properly.
func TestOpposite(t *testing.T) {
	routes := make(map[string][]string)
	for _, route := range githubAPI {
		if _, ok := routes[route.path]; !ok {
			routes[route.path] = []string{
				"GET", "DELETE", "PATCH", "POST", "PUT",
			}
		}
		methods := routes[route.path]
		for i, method := range methods {
			if method == route.method {
				routes[route.path] = append(methods[:i], methods[i+1:]...)
				break
			}
		}
	}
	for route := range routes {
		methods := routes[route]
		for _, method := range methods {
			requestNotFound(t, method, route)
		}
	}
}
