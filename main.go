package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
	m "github.com/antonosmond/ingress-mutating-webhook/pkg/mutate"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", handleHealth)
	mux.HandleFunc("/mutate", handleMutate)

	s := &http.Server{
		Addr:           ":8443",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1048576
	}

	log.Fatal(s.ListenAndServe())

}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func handleMutate(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// mutate the request
	mutated, err := m.Mutate(body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// and write it back
	w.WriteHeader(http.StatusOK)
	w.Write(mutated)
}
