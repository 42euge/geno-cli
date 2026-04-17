package ollama

import (
	"bufio"
	"encoding/json"
	"net/http"
)

func streamResponse(resp *http.Response, ch chan<- ChatResponse) {
	defer resp.Body.Close()
	defer close(ch)

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var chatResp ChatResponse
		if err := json.Unmarshal(line, &chatResp); err != nil {
			continue
		}

		ch <- chatResp

		if chatResp.Done {
			return
		}
	}
}
