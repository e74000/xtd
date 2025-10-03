package dag

import "context"

type ID string

type Task struct {
    ID   ID
    Deps []ID
    Fn   Runner
}

type Runner func(ctx context.Context) error
