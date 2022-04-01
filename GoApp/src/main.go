package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

const (
	uid = "uid"
)

type ContextKey string
type ContextKey1 string

const ContextUserKey ContextKey = "user"
const ContextUserKey1 ContextKey1 = "user"

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}
type 

func 

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareOne")

		ctx := context.WithValue(r.Context(), ContextUserKey, "the user")
		ctx1 := context.WithValue(r.Context(), ContextUserKey1, "the user second context")
		r = r.WithContext(ctx1)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
		log.Println("Executing middlewareOne again")
	})
}

func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareTwo")
		if r.URL.Path != "/" {
			return
		}
		next.ServeHTTP(w, r)
		log.Println("Executing middlewareTwo again")
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing finalHandler")
	uid := r.Context().Value(ContextUserKey)
	uid1 := r.Context().Value(ContextUserKey1)
	//uid = nil
	fmt.Println(uid.(string))
	fmt.Println(uid1.(string))
	w.WriteHeader(http.StatusBadRequest)
	fmt.Printf("%+v", w)
	lrw := NewLoggingResponseWriter(w)
	statusCode := lrw.statusCode
	fmt.Println(statusCode)
}
func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}
func main() {
	finalHandler := http.HandlerFunc(final)

	http.Handle("/", middlewareOne(middlewareTwo(finalHandler)))
	http.ListenAndServe(":3000", nil)
}
