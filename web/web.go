package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/jingweno/negroni-gorelic"
	"github.com/sjltaylor/stats-gopher/mq"
	"github.com/sjltaylor/stats-gopher/presence"
)

type OneLineLogger struct {
	*log.Logger
}

func newLogger() *OneLineLogger {
	return &OneLineLogger{log.New(os.Stdout, "[stats-gopher]\t", 0)}
}

func (l *OneLineLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)
	res := rw.(negroni.ResponseWriter)
	requestId := r.Header.Get("X-Request-ID")
	l.Printf("%s\t%s\t%s\t%v\t%s\t%v", r.Method, r.URL.Path, requestId, res.Status(), http.StatusText(res.Status()), time.Since(start))
}

func newHttpServer() *negroni.Negroni {
	return negroni.New(negroni.NewRecovery(), newLogger(), negroni.NewStatic(http.Dir("public")))
}

var monitors = presence.NewMonitorPool(map[string]time.Duration{
	"heartbeat":     time.Second * 45,
	"user-activity": time.Minute * 45,
})

// Start the web server
// the key is the api key for new relic monitoring
func Start(bind, key string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/gopher/", gopherEndpoint)
	mux.HandleFunc("/presence/", presenceEndpoint)
	mux.HandleFunc("/", helloEndpoint)

	n := newHttpServer()
	n.UseHandler(mux)

	n.Use(negroni.NewRecovery())

	if key != "" {
		n.Use(negronigorelic.New(key, "stats-gopher", true))
	}

	go func() {
		for {
			n := <-monitors.C
			event := map[string]interface{}{
				"eventType":        "UserAction",
				"key":              n.Key,
				"code":             n.Code,
				"start":            n.Start,
				"lastNotification": n.LastNotification,
				"end":              n.End,
				"duration":         int64(n.Duration) / 1e9,
				"wait":             int64(n.Wait) / 1e9,
			}
			mq.Send(event)
		}
	}()

	n.Run(bind)
}

func helloEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}

	respond("hello, send me something", w, r)
}

func gopherEndpoint(w http.ResponseWriter, r *http.Request) {
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

	go mq.Send(data)

	respond("nom nom nom", w, r)
}

func presenceEndpoint(w http.ResponseWriter, r *http.Request) {
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
	var notifications []presence.Notification

	if blob, err = ioutil.ReadAll(r.Body); err != nil {
		log.Println(fmt.Sprintf("web: error reading heartbeat body: %s", err))
		w.WriteHeader(500)
		return
	}

	if err = json.Unmarshal(blob, &notifications); err != nil {
		log.Printf("web: error parsing heartbeat json (%s): %s\n", err, string(blob))
		w.WriteHeader(500)
		return
	}

	for _, notification := range notifications {
		go monitors.Notify(&notification)
	}
}

func respond(message string, w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprint(w, fmt.Sprintf("%s\n", message)); err != nil {
		w.WriteHeader(500)
		log.Printf(fmt.Sprintf("web: couldn't write to the response: '%s', %s\n", r.URL.String(), err))
	}
}
