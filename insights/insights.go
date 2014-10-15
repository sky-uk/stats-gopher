package insights

import "github.com/sjltaylor/stats-gopher/mq"

// Listen blocks, repeatedly sending any input to insights
func Listen(key, endpoint string, input <-chan interface{}) {
	channelReader := mq.NewChannelReader(input, 16)
	sender := newSender(key, endpoint)

	go sender.run(channelReader.Out)

	channelReader.Read()
}
