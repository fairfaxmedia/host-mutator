package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/antonosmond/host-mutator/pkg/mutator"
)

var addr = ":8443"
var sourceDomains = os.Getenv("SOURCE_DOMAINS")
var targetDomain = os.Getenv("TARGET_DOMAIN")
var certPath = os.Getenv("SSL_CERT_PATH")
var keyPath = os.Getenv("SSL_KEY_PATH")

func init() {
	if sourceDomains == "" {
		log.Fatal("SOURCE_DOMAINS environment variable must not be empty")
	}
	if targetDomain == "" {
		log.Fatal("TARGET_DOMAIN environment variable must not be empty")
	}
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/mutate", handleMutate)

	s := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1048576
	}

	log.Printf("Server listening on %s\n", addr)
	log.Fatal(s.ListenAndServeTLS(certPath, keyPath))

}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.URL, r.RemoteAddr)
	w.WriteHeader(200)
}

func handleMutate(w http.ResponseWriter, r *http.Request) {

	log.Printf("%s %s\n", r.URL, r.RemoteAddr)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	mutated, err := mutator.Mutate(body, sourceDomains, targetDomain)
	if err != nil {
		log.Println(err)
		// check if the error is due to a bad request
		if _, ok := err.(*mutator.BadRequest); ok {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mutated)

}
