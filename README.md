# Medeina [![Build Status][1]][2] [![GoDoc](https://godoc.org/github.com/imdario/medeina?status.svg)](https://godoc.org/github.com/imdario/medeina) [![docs examples](https://sourcegraph.com/api/repos/github.com/imdario/medeina/.badges/docs-examples.png)](https://sourcegraph.com/github.com/imdario/medeina) [![dependencies](https://sourcegraph.com/api/repos/github.com/imdario/medeina/.badges/dependencies.png)](https://sourcegraph.com/github.com/imdario/medeina)

[1]: https://travis-ci.org/imdario/medeina.png
[2]: https://travis-ci.org/imdario/medeina

[Medeina](https://github.com/imdario/medeina) is a Go routing tree based on [HttpRouter](https://github.com/julienschmidt/httprouter) and inspired by Ruby's [Roda](http://roda.jeremyevans.net/) and [Cuba](http://cuba.is/). It allows to define the HTTP routes of your web application as a tree, operating on the current route at any point of the tree.

As stated in Roda's website, "this allows you to have much DRYer code". All the routes can have all the features of HttpRouter: named paramaters, catch-all parameters, etc.

Actually, Medeina inherits all the performance and flexibility you love in HttpRouter (great job, Julien!).

## Status

Medeina is fully functional and tested but it's still green and young. It may lack some useful functionality. If you use HttpRouter, give Medeina a try. All real world experience is welcome.

## Install

    go get github.com/imdario/medeina

    // use in your .go code
    import (
        "github.com/imdario/medeina"
    )

## Usage

Check the [docs](https://godoc.org/github.com/imdario/medeina) and these examples:

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

    log.Fatal(http.ListenAndServe(":8080", r))


And this for a quick and dirty "REST-like" API (ignore the fact that I'm using a single method to handle everything):

    endpoints := []string {
        "department", "employee", "project"
    }
    r := medeina.NewMedeina()
    r.GET(func() {
        r.Is("", Index)
    })
    for endpoint := range endpoints {
        r.On(endpoint, func() {
            r.GET(func() {
                r.Is(":id", Model)
            })
            r.POST(func() {
                r.Is("", Model)
            })
            r.PUT(func() {
                r.Is("", Model)
            })
            r.DELETE(func() {
                r.Is(":id", Model)
            })
        })
    }

    log.Fatal(http.ListenAndServe(":8080", r))

Or it's "shorter" (YMMV) form:

    endpoints := []string {
        "department", "employee", "project"
    }
    r := medeina.NewMedeina()
    r.Is("", Index, medeina.GET)
    for endpoint := range endpoints {
        r.On(endpoint, func() {
            r.Is("", Model, medeina.POST, medeina.PUT)
            r.Is(":id", Model, medeina.GET, medeina.DELETE)
        })
    }

    log.Fatal(http.ListenAndServe(":8080", r))

What does "" means? It matches the current path in your route. It's a way to easily match the scope itself as canonical URL.

## Why HttpRouter?

Because it's the most fast and flexible Go HTTP router around the town and a good one to start. If you want Medeina to work with your preferred option, patches are welcome!

If you don't know [HttpRouter](https://github.com/julienschmidt/httprouter), please check it out. You won't regret it.

## Why is it called Medeina?

From [Wikipedia](https://en.wikipedia.org/wiki/Medeina):

> Medeina or Medeinė (derived from medis (tree) and medė (forest)), [...] is one of the main deities in the Lithuanian mythology, similar to Latvian Meža Māte. She is a ruler of forests, trees and animals. Her sacred animal is a hare.

Hey, we were talking about trees. It fits right! Also, this project can join my other ones, also called by names starting by 'm': [Mergo](https://github.com/imdario/mergo), [Minshu](https://github.com/imdario/minshu), [mqqsig192](https://github.com/imdario/mqqsig192), etc. Don't ask, it wasn't on purpose.

## Contact me

If I can help you, you have an idea or you are using Medeina in your projects, don't hesitate to drop me a line (or a pull request): [@im_dario](https://twitter.com/im_dario)

## About

Written by [Dario Castañé](http://dario.im).

## License

[MIT](http://opensource.org/licenses/MIT) license.
