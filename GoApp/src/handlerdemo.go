package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// type hotdog int

// func (h hotdog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(w, "Any code")
// }
const CPartnerID = "PartnerID"

type IHandler struct{}
type SHandler struct{}

func main() {
	// var h hotdog
	// http.Handle("/", h)
	// http.ListenAndServe(":8080", nil)
	// t := func() map[string]int {
	// 	return map[string]int{"text": 1, "foo": 2}

	// }()
	// fmt.Println(t)
	router := mux.NewRouter()
	//router.HandleFunc("/test", signup).Methods("GET")
	//router.HandleFunc("/signin", signup).Methods("GET")
	router.HandleFunc("/protected", authMiddleware(protectedEndpoint)).Methods("GET")
	http.ListenAndServe(":8080", router)
	// go func() {
	// 	log.Fatal(http.ListenAndServe(":8080", IHandler{}))
	// }()
	// ssserv := http.Server{
	// 	Addr:    ":8181",
	// 	Handler: SHandler{},
	// }
	// ssserv.ListenAndServe()

}
func (ih IHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//pid := mux.Vars(r)[CPartnerID]
	//fmt.Println("signup invoked", r.URL.EscapedPath(), pid)
	//http.Redirect(w, r, "http://localhost:8181"+r.URL.Path, 302)
	http.Redirect(w, r, "https://integration.agent.exec.itsupport247.net/agent/v1/download/swagger.yaml", 302)
	fmt.Println(w)
	fmt.Println(r)
}
func (sh SHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "signin invoked")
}
func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("protected invoked")
}
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("authMiddleware invoked")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		private := q.Get("itsplatform_private")
		partner := q.Get("itsplatform_partner")
		client := q.Get("itsplatform_client")
		log.Println(private, partner, client)
		next.ServeHTTP(w, r)
		log.Println("Executing middlewareTwo again")
	})
}

func newMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(rw, r)
	})
}
