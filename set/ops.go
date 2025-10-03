package set

func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
    s3 := New[T]()

    for t, _ := range s1.m {
        s3.m[t] = struct{}{}
    }

    for t, _ := range s2.m {
        s3.m[t] = struct{}{}
    }

    return s3
}

func Intersect[T comparable](s1, s2 *Set[T]) *Set[T] {
    s3 := New[T]()

    for t, _ := range s1.m {
        if _, ex := s2.m[t]; ex {
            s3.m[t] = struct{}{}
        }
    }

    return s3
}

func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
    s3 := New[T]()

    for t, _ := range s1.m {
        if _, ex := s2.m[t]; !ex {
            s3.m[t] = struct{}{}
        }
    }

    return s3
}
