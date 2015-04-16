package main

import (
	"log"
	"net/http"
	"time"
)

func Logger(fn http.HandlerFunc, name string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fn.ServeHTTP(w, r)
		log.Printf("%s\t%s\t%s\t%s\n",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start))

	})
}
