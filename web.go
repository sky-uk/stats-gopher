package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sjltaylor/stats-gopher/insights"
	"github.com/sjltaylor/stats-gopher/mq"
	"github.com/sjltaylor/stats-gopher/web"
	"github.com/yvasiyarov/gorelic"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	port := port()

	newRelicKey := newRelicKey()
	newRelicEndpoint := newRelicEndpoint()

	newRelicMonitoring()

	go insights.Listen(newRelicKey, newRelicEndpoint, mq.Channel())

	fmt.Printf("Stats Gopher PORT=%s\n", port)
	web.Start(fmt.Sprintf(":%s", port))
}

func newRelicMonitoring() {
	key := os.Getenv("NEWRELIC_MONITORING_KEY")

	if key == "" {
		return
	}

	agent := gorelic.NewAgent()
	agent.Verbose = true
	agent.NewrelicLicense = key
	agent.Run()
}

func port() string {
	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "80"
	}

	return port
}

func newRelicKey() string {
	key := os.Getenv("NEWRELIC_KEY")

	if len(key) == 0 {
		panic("NEWRELIC_KEY not set")
	}

	return key
}

func newRelicEndpoint() string {
	endpoint := os.Getenv("NEWRELIC_ENDPOINT")

	if len(endpoint) == 0 {
		panic("NEWRELIC_ENDPOINT not set")
	}

	return endpoint
}
