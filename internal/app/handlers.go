package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func HandlerRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		URL, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		hash := md5.New()
		shr := hash.Sum(URL)
		resStruct := DB{len(LocalDB), string(URL), hex.EncodeToString(shr)}
		SaveDB(resStruct)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", hex.EncodeToString(shr))))
	case "GET":
		id := strings.TrimPrefix(r.RequestURI, "/")
		columns, err := FindByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		fmt.Println(columns.URL)
		w.Header().Set("Location", columns.URL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		http.Error(w, "Bad Gateway", http.StatusMethodNotAllowed)
		return
	}
}
