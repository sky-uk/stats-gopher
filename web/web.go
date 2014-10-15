package web

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/sjltaylor/stats-gopher/mq"
)

func Start(bind string) {
	http.HandleFunc("/gopher/", gopher)
	http.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(bind, nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}

	respond("hello, send me something", w, r)
}

func sid(w http.ResponseWriter, r *http.Request) string {
	var cookie *http.Cookie
	var err error

	if cookie, err = r.Cookie("sg-sid"); err != nil {
		hash := sha256.New()
		hash.Write([]byte(r.RemoteAddr))
		hash.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
		sum := hash.Sum(nil)
		sid := base64.StdEncoding.EncodeToString(sum)
		cookie = &http.Cookie{Name: "sg-sid", Value: sid}
		http.SetCookie(w, cookie)
	}

	return cookie.Value
}

func gopher(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	var blob []byte
	var err error
	var data interface{}

	if blob, err = ioutil.ReadAll(r.Body); err != nil {
		log.Println(fmt.Sprintf("web: error reading body: %s", err))
		w.WriteHeader(500)
		return
	}

	if err = json.Unmarshal(blob, &data); err != nil {
		log.Printf("web: error parsing json (%s): %s\n", err, string(blob))
		w.WriteHeader(500)
		return
	}

	go mq.Send(sid(w, r), data)

	respond("nom nom nom", w, r)
}

func respond(message string, w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprint(w, fmt.Sprintf("%s\n", message)); err != nil {
		w.WriteHeader(500)
		log.Printf(fmt.Sprintf("web: couldn't write to the response: '%s', %s\n", r.URL.String(), err))
	}
}
