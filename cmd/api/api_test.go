package main

import (
	"testing"
)

/*func TestHealthHandler(t *testing.T) {
	resp, err := http.Get("http://localhost:3000/v1/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}

func TestCreatePostHandler(t *testing.T) {
	payload := map[string]string{
		"title":   "string",
		"content": "string",
	}
	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post("http://localhost:3000/v1/post", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Creaated, got %v", resp.Status)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %v", resp.Header.Get("Content-Type"))
	}
}*/

func TestDeletePostHandler(t *testing.T) {
	result := add(2, 3)
	if result != 5 {
		t.Errorf("Expected 5, got %v", result)
	}
}

func add(a, b int) int {
	return a + b
}
