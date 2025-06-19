// Example 18
package main

import (
	"errors"
	"fmt"
)

// ImmutableList implements an immutable collection using generics
type ImmutableList[T any] struct {
	data []T
}

// NewImmutableList creates a new immutable list with the provided items
func NewImmutableList[T any](items ...T) ImmutableList[T] {
	data := make([]T, len(items))
	copy(data, items)
	return ImmutableList[T]{data: data}
}

// Append creates a new list with additional items without modifying the original
func (l ImmutableList[T]) Append(items ...T) ImmutableList[T] {
	newData := make([]T, len(l.data)+len(items))
	copy(newData, l.data)
	copy(newData[len(l.data):], items)
	return ImmutableList[T]{data: newData}
}

// Get retrieves an item at the specified index
func (l ImmutableList[T]) Get(index int) (T, error) {
	if index < 0 || index >= len(l.data) {
		var zero T
		return zero, errors.New("index out of bounds")
	}
	return l.data[index], nil
}

// String provides a string representation of the list
func (l ImmutableList[T]) String() string {
	return fmt.Sprintf("%v", l.data)
}

func main() {
	// Create a new immutable list with some integers
	list1 := NewImmutableList(1, 2, 3)
	fmt.Println("Original list:", list1)

	// Append items without modifying the original
	list2 := list1.Append(4, 5)
	fmt.Println("Original list after append:", list1)
	fmt.Println("New list after append:", list2)

	// Access elements
	val, err := list2.Get(2)
	if err == nil {
		fmt.Println("Value at index 2:", val)
	}

	// Try to access an out-of-bounds index
	_, err = list2.Get(10)
	if err != nil {
		fmt.Println("Error:", err)
	}
}