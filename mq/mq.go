package mq

import "log"

var receivers = make([]chan interface{}, 0)

// Send data to any listening receivers
// If the receivers buffered channel is full, the events are dropped
func Send(data interface{}) {
	var array []interface{}
	var isArray bool

	if array, isArray = data.([]interface{}); !isArray {
		array = []interface{}{data}
	}
	for _, receiver := range receivers {
		for _, datum := range array {
			select {
			case receiver <- datum:
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
