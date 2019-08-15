package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

type cconfig struct {
	Port    int    `json:"port"`
	Secret  string `json:"secret"`
	Protect bool   `json:"protect"`
}

func main() {
	box := packr.New("Configs", "./config")

	cconf, err := box.Find("__cconfig.json")
	if err != nil {
		panic(err)
	}

	var cc cconfig
	err = json.Unmarshal(cconf, &cc)
	if err != nil {
		panic(err)
	}

	if cc.Protect {
		var evalSecret string
		fmt.Println("Enter secret below to start: ")
		fmt.Scanln(&evalSecret)
		if evalSecret == cc.Secret {
			startServer(cc, box)
		} else {
			fmt.Println("Secret did not match with the one that it was built with.")
		}
		return
	}

	startServer(cc, box)
}

func startServer(c cconfig, box *packr.Box) {
	http.HandleFunc("/", authenticated(c, serveConfig(box)))
	fmt.Println(fmt.Sprintf("Circulator starting at port: %d", c.Port))
	panic(http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil))
}

func authenticated(c cconfig, aHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		secret := strings.Split(bearer, " ")[1]
		if secret != "" && c.Secret == secret {
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

		tempIface, err := accessConfig(application, box)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
			return
		}

		// json marshaling a non json is also good in this case
		// because if it is an interface that we try to marshal as type
		// we need to use gob.Encoder to convert into bytes
		// in this case json.Marshal does a pretty good job of returning
		// bytes of primitive type without the need to cast it
		resp, err := json.Marshal(tempIface)

		// TODO: this error thing is repeating
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
			return
		}

		if _, ok := tempIface.(map[string]interface{}); ok {
			w.Header().Set("Content-Type", "application/json")
		}
		w.Write(resp)
	}
}

func accessConfig(accessor string, box *packr.Box) (interface{}, error) {
	var fields = []string{}

	if strings.Contains(accessor, ".") {
		fields = strings.Split(accessor, ".")
	} else {
		fields = append(fields, accessor)
	}

	jsonString, err := box.FindString(fields[0] + ".json5")
	if err != nil {
		return nil, err
	}

	var tempIface map[string]interface{}
	err = json5.Unmarshal([]byte(jsonString), &tempIface)
	if err != nil {
		return nil, err
	}

	// then we have entered the arena of dot notations
	if len(fields) > 1 {
		var finalValue interface{}
		for _, f := range fields[1:] {
			if finalValue == nil {
				if tempIface[f] == nil {
					return nil, noSuchKeyErr(f, accessor)
				}
				finalValue = tempIface[f]
			} else {
				interm, ok := finalValue.(map[string]interface{})
				if !ok {
					return nil, noSuchKeyErr(f, accessor)
				}
				if interm[f] == nil {
					return nil, noSuchKeyErr(f, accessor)
				}
				finalValue = interm[f]
			}
		}
		return finalValue, nil
	}

	return tempIface, nil
}

func noSuchKeyErr(key, acc string) error {
	return errors.New("no such key: " + key + " for notation: " + acc)
}
