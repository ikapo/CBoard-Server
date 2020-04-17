package main

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"pages"
	"time"
)

// exitOnError logs an error then exists
func exitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// waitForAPI waits for the API to be ready for connections
func waitForAPI() {
	host := "api:8000"
	connected := false
	attempts := 0
	for attempts < 10 && !connected {
		// Trying to connect to API
		_, err := net.DialTimeout("tcp", host, time.Second*5)

		if err != nil {
			log.Println("API not ready for connections yet, sleeping 5s")
			attempts++
			time.Sleep(time.Second * 5)
		} else {
			connected = true
		}
	}

	if !connected {
		log.Fatal("Failed to connect to API, shutting down...")
	}

}

func main() {

	waitForAPI()

	// Time duration configs for timeout
	var wait time.Duration
	usage := "time server waits for existing connections to finish"
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, usage)
	flag.Parse()

	// Defining mux router
	r := mux.NewRouter()
	// Initializing templates
	pages.Init()

	r.HandleFunc(
		"/", pages.Index,
	)

	// Routing for a board, can only be 1-4
	r.HandleFunc(
		"/{board:1|2|3|4}/", pages.Board,
	)

	// Routing for seeing a post,
	// ID can only be a positive int
	r.HandleFunc(
		"/post/{id:[0-9]+}", pages.Post,
	)

	// Static file serving
	fs := http.FileServer(http.Dir("./static"))
	staticHandler := http.StripPrefix("/static/", fs)
	r.PathPrefix("/static/").Handler(staticHandler)

	srv := &http.Server{
		Addr: "0.0.0.0:80",
		// Setting timeouts to mitigate Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	//Starting server as a goroutine
	go func() {
		log.Println("Ready to receive connections")
		err := srv.ListenAndServe()
		exitOnError(err)
	}()

	c := make(chan os.Signal, 1)
	// Accepting graceful shutdowns when quit via SIGINT (Ctrl+C)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	log.Println("Shutting down")
	os.Exit(0)
}
