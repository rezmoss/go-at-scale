// Example 17
// Immutable User type
type User struct {
    name string
    age  int
}

// Constructor ensures initialization
func NewUser(name string, age int) User {
    return User{name: name, age: age}
}

// WithName returns a new User instead of modifying
func (u User) WithName(name string) User {
    return User{name: name, age: u.age}
}

// WithAge returns a new User instead of modifying
func (u User) WithAge(age int) User {
    return User{name: u.name, age: age}
}

// Usage
func main() {
    user1 := NewUser("Alice", 30)
    user2 := user1.WithAge(31)  // user1 remains unchanged
}