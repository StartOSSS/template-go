package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type createTodoResponse struct {
	Task struct {
		Id string `json:"id"`
	} `json:"task"`
}

func main() {
	endpoint := env("INTEGRATION_TARGET", "http://todo-app-api.todo.svc.cluster.local:8080")
	payload := map[string]any{"title": "integration", "description": "client"}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(endpoint+"/todo.v1.TodoService/CreateTodo", "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("unexpected status %d", resp.StatusCode))
	}
	var parsed createTodoResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		panic(err)
	}
	fmt.Printf("created todo %s\n", parsed.Task.Id)
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
