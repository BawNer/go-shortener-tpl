package app

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type postRequest struct {
	Url string `json:"url"`
}

func HandlerRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		URL, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		timestamp := time.Now()
		shr := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s %d", string(URL)[2:6], timestamp.Second())))
		resStruct := DB{len(LocalDB), string(URL), fmt.Sprintf("http://localhost:8080/%s", shr)}
		SaveDB(resStruct)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", shr)))
	case "GET":
		re, _ := regexp.Compile(`\d`)
		id, _ := strconv.Atoi(string(re.Find([]byte(r.RequestURI))))
		w.Header().Set("Content-Type", "text/plain")
		columns, err := FindByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Header().Set("Location", columns.URLShort)
		w.Write([]byte(columns.URLShort))
	default:
		http.Error(w, "Bad Gateway", http.StatusMethodNotAllowed)
		return
	}
}
