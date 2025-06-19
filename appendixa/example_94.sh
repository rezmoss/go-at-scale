// Example 94
// Example 94
==================
WARNING: DATA RACE
Read at 0x00c000192038 by goroutine 8:
  main.(*Counter).Increment()
      /app/main.go:12 +0x3a

Previous write at 0x00c000192038 by goroutine 11:
  main.(*Counter).Increment()
      /app/main.go:8 +0x56
==================