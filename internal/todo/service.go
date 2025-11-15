package todo

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	metric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Service implements the TodoService Connect handlers.
type Service struct {
	pool           *pgxpool.Pool
	logger         *slog.Logger
	tracer         trace.Tracer
	meter          metric.Meter
	requestCounter metric.Int64Counter
}

func NewService(pool *pgxpool.Pool, logger *slog.Logger) (*Service, error) {
	meter := otel.Meter("todo-service")
	counter, err := meter.Int64Counter("todo_requests_total")
	if err != nil {
		return nil, err
	}
	return &Service{
		pool:           pool,
		logger:         logger.With("component", "todo-service"),
		tracer:         otel.Tracer("todo-service"),
		meter:          meter,
		requestCounter: counter,
	}, nil
}

func (s *Service) CreateTodo(ctx context.Context, req *connect.Request[CreateTodoRequest]) (*connect.Response[CreateTodoResponse], error) {
	ctx, span := s.tracer.Start(ctx, "CreateTodo")
	defer span.End()

	s.requestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("method", "CreateTodo")))

	const q = `
        INSERT INTO todo_tasks (title, description)
        VALUES ($1, $2)
        RETURNING id, title, description, completed, created_at, updated_at;
    `
	row := s.pool.QueryRow(ctx, q, req.Msg.Title, req.Msg.Description)
	task, err := scanTask(row)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp := connect.NewResponse(&CreateTodoResponse{Task: task})
	resp.Header().Set("Todo-Version", "v1")
	return resp, nil
}

func (s *Service) GetTodo(ctx context.Context, req *connect.Request[GetTodoRequest]) (*connect.Response[GetTodoResponse], error) {
	ctx, span := s.tracer.Start(ctx, "GetTodo")
	defer span.End()

	const q = `
        SELECT id, title, description, completed, created_at, updated_at
        FROM todo_tasks WHERE id = $1;
    `
	row := s.pool.QueryRow(ctx, q, req.Msg.Id)
	task, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&GetTodoResponse{Task: task}), nil
}

func (s *Service) ListTodos(ctx context.Context, _ *connect.Request[ListTodosRequest]) (*connect.Response[ListTodosResponse], error) {
	ctx, span := s.tracer.Start(ctx, "ListTodos")
	defer span.End()

	const q = `
        SELECT id, title, description, completed, created_at, updated_at
        FROM todo_tasks ORDER BY created_at DESC;
    `
	rows, err := s.pool.Query(ctx, q)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	defer rows.Close()

	tasks := []*TodoTask{}
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		tasks = append(tasks, task)
	}

	return connect.NewResponse(&ListTodosResponse{Tasks: tasks}), nil
}

func (s *Service) UpdateTodo(ctx context.Context, req *connect.Request[UpdateTodoRequest]) (*connect.Response[UpdateTodoResponse], error) {
	ctx, span := s.tracer.Start(ctx, "UpdateTodo")
	defer span.End()

	const q = `
        UPDATE todo_tasks
        SET title = $1, description = $2, completed = $3, updated_at = NOW()
        WHERE id = $4
        RETURNING id, title, description, completed, created_at, updated_at;
    `
	row := s.pool.QueryRow(ctx, q, req.Msg.Title, req.Msg.Description, req.Msg.Completed, req.Msg.Id)
	task, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&UpdateTodoResponse{Task: task}), nil
}

func (s *Service) DeleteTodo(ctx context.Context, req *connect.Request[DeleteTodoRequest]) (*connect.Response[DeleteTodoResponse], error) {
	ctx, span := s.tracer.Start(ctx, "DeleteTodo")
	defer span.End()

	const q = `DELETE FROM todo_tasks WHERE id = $1;`
	cmd, err := s.pool.Exec(ctx, q, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if cmd.RowsAffected() == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("todo not found"))
	}
	return connect.NewResponse(&DeleteTodoResponse{}), nil
}

func (s *Service) HealthCheck(ctx context.Context, _ *connect.Request[HealthCheckRequest]) (*connect.Response[HealthCheckResponse], error) {
	ctx, span := s.tracer.Start(ctx, "HealthCheck")
	defer span.End()

	if err := s.pool.Ping(ctx); err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, err)
	}
	return connect.NewResponse(&HealthCheckResponse{Status: "ok"}), nil
}

func scanTask(row interface{ Scan(dest ...any) error }) (*TodoTask, error) {
	var (
		id, title, description string
		completed              bool
		createdAt, updatedAt   time.Time
	)
	if err := row.Scan(&id, &title, &description, &completed, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	return &TodoTask{
		Id:          id,
		Title:       title,
		Description: description,
		Completed:   completed,
		CreatedAt:   createdAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:   updatedAt.UTC().Format(time.RFC3339Nano),
	}, nil
}
