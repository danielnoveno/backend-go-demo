package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const defaultBaseURL = "http://127.0.0.1:8080/api"

type SensorData struct {
	ID        uint      `json:"id"`
	Status    string    `json:"status"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func GetLatestSensor() (*SensorData, error) {
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(baseURL + "/sensor/latest")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %d", resp.StatusCode)
	}

	var sensor SensorData
	if err := json.NewDecoder(resp.Body).Decode(&sensor); err != nil {
		return nil, err
	}

	return &sensor, nil
}
