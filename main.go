package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sjltaylor/stats-gopher/insights"
	"github.com/sjltaylor/stats-gopher/mq"
	"github.com/sjltaylor/stats-gopher/printer"
	"github.com/sjltaylor/stats-gopher/web"
	"github.com/yvasiyarov/gorelic"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	newRelicMonitoring()
	startInsights()
	startPrinter()
	startWebServer()
}

func startInsights() {
	if ep := os.Getenv("NEW_RELIC_INSIGHTS_ENDPOINT"); ep != "" {
		key := os.Getenv("NEW_RELIC_INSIGHTS_KEY")

		if len(key) == 0 {
			panic("NEW_RELIC_INSIGHTS_KEY not set")
		}

		go insights.Listen(key, ep, mq.Channel())
	} else {
		fmt.Println("Insights listener not started; NEW_RELIC_INSIGHTS_ENDPOINT is unset")
	}
}

func startPrinter() {
	if os.Getenv("STDOUT_LISTENER") == "1" {
		go printer.Listen(mq.Channel())
	}
}

func newRelicMonitoring() {
	key := os.Getenv("NEW_RELIC_LICENSE_KEY")

	if key == "" {
		return
	}

	agent := gorelic.NewAgent()
	agent.Verbose = true
	agent.NewrelicLicense = key
	agent.NewrelicName = "stats-gopher"
	agent.GCPollInterval = 60
	agent.Run()
}

func startWebServer() {
	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "80"
	}

	fmt.Printf("Stats Gopher PORT=%s\n", port)
	web.Start(fmt.Sprintf(":%s", port), os.Getenv("NEW_RELIC_LICENSE_KEY"))
}
