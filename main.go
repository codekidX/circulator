package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

type cconfig struct {
	Port   int    `json:"port"`
	Secret string `json:"secret"`
}

func main() {
	box := packr.New("Configs", "./config")

	cconfigStr, err := box.FindString("__cconfig.json")
	if err != nil {
		panic(err)
	}

	var cc cconfig
	err = json.Unmarshal([]byte(cconfigStr), &cc)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", authenticated(cc.Secret, serveConfig(box)))
	fmt.Println(fmt.Sprintf("Circulator starting at port: %d", cc.Port))
	panic(http.ListenAndServe(fmt.Sprintf(":%d", cc.Port), nil))
}

func authenticated(evalSecret string, aHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		secret := strings.Split(bearer, " ")[1]
		if secret != "" && evalSecret == secret {
			aHandler(w, r)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
	}
}

func serveConfig(box *packr.Box) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		application := r.URL.Query().Get("app")
		if application == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		jsonString, err := box.FindString(application + ".json5")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var tempIface interface{}
		err = json5.Unmarshal([]byte(jsonString), &tempIface)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}

		resp, err := json.Marshal(tempIface)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}
