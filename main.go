package main

import (
	"fmt"
	"log"
	"net/http"

	"gee"
)

func main() {
	r := gee.New()

	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %s\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {

		for key, val := range req.Header {
			fmt.Fprintf(w, "Header[%q]=%q\n", key, val)
		}
	})

	log.Fatal(r.Run(":8199"))
}
