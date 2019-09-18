package main

import (
	"flag"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var baseURL string

func main() {
	// set up our logger
	log = logrus.New()

	// set up our http server port
	var port string

	// parse any flags that came into the app
	flag.StringVar(&baseURL, "baseURL", "http://localhost", "what http location are you trying to wrap?")
	flag.StringVar(&port, "port", "9000", "what port should we listen on for incoming http requests?")
	flag.Parse()

	// This line of code maps all routes to a single handler function.
	//
	// The first parameter is sort of like a wildcard, so it'll catch all routes
	// that match this pattern. If the pattern is "/", it'll pretty much catch
	// everything.
	http.HandleFunc("/", catchAllHandler)

	// this line of code makes a specific route work.
	// you can put additional routes here and "catch" them in the new go app
	//
	// http.HandleFunc("/anything", dummyHandler)

	// let's give the user some context to the application
	log.Info("Wrapping " + baseURL)
	log.Info("Now listening on port: " + port)
	http.ListenAndServe(":"+port, nil)
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "this is a text response")
}

func catchAllHandler(w http.ResponseWriter, r *http.Request) {
	// Since we are catching pretty much everything here, if we are actually
	// trying to treat the "/" route from everything else, we can specifically
	// look for that route path here
	if r.URL.Path == "/" {
		dummyHandler(w, r)
		return
	}

	// here is where we can create a client, log the uncaught route, and "wrap"
	// an existing application by calling out to it
	log.Error("Uncaught route: " + r.URL.Path)

	// handle different HTTP verbs
	switch r.Method {
	case "GET":
		response, err := http.Get(baseURL + r.URL.Path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		for _, c := range response.Cookies() {
			log.Info(c)
		}
		w.Write(body)
	case "POST":
		response, err := http.Post(baseURL+r.URL.Path, r.Header.Get("Content-Type"), r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(body)
	}
	return
}
