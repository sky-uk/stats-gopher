package mq

import "log"

var receivers = make([]chan interface{}, 0)

// Send data to any listening receivers
// If the receivers buffered channel is full, the events are dropped
func Send(sid string, data interface{}) {
	var array []interface{}
	var isArray bool

	if array, isArray = data.([]interface{}); !isArray {
		array = []interface{}{data}
	}

	for _, datum := range array {
		m, ok := datum.(map[string]interface{})

		if !ok {
			log.Printf("mq: expected data to be a key-value map, but got: %v'n", datum)
			continue
		}

		m["sid"] = sid

		for _, receiver := range receivers {
			select {
			case receiver <- m:
			default:
				log.Println("Buffer full, dropping event")
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
