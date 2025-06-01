-- name: ListTodos :many
SELECT id, description, completed FROM todos ORDER BY id;

-- name: MarkTodoDone :exec
UPDATE todos SET completed = true WHERE id = $1;

-- name: CreateTodo :one
INSERT INTO todos (description) VALUES ($1)
RETURNING id, description, completed;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = $1;