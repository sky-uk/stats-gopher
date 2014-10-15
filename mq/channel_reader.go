package mq

// ChannelReader reads available input into chunks
type ChannelReader struct {
	in      <-chan interface{}
	buffer  []interface{}
	Out     <-chan []interface{}
	out     chan []interface{}
	maxSize int
}

// NewChannelReader initializes a new reader
// maximum sizes of less than two result in an effective maximum size of two
func NewChannelReader(ch <-chan interface{}, maxSize int) *ChannelReader {
	out := make(chan []interface{}, 5)
	return &ChannelReader{
		in:      ch,
		buffer:  make([]interface{}, 0, 128),
		Out:     out,
		out:     out,
		maxSize: maxSize,
	}
}

// Read infinitely drains available input into buffers which are
// passed into the output channel
func (r *ChannelReader) Read() []interface{} {

	for {
		r.buffer = append(r.buffer, <-r.in)

		for stop := false; !stop; {
			select {
			case e := <-r.in:
				r.buffer = append(r.buffer, e)
				if len(r.buffer) >= r.maxSize {
					stop = true
				}
			default:
				stop = true
			}
		}

		r.flush()
	}
}

func (r *ChannelReader) flush() {
	dst := make([]interface{}, len(r.buffer))

	copy(dst, r.buffer)

	r.buffer = r.buffer[:0]
	r.out <- dst
}
