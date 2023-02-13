package main

import (
	"fmt"
	"log"
	"net/http"
)

func server1_main() {
	// Amennyiben a pattern végén "/" áll a pattern az összes url-re illeszkedik,
	// ha nincs jobbabn illeszkedő pattern. Ha a pattern nem "/"-re végződik, csak
	// pontos illeszkedés az elfogadott.
	http.HandleFunc("/", landingPageHandler)
	http.HandleFunc("/hello/", helloPageHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func landingPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func helloPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}
