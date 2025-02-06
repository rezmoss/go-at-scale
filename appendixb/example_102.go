// Example 102
// Pitfall 1: Hidden memory leaks in slices
func getBigSlicePitfall(data []int) []int {
    return data[:2]  // Still holds reference to original array!
}

// Solution: Copy what you need
func getBigSliceSolution(data []int) []int {
    result := make([]int, 2)
    copy(result, data)
    return result
}

// Pitfall 2: Unexpected slice behavior
func appendPitfall() {
    x := make([]int, 0, 10)
    x = append(x, 1)
    y := x
    y = append(y, 2)
    // x[0] is now also 2!
}

// Solution: Use full slices or copy
func appendSolution() {
    x := make([]int, 0, 10)
    x = append(x, 1)
    y := make([]int, len(x))
    copy(y, x)
    y = append(y, 2)
}

// Pitfall 3: Growing slices inefficiently
func growSlicePitfall(items []int) []int {
    for i := 0; i < 1000; i++ {
        items = append(items, i)  // Grows one by one
    }
    return items
}

// Solution: Pre-allocate when size is known
func growSliceSolution(items []int) []int {
    newItems := make([]int, len(items), len(items)+1000)
    copy(newItems, items)
    for i := 0; i < 1000; i++ {
        newItems = append(newItems, i)
    }
    return newItems
}