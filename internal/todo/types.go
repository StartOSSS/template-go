package todo

// Request and response types mirror proto schema. They intentionally live in
// Go code so that we can bootstrap the repo without requiring protoc tooling.
type TodoTask struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateTodoResponse struct {
	Task *TodoTask `json:"task"`
}

type GetTodoRequest struct {
	Id string `json:"id"`
}

type GetTodoResponse struct {
	Task *TodoTask `json:"task"`
}

type ListTodosRequest struct{}

type ListTodosResponse struct {
	Tasks []*TodoTask `json:"tasks"`
}

type UpdateTodoRequest struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type UpdateTodoResponse struct {
	Task *TodoTask `json:"task"`
}

type DeleteTodoRequest struct {
	Id string `json:"id"`
}

type DeleteTodoResponse struct{}

type HealthCheckRequest struct{}

type HealthCheckResponse struct {
	Status string `json:"status"`
}
