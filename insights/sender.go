package insights

import (
	"log"

	"github.com/sjltaylor/stats-gopher/retry"
)

type sender struct {
	key      string
	endpoint string
}

func newSender(key, endpoint string) *sender {
	return &sender{
		key:      key,
		endpoint: endpoint,
	}
}

func (sender *sender) run(ch <-chan []interface{}) {
	for payload := range ch {
		sender.send(payload)
	}
	log.Println("insights: sender finished")
}

func (sender *sender) send(chunk []interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("insights: %v\n", err)
		}
	}()

	request := newRequest(sender.key, sender.endpoint, chunk)
	retry := retry.NewRetry()
	go handleErrors(retry)
	retry.Execute(request)
}

func handleErrors(retry *retry.Retry) {
	count := 0
	for err := range retry.Errors {
		log.Println(err)
		count++
		if count == 15 {
			retry.Stop()
		}
	}
}
