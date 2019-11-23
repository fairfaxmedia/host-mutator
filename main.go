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
var certPath = os.Getenv("SSL_CERT_PATH")
var keyPath = os.Getenv("SSL_KEY_PATH")

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

	mutated, err := mutator.Mutate(body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mutated)

}
