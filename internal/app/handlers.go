package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
		decoder := json.NewDecoder(r.Body)
		var t postRequest
		err := decoder.Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		timestamp := time.Now()
		shr := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s %d", t.Url[2:6], timestamp.Second())))
		resStruct := DB{len(LocalDB), t.Url, fmt.Sprintf("http://localhost:8080/%s", shr)}
		written, _ := SaveDB(resStruct)
		res, err := json.Marshal(written)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	case "GET":
		re, _ := regexp.Compile(`\d`)
		id, _ := strconv.Atoi(string(re.Find([]byte(r.RequestURI))))
		w.Header().Set("Content-Type", "application/json")
		columns, err := FindById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(columns)
		w.Write(response)
	default:
		http.Error(w, "Bad Gateway", http.StatusMethodNotAllowed)
		return
	}
}
