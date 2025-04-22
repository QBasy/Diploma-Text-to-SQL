package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GenerateSQLWithGroc(client *http.Client, prompt, systemMessage string, model string, apiKey string, baseURL string) (string, error) {
	requestBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": systemMessage},
			{"role": "user", "content": prompt},
		},
	}
	log.Printf(prompt)
	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: received non-200 status code %d", resp.StatusCode)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices returned by Groq API")
	}

	sqlQuery, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	if !ok || sqlQuery == "" {
		return "", fmt.Errorf("empty SQL response")
	}

	log.Printf("Generated SQL: %s", sqlQuery)

	return sqlQuery, nil
}
