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
	basePath := "/" + serviceName + "/"

	register := func(name string, handler *connect.Handler) {
		mux.Handle(basePath+name, handler)
	}

	register("CreateTodo", connect.NewUnaryHandler(basePath+"CreateTodo", svc.CreateTodo, opts...))
	register("GetTodo", connect.NewUnaryHandler(basePath+"GetTodo", svc.GetTodo, opts...))
	register("ListTodos", connect.NewUnaryHandler(basePath+"ListTodos", svc.ListTodos, opts...))
	register("UpdateTodo", connect.NewUnaryHandler(basePath+"UpdateTodo", svc.UpdateTodo, opts...))
	register("DeleteTodo", connect.NewUnaryHandler(basePath+"DeleteTodo", svc.DeleteTodo, opts...))
	register("HealthCheck", connect.NewUnaryHandler(basePath+"HealthCheck", svc.HealthCheck, opts...))
	return basePath, mux
}
