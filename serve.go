package main

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// func replyPing(w http.ResponseWriter, r *http.Request) {
// fmt.Fprintf(w, "Bee is alive")
// }

func (C Config) setupWebServers() chan interface{} {
	ret := make(chan interface{})

	///
	/// STARTING HTTP SERVER
	///
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	C.setupHttpHandlersx(myRouter)
	if C.HTTPport != "" {

		//myRouter.HandleFunc("/ping", replyPing)
		srv := http.Server{
			Addr:         C.HTTPport,
			ReadTimeout:  time.Duration(C.HTTPXreadTimeout) * time.Second,
			WriteTimeout: time.Duration(C.HTTPXwriteTimeout) * time.Second,
			Handler:      myRouter,
		}
		go func() {
			Info.Printf("Starting HTTP server on port %s", C.HTTPport)
			if err := srv.ListenAndServe(); err != nil {
				Error.Printf("Server HTTP failed: %s\n", err)
			}
		}()
		go func(srv http.Server, done chan interface{}) {
			Trace.Printf("Listening to HTTP server done channel")
			select {
			case <-done:
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := srv.Shutdown(ctx); err != nil {
					Error.Printf("Server HTTP unable to shutdown: %s\n", err)
				}
				Info.Printf("HTTP server stopped")
			}
		}(srv, ret)

	} else {
		Info.Println("HTTP server is disabled")
	}
	///
	/// STARTING HTTPS SERVER
	///
	if C.HTTPSport != "" {

		srvs := http.Server{
			Addr:         C.HTTPSport,
			ReadTimeout:  time.Duration(C.HTTPXreadTimeout) * time.Second,
			WriteTimeout: time.Duration(C.HTTPXwriteTimeout) * time.Second,
			Handler:      myRouter,
		}

		go func() {
			Info.Printf("Starting HTTPs server on port %s", C.HTTPSport)
			if err := srvs.ListenAndServeTLS(C.HTTPScert, C.HTTPSkey); err != nil {
				Error.Printf("Server HTTPS  failed: %s\n", err)
			}
		}()
		go func(srvs http.Server, done chan interface{}) {

			select {
			case <-done:
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := srvs.Shutdown(ctx); err != nil {
					Error.Printf("Server HTTPS unable to shutdown: %s\n", err)
				}
				Info.Printf("HTTPs server stopped")
			}
		}(srvs, ret)

	} else {
		Info.Println("HTTPs server is disabled")
	}

	return ret
}

func slowHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(2 * time.Second)
	io.WriteString(w, "I am slow!\n")
}
