package main

import (
	"log"
	"net/http"
	"os"

	"kelp"

	"github.com/gorilla/mux"
)

func init_log() {
	log_file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(log_file)
}

func main() {
	init_log()

	// router
	r := mux.NewRouter()

	// frontend handler
	r.HandleFunc("/", kelp.RootHandler)

	// file handler
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	err := http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal(err)
	}
	//log.Println("Hello, world!")
	//fmt.Println("Hello, Arch!")
}
