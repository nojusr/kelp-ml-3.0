package kelp

import (
	"net/http"
)

// get global site stats.
// params: none
func fetchSiteStats(w http.ResponseWriter, r *http.Request) {

}

// get user stats
// params:
// 'api_key': user's API key
func fetchUserStats(w http.ResponseWriter, r *http.Request) {

}

// get user files
// params:
// 'api_key': user's API key
func fetchUserUploads(w http.ResponseWriter, r *http.Request) {

}

// get user pastes
func fetchUserPastes(w http.ResponseWriter, r *http.Request) {

}

// upload a file
// params:
// 'api_key': user's API key
// 'u_file': the file to be uploaded
func uploadFile(w http.ResponseWriter, r *http.Request) {

}

// delete a file
func deleteFile(w http.ResponseWriter, r *http.Request) {

}

// delete all files
func deleteAllFiles(w http.ResponseWriter, r *http.Request) {

}

// upload a paste
func uploadPaste(w http.ResponseWriter, r *http.Request) {

}

// delete a paste
func deletePaste(w http.ResponseWriter, r *http.Request) {

}

// delete all pastes
func deleteAllPastes(w http.ResponseWriter, r *http.Request) {

}
