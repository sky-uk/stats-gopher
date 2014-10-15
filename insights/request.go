package insights

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type request struct {
	key      string
	endpoint string
	payload  []byte
}

func newRequest(key, endpoint string, data []interface{}) *request {
	var err error
	var payload []byte

	if payload, err = json.Marshal(data); err != nil {
		panic(fmt.Sprintf("insights :could not encode data to json (%v): '%v'\n", err, data))
	}

	return &request{key: key, endpoint: endpoint, payload: payload}
}

func (request *request) Try() error {
	var err error
	var httpRequest *http.Request

	if httpRequest, err = http.NewRequest("POST", request.endpoint, bytes.NewBuffer(request.payload)); err != nil {
		panic(fmt.Sprintf("insights: could not create request (%v)\n", err))
	}

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("X-Insert-Key", request.key)

	var response *http.Response

	var client = &http.Client{}

	if response, err = client.Do(httpRequest); err != nil {
		return fmt.Errorf("error sending request (%v)\n", err)
	}

	status := response.StatusCode

	if status < 200 { // 0-199
		return fmt.Errorf("unexpected response (%d)\n", status)
	}

	if status < 300 { // 200-299
		return nil
	}

	if status < 400 { // 300-399
		return fmt.Errorf("unexpected redirect response (%d)\n", status)
	}

	if status < 500 { // 400-499
		return fmt.Errorf("client error (%d)\n", status)
	}

	if status < 600 { // 500-599
		return fmt.Errorf("server error (%d)\n", status)
	}

	// 600+
	return fmt.Errorf("unexpected response (%d)\n", status)
}
