// Example 18
type ImmutableList[T any] struct {
    data []T
}

func NewImmutableList[T any](items ...T) ImmutableList[T] {
    data := make([]T, len(items))
    copy(data, items)
    return ImmutableList[T]{data: data}
}

func (l ImmutableList[T]) Append(items ...T) ImmutableList[T] {
    newData := make([]T, len(l.data)+len(items))
    copy(newData, l.data)
    copy(newData[len(l.data):], items)
    return ImmutableList[T]{data: newData}
}

func (l ImmutableList[T]) Get(index int) (T, error) {
    if index < 0 || index >= len(l.data) {
        var zero T
        return zero, errors.New("index out of bounds")
    }
    return l.data[index], nil
}