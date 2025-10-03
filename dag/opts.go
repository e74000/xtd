package dag

import "context"

func defaultOpts() *Opts {
    return &Opts{
        Context:    context.Background(),
        MaxWorkers: 1,
    }
}

type Opts struct {
    Context    context.Context
    MaxWorkers int
}

func (o *Opts) apply(opts ...Opt) *Opts {
    for _, opt := range opts {
        opt(o)
    }

    return o
}

type Opt func(*Opts)

func WithContext(ctx context.Context) Opt {
    return func(o *Opts) {
        o.Context = ctx
    }
}

func WithMaxWorkers(max int) Opt {
    return func(o *Opts) {
        o.MaxWorkers = max
    }
}
