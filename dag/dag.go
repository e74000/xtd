package dag

import (
    "fmt"
    "sync"

    "github.com/e74000/xtd/set"
    "golang.org/x/sync/errgroup"
)

var (
    ErrDuplicateTask = fmt.Errorf("duplicate task")
    ErrSelfDep       = fmt.Errorf("self dependency")
    ErrUnknownTask   = fmt.Errorf("unknown")
    ErrCycle         = fmt.Errorf("cycle")
    ErrTask          = fmt.Errorf("task")
)

func newDag(tasks []Task) (*dag, error) {
    d := &dag{
        m:   make(map[ID]Task),
        adj: make(map[ID][]ID),
        deg: make(map[ID]int),
    }

    for _, t := range tasks {
        if _, ex := d.m[t.ID]; ex {
            return nil, fmt.Errorf("%w: %s", ErrDuplicateTask, t.ID)
        }
        d.m[t.ID] = t
        d.deg[t.ID] = 0
    }

    for _, t := range tasks {
        seen := set.New[ID]()
        for _, dep := range t.Deps {
            if t.ID == dep {
                return nil, fmt.Errorf("%w: %s", ErrSelfDep, t.ID)
            }
            if _, ex := d.m[dep]; !ex {
                return nil, fmt.Errorf("%w: %s", ErrUnknownTask, dep)
            }
            if seen.Has(dep) {
                continue
            }
            seen.Add(dep)
            d.deg[dep]++
            d.adj[dep] = append(d.adj[dep], t.ID)
        }
    }

    return d, nil
}

type dag struct {
    m   map[ID]Task
    adj map[ID][]ID
    deg map[ID]int
}

func Run(tasks []Task, options ...Opt) error {
    if len(tasks) == 0 {
        return nil
    }

    opts := defaultOpts().apply(options...)

    d, err := newDag(tasks)
    if err != nil {
        return err
    }

    ready := make([]ID, 0, len(d.m))
    for id, deg := range d.deg {
        if deg == 0 {
            ready = append(ready, id)
        }
    }

    if len(ready) == 0 {
        return ErrCycle
    }

    g, gctx := errgroup.WithContext(opts.Context)
    if opts.MaxWorkers > 0 {
        g.SetLimit(opts.MaxWorkers)
    }

    var mu sync.Mutex
    finished := 0

    var schedule func(ID)
    schedule = func(id ID) {
        task := d.m[id]
        g.Go(func() error {
            select {
            case <-gctx.Done():
                return gctx.Err()
            default:
            }

            if err := task.Fn(gctx); err != nil {
                return fmt.Errorf("%w %s: %w", ErrTask, task.ID, err)
            }

            var nowReady []ID
            mu.Lock()
            finished++
            for _, dep := range d.adj[id] {
                d.deg[dep]--
                if d.deg[dep] == 0 {
                    nowReady = append(nowReady, dep)
                }
            }
            mu.Unlock()

            for _, t := range nowReady {
                schedule(t)
            }

            return nil
        })
    }

    for _, id := range ready {
        schedule(id)
    }

    if err := g.Wait(); err != nil {
        return err
    }

    if finished != len(tasks) {
        var stuck []ID
        for id, deg := range d.deg {
            if deg > 0 {
                stuck = append(stuck, id)
            }
        }
        return fmt.Errorf("%w: %v", ErrCycle, stuck)
    }

    return nil
}
