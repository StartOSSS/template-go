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
	handlerOpts := append([]connect.HandlerOption{}, opts...)
	for _, codec := range codecsForJSON() {
		handlerOpts = append(handlerOpts, connect.WithCodec(codec))
	}

	register := func(name string, handler *connect.Handler) {
		mux.Handle(basePath+name, handler)
	}

	register("CreateTodo", connect.NewUnaryHandler(basePath+"CreateTodo", svc.CreateTodo, handlerOpts...))
	register("GetTodo", connect.NewUnaryHandler(basePath+"GetTodo", svc.GetTodo, handlerOpts...))
	register("ListTodos", connect.NewUnaryHandler(basePath+"ListTodos", svc.ListTodos, handlerOpts...))
	register("UpdateTodo", connect.NewUnaryHandler(basePath+"UpdateTodo", svc.UpdateTodo, handlerOpts...))
	register("DeleteTodo", connect.NewUnaryHandler(basePath+"DeleteTodo", svc.DeleteTodo, handlerOpts...))
	register("HealthCheck", connect.NewUnaryHandler(basePath+"HealthCheck", svc.HealthCheck, handlerOpts...))
	return basePath, mux
}
