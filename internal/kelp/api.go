package kelp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// get global site stats.
// params: none
func FetchSiteStats(w http.ResponseWriter, r *http.Request) {

}

// get user stats
// params:
// 'api_key': user's API key
func FetchUserStats(w http.ResponseWriter, r *http.Request) {

}

// get user files
// params:
// 'api_key': user's API key
func FetchUserUploads(w http.ResponseWriter, r *http.Request) {
	user := KelpUser{}

	err := Db.Where(&KelpUser{ApiKey: r.PostFormValue("api_key")}).First(&user)

	if err != nil {
		respondError(w, 404, "user not found")
		return
	}

	//Db.Find(&projects)
	files := []KelpFile{}

	err = Db.Where(&KelpFile{UserId: int(user.ID)}).Find(&files)

	if err != nil {
		respondError(w, 500, "failed to retreive files")
		return
	}

	respondJSON(w, 200, files)
}

// get user pastes
func FetchUserPastes(w http.ResponseWriter, r *http.Request) {
	user := KelpUser{}

	err := Db.Where(&KelpUser{ApiKey: r.PostFormValue("api_key")}).First(&user).Error

	log.Println(fmt.Sprintf("user chk: %s, id: %d", user.Username, user.ID))

	if err != nil {
		respondError(w, 404, "user not found")
		return
	}

	//Db.Find(&projects)
	pastes := []KelpPaste{}
	err = Db.Where(&KelpPaste{UserId: int(user.ID)}).Find(&pastes).Error

	if err != nil {
		respondError(w, 500, "failed to retreive pastes")
		return
	}

	respondJSON(w, 200, pastes)
}

// upload a file
// params:
// 'api_key': user's API key
// 'u_file': the file to be uploaded
func UploadFile(w http.ResponseWriter, r *http.Request) {
	user := KelpUser{}

	Db.Where(&KelpUser{ApiKey: r.PostFormValue("api_key")}).First(&user)

	r.ParseMultipartForm(100 << 20)

	file, _, err := r.FormFile("u_file")

	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
}

// delete a file
func DeleteFile(w http.ResponseWriter, r *http.Request) {

}

// delete all files
func DeleteAllFiles(w http.ResponseWriter, r *http.Request) {

}

// upload a paste
func UploadPaste(w http.ResponseWriter, r *http.Request) {

}

// delete a paste
func DeletePaste(w http.ResponseWriter, r *http.Request) {

}

// delete all pastes
func DeleteAllPastes(w http.ResponseWriter, r *http.Request) {

}

// common call to return JSON
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// common call to return a JSON error
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}
