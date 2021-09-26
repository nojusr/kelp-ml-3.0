package kelp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
	"unsafe"
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

	err := Db.Where(&KelpUser{ApiKey: r.PostFormValue("api_key")}).First(&user).Error

	if err != nil {
		respondError(w, 404, "user not found")
		return
	}

	err = r.ParseMultipartForm(100 << 20)

	if err != nil {
		respondError(w, 500, "file too large (limit: 100mb)")
		return
	}

	file, handler, err := r.FormFile("u_file")

	if err != nil {
		respondError(w, 500, "file not found")
		return
	}

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		respondError(w, 500, "failed to read file")
		return
	}

	filename := generateRandomString(6)
	fileType := trimFirstRune(filepath.Ext(handler.Filename)) // trimmed in order to preserve db structure
	orgName := strings.Split(handler.Filename, ".")[0]

	err = ioutil.WriteFile(fmt.Sprintf("./static/u/%s.%s", filename, fileType), fileBytes, 0444)

	if err != nil {
		respondError(w, 500, "failed to write file")
		return
	}

	newFile := KelpFile{
		UserId:  int(user.ID),
		Type:    fileType,
		Name:    filename,
		OrgName: orgName,
	}

	err = Db.Create(&newFile).Error

	if err != nil {
		respondError(w, 500, "failed to add db entry")
		os.Remove(fmt.Sprintf("./static/u/%s.%s", filename, fileType))
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

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
// generates random str of length n, and is apparantly super fast. i dont care enough to go in-depth about this
// so i'm keeping this as-is.
func generateRandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
