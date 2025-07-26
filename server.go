package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var isReady = true // simple global variable

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got / request")
	io.WriteString(w, "This is my website!\n")
}

func getHealthz(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /healthz request")
	io.WriteString(w, "Healthz OK!\n")
}

func getReadiness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /readiness request")
	if !isReady {
		w.WriteHeader(http.StatusServiceUnavailable)
		io.WriteString(w, "Not Ready\n")
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Ready\n")
}

func setReadinessFail(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /readiness/fail request")
	isReady = false
	w.WriteHeader(http.StatusServiceUnavailable)
	io.WriteString(w, "Readiness set to false\n")
}

func setReadinessOk(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /readiness/ok request")
	isReady = true
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Readiness set to true\n")
}

func main() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/healthz", getHealthz)
	http.HandleFunc("/readiness", getReadiness)
	http.HandleFunc("/readiness/fail", setReadinessFail)
	http.HandleFunc("/readiness/ok", setReadinessOk)

	fmt.Println("Web server started at port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
