package main

import (
	"internal/kelp"
	"log"
	"net/http"
	"os"

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
	kelp.InitializeDatabase()

	// router
	r := mux.NewRouter()

	// API handlers
	r.HandleFunc("/api/fetch/stats", kelp.FetchSiteStats)
	r.HandleFunc("/api/fetch/user", kelp.FetchUserStats)
	r.HandleFunc("/api/fetch/user/files", kelp.FetchUserUploads)
	r.HandleFunc("/api/fetch/user/pastes", kelp.FetchUserPastes)
	r.HandleFunc("/api/upload", kelp.UploadFile)
	r.HandleFunc("/api/upload/delete", kelp.DeleteFile)
	r.HandleFunc("/api/upload/delete/all", kelp.DeleteAllFiles)
	r.HandleFunc("/api/paste", kelp.UploadPaste)
	r.HandleFunc("/api/paste/delete", kelp.DeletePaste)
	r.HandleFunc("/api/paste/delete/all", kelp.DeleteAllPastes)

	// frontend handlers
	r.HandleFunc("/", kelp.RootHandler)

	// file handler
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	err := http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal(err)
	}
}
