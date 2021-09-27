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

	err := Db.Where(&KelpUser{ApiKey: r.PostFormValue("api_key")}).First(&user).Error

	if err != nil {
		respondError(w, 404, "user not found")
		return
	}

	//Db.Find(&projects)
	files := []KelpFile{}

	err = Db.Where(&KelpFile{UserId: int(user.ID)}).Find(&files).Error

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

	fileName := generateRandomString(6)
	fileType := trimFirstRune(filepath.Ext(handler.Filename)) // trimmed in order to preserve db structure
	orgName := strings.Split(handler.Filename, ".")[0]

	err = ioutil.WriteFile(fmt.Sprintf("./static/u/%s.%s", fileName, fileType), fileBytes, 0444)

	if err != nil {
		respondError(w, 500, "failed to write file")
		return
	}

	newFile := KelpFile{
		UserId:  int(user.ID),
		Type:    fileType,
		Name:    fileName,
		OrgName: orgName,
	}

	err = Db.Create(&newFile).Error

	if err != nil {
		respondError(w, 500, "failed to add db entry")
		os.Remove(fmt.Sprintf("./static/u/%s.%s", fileName, fileType))
		return
	}

	defer file.Close()

	respondJSON(w, 200, map[string]string{
		"filesize": byteCountSI(int64(len(fileBytes))),
		"file_id":  fileName,
		"filename": fmt.Sprintf("%s.%s", fileName, fileType),
		"link":     fmt.Sprintf("%s/u/%s.%s", r.Host, fileName, fileType),
	})
}

// delete a file
func DeleteFile(w http.ResponseWriter, r *http.Request) {
	user := KelpUser{}
	fileToDelete := KelpFile{}

	err := Db.Where(&KelpUser{ApiKey: r.PostFormValue("api_key")}).First(&user).Error

	if err != nil {
		respondError(w, 404, "user not found")
		return
	}

	err = Db.Where(&KelpFile{Name: r.PostFormValue("file_id")}).First(&fileToDelete).Error

	if err != nil {
		respondError(w, 404, "file not found")
		return
	}

	deleteFileIfExists(fmt.Sprintf("./static/u/%s.%s", fileToDelete.Name, fileToDelete.Type))

	respondJSON(w, 200, map[string]string{
		"success": "true",
	})

}

// delete all files
func DeleteAllFiles(w http.ResponseWriter, r *http.Request) {
	user := KelpUser{}
	filesToDelete := []KelpFile{}

	err := Db.Where(&KelpUser{ApiKey: r.PostFormValue("api_key")}).First(&user).Error

	if err != nil {
		respondError(w, 404, "user not found")
		return
	}

	err = Db.Where(&KelpFile{UserId: int(user.ID)}).Find(&filesToDelete).Error

	if err != nil {
		respondError(w, 500, "failed to retreive files for deletion")
		return
	}

	for _, file := range filesToDelete {
		deleteFileIfExists(fmt.Sprintf("./static/u/%s.%s", file.Name, file.Type))
	}

	Db.Where(&KelpFile{UserId: int(user.ID)}).Delete(KelpFile{})

	respondJSON(w, 200, map[string]string{
		"success": "true",
	})

}

// upload a paste
func UploadPaste(w http.ResponseWriter, r *http.Request) {
	user := KelpUser{}

	pasteID := generateRandomString(6)

	newPaste := KelpPaste{
		Name:    pasteID,
		Title:   r.PostFormValue("paste_name"),
		Content: r.PostFormValue("u_paste"),
	}

	err := Db.Where(&KelpUser{ApiKey: r.PostFormValue("api_key")}).First(&user).Error

	if err != nil {
		respondError(w, 404, "user not found")
		return
	}

	err = Db.Create(newPaste).Error

	if err != nil {
		respondError(w, 500, "failed to create paste")
		return
	}

	respondJSON(w, 200, map[string]string{
		"paste_id": pasteID,
		"link":     fmt.Sprintf("%s/p/%s", r.Host, pasteID),
		"raw_link": fmt.Sprintf("%s/p/raw/%s", r.Host, pasteID),
	})

}

// delete a paste
func DeletePaste(w http.ResponseWriter, r *http.Request) {
	user := KelpUser{}
	pasteToDelete := KelpPaste{}

	err := Db.Where(&KelpUser{ApiKey: r.PostFormValue("api_key")}).First(&user).Error

	if err != nil {
		respondError(w, 404, "user not found")
		return
	}

	err = Db.Where(&KelpPaste{Name: r.PostFormValue("paste_id")}).First(&pasteToDelete).Error

	if err != nil {
		respondError(w, 500, "paste not found")
		return
	}

	Db.Delete(&pasteToDelete)

	respondJSON(w, 200, map[string]string{
		"success": "true",
	})
}

// delete all pastes
func DeleteAllPastes(w http.ResponseWriter, r *http.Request) {
	user := KelpUser{}

	err := Db.Where(&KelpUser{ApiKey: r.PostFormValue("api_key")}).First(&user).Error

	if err != nil {
		respondError(w, 404, "user not found")
		return
	}

	err = Db.Where(&KelpPaste{UserId: int(user.ID)}).Delete(KelpPaste{}).Error

	if err != nil {
		respondError(w, 500, "failed to delete pastes")
		return
	}

	respondJSON(w, 200, map[string]string{
		"success": "true",
	})
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

// https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func byteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func deleteFileIfExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(path)
}
