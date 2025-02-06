// Example 94
==================
WARNING: DATA RACE
Read at 0x00c0000b4000 by goroutine 7:
  main.(*Counter).Value()
      /app/main.go:12 +0x3a

Previous write at 0x00c0000b4000 by goroutine 6:
  main.(*Counter).Increment()
      /app/main.go:8 +0x56
==================