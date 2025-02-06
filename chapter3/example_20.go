// Example 20
type Optional[T any] struct {
    value T
    present bool
}

func Some[T any](value T) Optional[T] {
    return Optional[T]{value: value, present: true}
}

func None[T any]() Optional[T] {
    return Optional[T]{present: false}
}

func (o Optional[T]) Map[U any](f func(T) U) Optional[U] {
    if !o.present {
        return None[U]()
    }
    return Some(f(o.value))
}

func (o Optional[T]) FlatMap[U any](f func(T) Optional[U]) Optional[U] {
    if !o.present {
        return None[U]()
    }
    return f(o.value)
}

// Example usage
func divide(a, b float64) Optional[float64] {
    if b == 0 {
        return None[float64]()
    }
    return Some(a / b)
}

func main() {
    result := divide(10, 2).
        Map(func(x float64) float64 { return x * 2 }).
        Map(func(x float64) float64 { return x + 1 })
}