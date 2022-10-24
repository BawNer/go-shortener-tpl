package app

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func HandlerRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	switch r.Method {
	case "POST":
		URL, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		timestamp := time.Now()
		shr := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s %d", string(URL)[2:6], timestamp.Second())))
		resStruct := DB{len(LocalDB), string(URL), shr}
		SaveDB(resStruct)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", shr)))
	case "GET":
		id := strings.TrimPrefix(r.RequestURI, "/")
		columns, err := FindByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Location", columns.URL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(fmt.Sprintf("http://localhosy:8080/%s", columns.URLShort)))
	default:
		http.Error(w, "Bad Gateway", http.StatusMethodNotAllowed)
		return
	}
}
