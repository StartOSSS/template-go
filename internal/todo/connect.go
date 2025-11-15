package todo

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"
)

const serviceName = "todo.v1.TodoService"

// Handler exposes the TodoService Connect surface.
type Handler interface {
	CreateTodo(ctx context.Context, req *connect.Request[CreateTodoRequest]) (*connect.Response[CreateTodoResponse], error)
	GetTodo(ctx context.Context, req *connect.Request[GetTodoRequest]) (*connect.Response[GetTodoResponse], error)
	ListTodos(ctx context.Context, req *connect.Request[ListTodosRequest]) (*connect.Response[ListTodosResponse], error)
	UpdateTodo(ctx context.Context, req *connect.Request[UpdateTodoRequest]) (*connect.Response[UpdateTodoResponse], error)
	DeleteTodo(ctx context.Context, req *connect.Request[DeleteTodoRequest]) (*connect.Response[DeleteTodoResponse], error)
	HealthCheck(ctx context.Context, req *connect.Request[HealthCheckRequest]) (*connect.Response[HealthCheckResponse], error)
}

func NewConnectHandler(svc Handler, opts ...connect.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(connect.NewUnaryHandler("/"+serviceName+"/CreateTodo", svc.CreateTodo, opts...))
	mux.Handle(connect.NewUnaryHandler("/"+serviceName+"/GetTodo", svc.GetTodo, opts...))
	mux.Handle(connect.NewUnaryHandler("/"+serviceName+"/ListTodos", svc.ListTodos, opts...))
	mux.Handle(connect.NewUnaryHandler("/"+serviceName+"/UpdateTodo", svc.UpdateTodo, opts...))
	mux.Handle(connect.NewUnaryHandler("/"+serviceName+"/DeleteTodo", svc.DeleteTodo, opts...))
	mux.Handle(connect.NewUnaryHandler("/"+serviceName+"/HealthCheck", svc.HealthCheck, opts...))
	return "/" + serviceName + "/", mux
}
