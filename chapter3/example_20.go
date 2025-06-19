// Example 20
package main

import "fmt"

// Optional Type Implementation
type Optional[T any] struct {
	value   T
	present bool
}

func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, present: true}
}

func None[T any]() Optional[T] {
	return Optional[T]{present: false}
}


func (o Optional[T]) Map(f func(T) interface{}) Optional[interface{}] {
	if !o.present {
		return None[interface{}]()
	}
	return Some(f(o.value))
}

func (o Optional[T]) FlatMap(f func(T) Optional[interface{}]) Optional[interface{}] {
	if !o.present {
		return None[interface{}]()
	}
	return f(o.value)
}

func (o Optional[T]) String() string {
	if !o.present {
		return "None"
	}
	return fmt.Sprintf("Some(%v)", o.value)
}

// Example usage
func divide(a, b float64) Optional[float64] {
	if b == 0 {
		return None[float64]()
	}
	return Some(a / b)
}

func main() {

	fmt.Println("Dividing 10 by 2:")
	result1 := divide(10, 2)
	fmt.Println("Result of division:", result1)

	result2 := result1.Map(func(x float64) interface{} { return x * 2 })
	fmt.Println("After multiplying by 2:", result2)

	result3 := result2.Map(func(x interface{}) interface{} {
		return x.(float64) + 1
	})
	fmt.Println("After adding 1:", result3)

	// Show division by zero case
	fmt.Println("\nDividing 10 by 0:")
	result4 := divide(10, 0)
	fmt.Println("Result of division by zero:", result4)

	// Chain operations with None result
	result5 := result4.Map(func(x float64) interface{} {
		return x * 2
	}).Map(func(x interface{}) interface{} {
		return x.(float64) + 1
	})
	fmt.Println("After chained operations on None:", result5)
}