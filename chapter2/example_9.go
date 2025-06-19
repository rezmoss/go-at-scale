// Example 9
package main

import (
	"fmt"
	"strings"
)

type StringTransformer func(string) string

func compose(f, g StringTransformer) StringTransformer {
	return func(s string) string {
		return f(g(s))
	}
}

// Example transformers
func removeSpaces(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func toLowerCase(s string) string {
	return strings.ToLower(s)
}

// Usage
func main() {
	cleanText := compose(removeSpaces, toLowerCase)
	result := cleanText("Hello World") // "helloworld"
	fmt.Println(result)               
}