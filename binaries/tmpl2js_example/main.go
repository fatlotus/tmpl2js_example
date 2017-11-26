package main

import (
	"github.com/fatlotus/tmpl2js_example"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", &tmpl2js_example.Counter{})
}
