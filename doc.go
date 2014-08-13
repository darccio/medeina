// Copyright (c) 2014 Dario Castañé. Licensed under the MIT License.

/*
Medeina is a routing tree based on httprouter inspired by Ruby's Roda and Cuba. It allows to define your routes as a tree, operating on the current one at any point of the routing tree.

Usage

	// From Roda's site
	r := medeina.NewMedeina()
	r.GET(func() {
		r.Is("", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			http.Redirect(w, r, "/hello", http.StatusFound)
		})
		r.On("hello", func() {
			r.Is("world", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				fmt.Fprintf(w, "Hello world!")
			})
			r.Is("", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				fmt.Fprintf(w, "Hello!")
			})
		})
	})

If you want to check a more convoluted example, feel free to check github_test.go. It uses several possible styles of creating routes using Medeina in order to avoid pitfalls.
*/
package medeina
