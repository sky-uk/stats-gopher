package mq

import (
	"fmt"
	"log"
)

var receivers = make([]chan interface{}, 0)

// Send data to any listening receivers
// If the receivers buffered channel is full, the events are dropped
func Send(data interface{}) {
	var array []interface{}
	var isArray bool

	if array, isArray = data.([]interface{}); !isArray {
		array = []interface{}{data}
	}

	for _, datum := range array {
		m, ok := datum.(map[string]interface{})

		if !ok {
			log.Printf(fmt.Sprintf("mq: expected data to be a key-value map, but got: %v", datum))
			continue
		}

		if sid, ok := m["sid"]; !ok || sid == "" {
			log.Println(fmt.Sprintf("mq: event dropped because it has no session id: %v", datum))
			continue
		}

		for _, receiver := range receivers {
			// TODO should really pass a copy of the event into the receiver so that
			// mutations aren't visible between receivers. For now, nothing is expected
			// to mutate the data
			select {
			case receiver <- m:
			default:
				log.Println("mq: buffer full, dropping event")
			}
		}
	}
}

// Channel returns a new channed into which all received data will be sent
func Channel() <-chan interface{} {
	channel := make(chan interface{}, 65536)
	receivers = append(receivers, channel)
	return channel
}
