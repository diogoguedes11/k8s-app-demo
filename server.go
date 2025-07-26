package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var isReady = true 

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "Method not allowed")
		return
	}
	fmt.Println("got / request")
	w.Header().Set("Content-Type","text/plain")
	w.Header().Set("Cache-Control","public,max-age=300") // performance headers
	io.WriteString(w,"This is my website!\n")
}

func getHealthz(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        io.WriteString(w, "Method not allowed\n")
        return
    }
    fmt.Println("got /healthz request")
    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

    io.WriteString(w, "Healthz OK!\n")
}


func getReadiness(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        io.WriteString(w, "Method not allowed\n")
        return
    }
    fmt.Println("got /readiness request")

    w.Header().Set("Content-type","text/plain")
    w.Header().Set("Cache-control","no-cache, no-store, must-revalidate")
    if !isReady {
        w.WriteHeader(http.StatusServiceUnavailable)
        io.WriteString(w, "Not Ready\n")
        return
    }
    w.WriteHeader(http.StatusOK)
    io.WriteString(w, "Ready\n")
}

func setReadinessFail(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        io.WriteString(w, "Method not allowed\n")
        return
    }
    fmt.Println("Setting readiness to FAIL")
    isReady = false
    
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    io.WriteString(w, "Readiness set to FAIL\n")
}

func setReadinessOk(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        io.WriteString(w, "Method not allowed\n")
        return
    }
    fmt.Println("Setting readiness to OK")
    isReady = true
    
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    io.WriteString(w, "Readiness set to OK\n")
}

func main() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/healthz", getHealthz)
	http.HandleFunc("/readiness", getReadiness)
	http.HandleFunc("/readiness/fail", setReadinessFail)
	http.HandleFunc("/readiness/ok", setReadinessOk)
	port := getPort()
	fmt.Printf("Web server started at port %s\n",port)
	server := &http.Server{ // clients cant block our server forever :) 
		Addr: port,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second, 
		IdleTimeout: 15 * time.Second, // closes idle connections

	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
