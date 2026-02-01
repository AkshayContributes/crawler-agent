package brain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OllamaStrategy struct {
	// Implementation details would go here
	Model string
	URL   string
}

func NewOllamaStrategy(model string) *OllamaStrategy {
	return &OllamaStrategy{
		Model: model,
		URL:   "http://localhost:11434/api/generate",
	}
}

func (o *OllamaStrategy) Summarize(text string) (string, error) {

	if len(text) > 6000 {
		text = text[:6000] + "..."
	}

	reqBody := map[string]interface{}{
		"model":  o.Model,
		"prompt": fmt.Sprintf("Summarize this into 3 technical bullet points: %s", text),
		"stream": false,
	}

	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(o.URL, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	var result struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}

func (o *OllamaStrategy) Close() error {
	// No resources to clean up in this implementation
	return nil
}
