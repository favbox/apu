package api

import (
	"log"
	"net/http"
)

func MyMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("running before MyMiddleware")
	next(rw, r)
	log.Println("running after MyMiddleware")
}
