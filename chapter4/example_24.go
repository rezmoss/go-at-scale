// Example 24

//use errgroup to wait for all goroutines to finish

import (
    "golang.org/x/sync/errgroup"
    // ...
)

// Instead of manual wait groups and channels, use errgroup:

g, ctx := errgroup.WithContext(ctx)

for _, item := range items {
    it := item // local copy
    g.Go(func() error {
        return processItem(ctx, it)
    })
}

// If any goroutine returns an error, g.Wait() returns it.
if err := g.Wait(); err != nil {
    // Handle the first error from any goroutine
    log.Printf("Error processing items: %v", err)
} else {
    // All goroutines succeeded
}