package printer

import "fmt"

// Listen blocks, repeatedly printing any input data
func Listen(input <-chan interface{}) {
	for e := range input {
		fmt.Printf("received: %v\n", e)
	}
}
