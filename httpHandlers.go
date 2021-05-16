package main

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

func (C Config) setupHttpHandlersx(myRouter *mux.Router) {
	myRouter.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { C.replyPing(w, r) })
	myRouter.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) { C.update(w, r) })
	//myRouter.HandleFunc("/updatexmlrpc", func(w http.ResponseWriter, r *http.Request) { C.replyPing(w, r) })

}

func (C Config) replyPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bee is alive")
	return
}

func checkKey(w http.ResponseWriter, r *http.Request) {

	return
}

func (C Config) update(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()        // Parses the request body
	x := r.Form.Get("G") // x will be "" if parameter is not set
	fmt.Println(x)

	var zzz = beeObj{UID: "MYSELF", AuthKey: "AZZEERT"}

	t := reflect.TypeOf(zzz)

	for i := 0; i < t.NumField(); i++ {
		if value, ok := t.Field(i).Tag.Lookup("private"); ok {
			fmt.Println(value)
		} else {
			fmt.Print(" NOPE ")
		}
	}

	return
}
