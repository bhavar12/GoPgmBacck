package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type RequestBody struct {
	roman string `json:"roman`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/add", handlerPost).Methods(http.MethodPost)
	http.ListenAndServe(":8080", r)
}

func handlerPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in a post call")
	app := []map[string]int{{"i": 1}, {"ii": 2}, {"iii": 3}}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var req RequestBody
	err = json.Unmarshal(data, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newReq := req.roman
	fmt.Println(newReq)
	for _, val := range app {
		if val1, ok := val[newReq]; ok {
			fmt.Println(val1)
			fmt.Fprint(w, val1)
			return
		}
	}
}
