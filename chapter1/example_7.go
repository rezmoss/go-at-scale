// Example 7
import "errors"

var ErrNotFound = errors.New("resource not found")

func fetchResource(id string) error {
    // ... implementation ...
    return fmt.Errorf("fetching resource %s: %w", id, ErrNotFound)
}