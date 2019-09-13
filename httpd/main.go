package main

import (
	"flag"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var baseURL string

func main() {
	log = logrus.New()

	flag.StringVar(&baseURL, "baseURL", "http://localhost", "what http location are you trying to wrap?")
	flag.Parse()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/anything", anythingHandler)

	http.ListenAndServe(":9000", nil)
}

func anythingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is a normal response"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// here is where we can create a client, and "wrap" an existing application
		log.Error("Uncaught route: " + r.URL.Path)
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
		w.Write(body)
		return
	}

	w.Write([]byte("howdy!"))
}
