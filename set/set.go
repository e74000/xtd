package set

func New[T comparable]() *Set[T] {
    return &Set[T]{m: make(map[T]struct{})}
}

type Set[T comparable] struct {
    m map[T]struct{}
}

func (s *Set[T]) Add(v T) {
    s.m[v] = struct{}{}
}

func (s *Set[T]) Has(v T) bool {
    _, ex := s.m[v]
    return ex
}

func (s *Set[T]) Delete(v T) {
    delete(s.m, v)
}

func (s *Set[T]) Values() []T {
    vs := make([]T, 0, len(s.m))
    for v, _ := range s.m {
        vs = append(vs, v)
    }
    return vs
}

func (s *Set[T]) Len() int {
    return len(s.m)
}

func (s *Set[T]) Clear() {
    clear(s.m)
}
